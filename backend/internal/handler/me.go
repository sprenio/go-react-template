package handler

import (
	"backend/internal/contexthelper"
	"backend/internal/repository"
	"backend/internal/response"
	"backend/internal/service"
	"backend/pkg/logger"
	"net/http"
)

func (h *Handler) MeHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userID, ok := contexthelper.GetUserId(ctx)
	logger.InfoCtx(ctx, "[ME] User ID from context: %d, present: %v", userID, ok)
	if !ok {
		response.UnauthorizedErrorResponse(w, "User not found")
		return
	}
	userRepo := repository.NewUserRepository(h.db)
	userService := service.NewUserService(userRepo)

	user, err := userService.GetUserResponseData(ctx, userID)
	if err != nil {
		logger.ErrorCtx(ctx, "Get user failed: %v", err)
		response.InternalServerError(w)
		return
	}
	response.SetMeSuccessResponse(w, ctx, user)
}
