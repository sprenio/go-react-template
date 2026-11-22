package logger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

type requestIdGetter = func(ctx context.Context) string

var (
	slogLogger = slog.New(&CustomHandler{
		opts: slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
		out:    os.Stdout,
		errOut: os.Stderr,
	})
	getRequestId requestIdGetter
)

const RequestIdLogPrefix = "reqID"

// CustomHandler implementuje interfejs slog.Handler
type CustomHandler struct {
	opts         slog.HandlerOptions
	out          io.Writer
	errOut       io.Writer
	preAttrs     []slog.Attr
	getRequestId requestIdGetter
}

type ErrorAttr struct {
	Key   string
	Value any
}

type errorWithAttrs interface {
	Attrs() []ErrorAttr
}
type errorWithCaller interface {
	Caller() string
}

func extractErrorAttrs(err error) []any {
	var allAttrs []any
	seenKeys := make(map[string]bool)

	for err != nil {
		// Sprawdzamy czy błąd implementuje interfejs z atrybutami
		if ew, ok := err.(errorWithAttrs); ok {
			for _, attr := range ew.Attrs() {
				if !seenKeys[attr.Key] {
					allAttrs = append(allAttrs, attr.Key, attr.Value)
					seenKeys[attr.Key] = true
				}
			}
		}
		err = errors.Unwrap(err)
	}
	return allAttrs
}

// extractErrorStack zbiera ślad ze wszystkich poziomów łańcucha błędów
func extractErrorStack(err error) []string {
	var stack []string
	for err != nil {
		if ew, ok := err.(errorWithCaller); ok {
			if loc := ew.Caller(); loc != "" {
				stack = append(stack, loc)
			}
		}
		err = errors.Unwrap(err)
	}
	return stack
}

func (h *CustomHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	// 1. Czas w UTC RFC3339Nano
	timestamp := r.Time.UTC().Format(time.RFC3339Nano)

	// 2. Poziom logowania
	levelStr := "[" + strings.ToUpper(r.Level.String()) + "]"

	// 3. Request ID z kontekstu
	requestID := "-"
	if h.getRequestId != nil && ctx != nil {
		if reqId := h.getRequestId(ctx); reqId != "" {
			requestID = reqId
		}
	}

	// 4. Lokalizacja kodu
	location := getCallerLocation()

	// 5. WSZYSTKE przekazane parametry (key-value)
	var attrStrs []string

	formatAttr := func(a slog.Attr) {
		// Sprawdzamy, czy wartością atrybutu jest błąd (np. "err": err)
		if errVal, ok := a.Value.Any().(error); ok {
			// 1. Wyciągamy atrybuty (query, userID itp.)
			errorAttrs := extractErrorAttrs(errVal)
			for i := 0; i < len(errorAttrs); i += 2 {
				k := errorAttrs[i].(string)
				v := errorAttrs[i+1]
				attrStrs = append(attrStrs, fmt.Sprintf("%s=%v", k, v))
			}
			// 2. Wyciągamy backtrace z błędów
			if stack := extractErrorStack(errVal); len(stack) > 0 {
				// Zapisujemy jako string z rozdzieleniem, np. callerStack="service.go:40 -> db_query.go:25"
				attrStrs = append(attrStrs, fmt.Sprintf("errorStack=%q", strings.Join(stack, " -> ")))
			}
		}
		attrStrs = append(attrStrs, fmt.Sprintf("%s=%v", a.Key, a.Value.Any()))
	}

	for _, a := range h.preAttrs {
		formatAttr(a)
	}
	r.Attrs(func(a slog.Attr) bool {
		formatAttr(a)
		return true
	})

	// 6. Finalny format loga
	logLine := fmt.Sprintf("%s %s [%s] %s :: %s",
		timestamp,
		levelStr,
		requestID,
		location,
		r.Message,
	)
	if len(attrStrs) > 0 {
		extras := strings.Join(attrStrs, " ")
		logLine += " :: " + extras
	}
	logLine += "\n"

	target := h.out
	if r.Level >= slog.LevelError {
		target = h.errOut
	}

	_, err := target.Write([]byte(logLine))
	return err
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.preAttrs)+len(attrs))
	copy(newAttrs, h.preAttrs)
	copy(newAttrs[len(h.preAttrs):], attrs)
	return &CustomHandler{
		opts:         h.opts,
		out:          h.out,
		errOut:       h.errOut,
		preAttrs:     newAttrs,
		getRequestId: h.getRequestId,
	}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return h
}

func getCallerLocation() string {
	pcs := make([]uintptr, 15)
	l := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:l])
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "pkg/logger/") && !strings.Contains(frame.File, "log/slog") {
			return fmt.Sprintf("%s:%d", frame.File, frame.Line)
		}
		if !more {
			break
		}
	}
	return "unknown"
}

func Init(level string, rqIdGetter requestIdGetter) {
	getRequestId = rqIdGetter

	var slogLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	handler := &CustomHandler{
		opts: slog.HandlerOptions{
			Level: slogLevel,
		},
		out:          os.Stdout,
		errOut:       os.Stderr,
		getRequestId: rqIdGetter,
	}

	slogLogger = slog.New(handler)
	Info("Logger initialized with level: " + level)
}

// --- Publiczne funkcje pomocnicze (bez formatMsg, czyste przekazanie args...) ---

func Info(msg string, args ...any) {
	slogLogger.Log(context.Background(), slog.LevelInfo, msg, args...)
}
func Error(msg string, args ...any) {
	slogLogger.Log(context.Background(), slog.LevelError, msg, args...)
}
func Debug(msg string, args ...any) {
	slogLogger.Log(context.Background(), slog.LevelDebug, msg, args...)
}
func Fatal(msg string, args ...any) {
	slogLogger.Log(context.Background(), slog.LevelError, msg, args...)
	os.Exit(1)
}
func Warn(msg string, args ...any) {
	slogLogger.Log(context.Background(), slog.LevelWarn, msg, args...)
}

func InfoCtx(ctx context.Context, msg string, args ...any) {
	slogLogger.Log(ctx, slog.LevelInfo, msg, args...)
}
func ErrorCtx(ctx context.Context, msg string, args ...any) {
	slogLogger.Log(ctx, slog.LevelError, msg, args...)
}
func DebugCtx(ctx context.Context, msg string, args ...any) {
	slogLogger.Log(ctx, slog.LevelDebug, msg, args...)
}
func FatalCtx(ctx context.Context, msg string, args ...any) {
	slogLogger.Log(ctx, slog.LevelError, msg, args...)
	os.Exit(1)
}
func WarnCtx(ctx context.Context, msg string, args ...any) {
	slogLogger.Log(ctx, slog.LevelWarn, msg, args...)
}
