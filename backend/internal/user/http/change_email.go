package http

import (
	"backend/internal/apierrors"
	"backend/internal/contexthelper"
	"backend/internal/helper"
	"backend/internal/response"
	"backend/internal/user/dto"
	"backend/internal/user/service"
	"backend/pkg/logger"
	"backend/pkg/validation"
	"context"
	"net/http"
)

type EmailChangeRequest struct {
	Email string `json:"email"`
}

func (req *EmailChangeRequest) Validate(ctx context.Context) error {
	req.Email = helper.CleanEmail(req.Email)

	err := validation.IsEmailValid(req.Email)
	if err != nil {
		return apierrors.NewInvalidInputError("Email", err.Error(), req)
	}
	userId, ok := contexthelper.GetUserId(ctx)
	if !ok {
		logger.ErrorCtx(ctx, "User not set")
		return apierrors.NewGeneralCustomError("User not authenticated")
	}
	userService := service.NewService(ctx)
	isEmailAvailable, err := userService.IsEmailAvailableForUser(ctx, userId, req.Email)
	if err != nil {
		return err
	}
	if !isEmailAvailable {
		return apierrors.NewInvalidInputError("Email", "email is not available", &req)
	}

	return nil
}

func (h *Handler) ChangeEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := contexthelper.GetUserId(ctx)
	if !ok {
		logger.ErrorCtx(ctx, "User not set")
		response.NotFoundErrorResponse(w)
		return
	}
	rqData, err := helper.GetRequestData[EmailChangeRequest](r, w)
	if err != nil {
		response.InvalidInputResponse(w, err)
		return
	}
	userService := service.NewService(ctx)
	logger.DebugCtx(ctx, "ChangeEmail called with data: %v", rqData)

	err = userService.EnqueueChangeEmail(ctx, userId, rqData.Email)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to change email", "error", err)
		response.InternalServerError(w)
		return
	}
	logger.InfoCtx(ctx, "Change email successfully", "email", rqData.Email)
	dto.SetChangeEmailSuccessResponse(w, ctx)
}
