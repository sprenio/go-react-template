package middleware

import (
	"net/http"

	"backend/internal/contexthelper"
	"backend/pkg/uuidstr"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuidstr.GetUniqBase36(32)
		// add to context
		ctx := contexthelper.SetRequestID(r.Context(), requestID)
		// optionally: add to response headers
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

