package middleware

import (
	"net/http"

	"backend/internal/contexthelper"
)

func Version(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cfg := contexthelper.GetConfig(ctx)
		w.Header().Set("X-App-Sha", cfg.CommitSha)
		w.Header().Set("X-App-Version", cfg.GitTag)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
