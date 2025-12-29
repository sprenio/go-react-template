package middleware

import (
	"net/http"

	"backend/internal/contexthelper"
	"backend/pkg/uuidstr"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuidstr.GetUniqBase36(32)
		// dodaj do kontekstu
		ctx := contexthelper.SetRequestID(r.Context(), requestID)
		// opcjonalnie: dodaj do nagłówków odpowiedzi
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

