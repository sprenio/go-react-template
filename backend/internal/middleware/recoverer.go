package middleware

import (
	"backend/internal/response"
	"backend/pkg/logger"
	"net/http"
)

func Recoverer(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
				logger.ErrorCtx(r.Context(), "Panic recovered: %v", err)
				response.InternalServerError(w)
                // tu możesz zalogować err + stacktrace
            }
        }()
        next.ServeHTTP(w, r)
    })
}
