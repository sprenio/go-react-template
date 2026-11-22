package http

import (
	"backend/internal/contexthelper"
	"backend/internal/response"
	"backend/internal/user/dto"
	"backend/internal/user/service"
	"backend/pkg/logger"
	"net/http"
)

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userID, ok := contexthelper.GetUserId(ctx)
	logger.InfoCtx(ctx, "[ME] User ID from context", "userId", userID, "present", ok)
	if !ok {
		response.UnauthorizedErrorResponse(w, "User not found")
		return
	}

	service := service.NewService(ctx)

	user, err := service.GetUserResponseData(ctx, userID)
	if err != nil {
		logger.ErrorCtx(ctx, "Get user failed", "error", err)
		response.InternalServerError(w)
		return
	}
	dto.SetMeSuccessResponse(w, ctx, user)
}
