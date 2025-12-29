package logger

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type callerAwareWriter struct {
	level string
	out   io.Writer
}

func (w *callerAwareWriter) Write(p []byte) (n int, err error) {
	location := getCallerLocaion();

	timestamp := time.Now().UTC().Format(time.RFC3339Nano)

	// Usuń newline z p[], bo dodamy własny
	rawMessage := strings.TrimSuffix(string(p), "\n")
	// Extract and remove request ID
	requestID := extractRequestID(rawMessage)
	message := stripRequestID(rawMessage)

	// Finalny format loga
	logLine := fmt.Sprintf("%s [%s] [%s] %s :: %s\n",
		timestamp,                // 2025-07-23T13:10:45.123Z
		strings.ToUpper(w.level), // [INFO]
		requestID,                // [reqID: ...]
		location,                 // handler.go:25
		message,                  // faktyczna treść
	)

	return w.out.Write([]byte(logLine))
}

func extractRequestID(msg string) string {
	start := strings.Index(msg, "["+RequestIdLogPrefix)
	if start == -1 {
		return "-"
	}
	end := strings.Index(msg[start:], "]")
	if end == -1 {
		return "-"
	}
	return msg[start+8 : start+end]
}

func stripRequestID(msg string) string {
	start := strings.Index(msg, "["+RequestIdLogPrefix)
	if start == -1 {
		return msg
	}
	end := strings.Index(msg[start:], "]")
	if end == -1 {
		return msg
	}
	return strings.TrimSpace(msg[:start] + msg[start+end+1:])
}

func getCallerLocaion() string {
	pcs := make([]uintptr, 10)          // tablica na 10 ramek
	l := runtime.Callers(2, pcs)        // skip=2 → pomija Callers i getCallerLocaion
	frames := runtime.CallersFrames(pcs[:l]) // uzyskaj ramki
	location := "unknown"
	for {
		frame, more := frames.Next()
		location = fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
		// Pomijaj ramki z tego pakietu (logger)
		if strings.HasPrefix(frame.File, "/app/") && !strings.Contains(frame.File, "pkg/logger/") {
			return fmt.Sprintf("%s:%d", frame.File, frame.Line)
		}
		if !more {
			break
		}
	}
	return location
}