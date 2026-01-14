package handler

import (
	"backend/internal/contexthelper"
	"backend/internal/cookie"
	"backend/internal/repository"
	"backend/internal/response"
	"backend/internal/service"
	"backend/pkg/logger"

	"net/http"
)

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := contexthelper.GetUserId(ctx)
	logger.DebugCtx(ctx, "logout userId: %d", userId)
	if !ok {
		response.LogoutFailedResponse(w)
		return
	}
	db := contexthelper.GetDb(ctx)
	sessionRepo := repository.NewUserSessionsRepository(db)
	sessionService := service.NewSessionService(sessionRepo, w)
	token := cookie.GetRefreshToken(r)
	logger.DebugCtx(ctx, "logout token: %s", token)
	err := sessionService.Logout(ctx, token)
	accessTokenData, ctx := contexthelper.GetAccessTokenData(ctx)
	accessTokenData.SetCookies = false
	
	if err != nil {
		logger.ErrorCtx(ctx, "Logout failed: %v", err)
		response.LogoutFailedResponse(w)
		return
	}
	logger.InfoCtx(ctx, "User %d logged out successfully", userId)
	response.SetLogoutSuccessResponse(w, ctx)
}
