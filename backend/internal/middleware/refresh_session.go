package middleware

import (
	"backend/internal/contexthelper"
	"backend/internal/cookie"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/logger"
	"net/http"
)

func RefreshSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sw := &statusResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK, // domyślny
		}
		accessTokenData, ctx := contexthelper.GetAccessTokenData(ctx)
		refreshToken := cookie.GetRefreshToken(r)
		logger.DebugCtx(ctx, "middleware refreshToken: %s", refreshToken)
		accessTokenData.RefreshToken = refreshToken
		userId, ok := contexthelper.GetUserId(ctx)
		logger.DebugCtx(ctx, "middleware accesTokenData: %v, userid: %d", accessTokenData, userId)
		if !ok && refreshToken != "" {
			var err error
			db := contexthelper.GetDb(ctx)
			sessionRepo := repository.NewUserSessionsRepository(db)
			sessionService := service.NewSessionService(sessionRepo, w)
			userId, err = sessionService.Login(ctx, refreshToken)
			if err != nil {
				logger.ErrorCtx(ctx, "Login user by refresh token %v failed: %v", refreshToken, err)
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
