package http

import (
	"backend/internal/contexthelper"
	"backend/internal/cookie"
	"backend/internal/session/service"
	"backend/internal/user/dto"
	"backend/pkg/logger"

	"net/http"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := contexthelper.GetUserId(ctx)
	logger.DebugCtx(ctx, "logout userId: %d", userId)
	if !ok {
		dto.LogoutFailedResponse(w)
		return
	}
	sessionService := service.NewService(ctx, w)
	token := cookie.GetRefreshToken(r)
	logger.DebugCtx(ctx, "logout token: %s", token)
	err := sessionService.Logout(ctx, token)
	accessTokenData, ctx := contexthelper.GetAccessTokenData(ctx)
	accessTokenData.SetCookies = false

	if err != nil {
		logger.ErrorCtx(ctx, "Logout failed", "error", err)
		dto.LogoutFailedResponse(w)
		return
	}
	logger.InfoCtx(ctx, "User logged out successfully", "userId", userId)
	dto.SetLogoutSuccessResponse(w, ctx)
}
