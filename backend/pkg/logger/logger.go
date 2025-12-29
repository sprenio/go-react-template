package logger

import (
	"backend/internal/contexthelper"
	"context"
	"log"
	"os"
)

var (
	infoLogger  = log.New(&callerAwareWriter{"info", os.Stdout}, "", 0)
	errorLogger = log.New(&callerAwareWriter{"error", os.Stderr}, "", 0)
	debugLogger = log.New(&callerAwareWriter{"debug", os.Stdout}, "", 0)
	fatalLogger = log.New(&callerAwareWriter{"fatal", os.Stdout}, "", 0)
	warnLogger  = log.New(&callerAwareWriter{"warn", os.Stdout}, "", 0)

	logLevel = "info" // default
)

const RequestIdLogPrefix = "reqID:"

func Init(level string) {
	logLevel = level
	infoLogger.Println("Logger initialized with level:", level)
}

func Info(msg string, args ...any) {
	if logLevel == "info" || logLevel == "debug" {
		infoLogger.Printf(msg, args...)
	}
}

func Error(msg string, args ...any) {
	errorLogger.Printf(msg, args...)
}

func Debug(msg string, args ...any) {
	if logLevel == "debug" {
		debugLogger.Printf(msg, args...)
	}
}
func Fatal(msg string, args ...any) {
	fatalLogger.Fatalf(msg, args...)
}
func Warn(msg string, args ...any) {
	if logLevel == "warn" || logLevel == "debug" || logLevel == "info" {
		warnLogger.Printf(msg, args...)
	}
}

// --- nowe funkcje z contextem ---
func InfoCtx(ctx context.Context, msg string, args ...any) {
	if logLevel == "info" || logLevel == "debug" {
		prefix := prefixFromContext(ctx)
		infoLogger.Printf(prefix+msg, args...)
	}
}

func ErrorCtx(ctx context.Context, msg string, args ...any) {
	prefix := prefixFromContext(ctx)
	errorLogger.Printf(prefix+msg, args...)
}

func DebugCtx(ctx context.Context, msg string, args ...any) {
	if logLevel == "debug" {
		prefix := prefixFromContext(ctx)
		debugLogger.Printf(prefix+msg, args...)
	}
}

func FatalCtx(ctx context.Context, msg string, args ...any) {
	prefix := prefixFromContext(ctx)
	fatalLogger.Printf(prefix+msg, args...)
}

func WarnCtx(ctx context.Context, msg string, args ...any) {
	if logLevel == "warn" || logLevel == "debug" || logLevel == "info" {
		prefix := prefixFromContext(ctx)
		warnLogger.Printf(prefix+msg, args...)
	}
}

func prefixFromContext(ctx context.Context) string {
	if reqID, ok := ctx.Value(contexthelper.RequestIDKey).(string); ok {
		return "[" + RequestIdLogPrefix + " " + reqID + "] "
	}
	return ""
}
