package middleware

import (
	"backend/pkg/logger"
	"net/http"
	"time"
)

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		frontendBase := r.Header.Get("X-Frontend-Base-URL")
		logger.InfoCtx(r.Context(), r.Method+" "+frontendBase+r.URL.Path)
		next.ServeHTTP(w, r)
		duration := time.Since(startTime)
		logger.InfoCtx(r.Context(), r.Method+" "+frontendBase+r.URL.Path+" :: Executed in %v", duration)
	})
}
