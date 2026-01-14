package middleware

import (
	"backend/internal/contexthelper"
	"backend/internal/response"
	"backend/pkg/logger"
	"net/http"
)

func AuthOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userId, _ := contexthelper.GetUserId(ctx)
		if userId == 0 {
			logger.WarnCtx(ctx, "The required UserId value is not present in context")
			response.UnauthorizedErrorResponse(w, "Access token is missing or malformed")
			return
		}
		logger.DebugCtx(ctx, "The required UserId value is present in context : %d", userId)
		next.ServeHTTP(w, r)
	})
}
