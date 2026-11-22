package middleware

import (
	"backend/internal/contexthelper"
	"backend/internal/cookie"
	"backend/internal/session/service"
	"backend/pkg/logger"
	"net/http"
)

func RefreshSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sw := &statusResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK, // default
		}
		accessTokenData, ctx := contexthelper.GetAccessTokenData(ctx)
		refreshToken := cookie.GetRefreshToken(r)
		accessTokenData.RefreshToken = refreshToken
		userId, ok := contexthelper.GetUserId(ctx)
		if !ok && refreshToken != "" {
			var err error
			sessionService := service.NewService(ctx, w)
			userId, err = sessionService.Login(ctx, refreshToken)
			if err != nil {
				logger.ErrorCtx(ctx, "Login user by refresh token failed", "refreshToken", refreshToken, "error", err)
				sessionService.Logout(ctx, refreshToken)
				refreshToken = ""
			}
			ctx = contexthelper.SetUserId(ctx, userId)
		}
		accessTokenData.UserId = userId
		accessTokenData.SetCookies = userId > 0
		next.ServeHTTP(sw, r.WithContext(ctx))
	})
}
