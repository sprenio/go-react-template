package middleware

import (
	"backend/internal/httpx"
	"backend/internal/response"
	"backend/pkg/logger"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if _, ok := err.(httpx.AbortRequest); ok {
					return
				}
				filteredStack := getFilteredStack()
				logger.ErrorCtx(
					r.Context(),
					fmt.Sprintf("Panic recovered during %s %s: %v\nFiltered Stack:\n%s",
						r.Method,
						r.URL.RequestURI(),
						err,
						filteredStack),
				)
				response.InternalServerError(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func getFilteredStack() string {
	pcs := make([]uintptr, 64)
	n := runtime.Callers(3, pcs) // Skip Recoverer and debug.Stack equivalents
	frames := runtime.CallersFrames(pcs[:n])

	var filteredStack strings.Builder
	count := 0
	for {
		frame, more := frames.Next()

		if strings.Contains(frame.File, "backend/internal") || strings.Contains(frame.File, "/app/internal") {
			fmt.Fprintf(&filteredStack, "%s:%d %s\n", frame.File, frame.Line, frame.Function)
			count++
		}
		if !more || count >= 15 {
			break
		}
	}
	return filteredStack.String()
}
