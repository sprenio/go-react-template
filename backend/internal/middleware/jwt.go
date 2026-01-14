package middleware

import (
	"backend/internal/contexthelper"
	"backend/internal/cookie"

	"backend/pkg/logger"
	"net/http"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userId, err := cookie.GetUserIdFromJwtToken(r)
		accessTokenData, ctx := contexthelper.GetAccessTokenData(ctx)
		if err == nil {
			logger.InfoCtx(ctx, "Authenticated user ID: %d", userId)
			// dodaj user_id do kontekstu
			ctx = contexthelper.SetUserId(ctx, userId)
			accessTokenData.SetCookies = true
			accessTokenData.UserId = userId
		} else {
			logger.ErrorCtx(ctx, "Get User Id from JWT token failed: %v", err)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
