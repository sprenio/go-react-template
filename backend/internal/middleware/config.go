package middleware

import (
	"backend/config"
	"backend/internal/contexthelper"
	"net/http"
)
func Config(cfg *config.Config) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := contexthelper.SetConfig(r.Context(), cfg)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}