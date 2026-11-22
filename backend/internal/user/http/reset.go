package http

import (
	"context"
	"net/http"

	"backend/internal/apierrors"
	"backend/internal/helper"
	"backend/internal/response"
	"backend/internal/user/dto"
	"backend/internal/user/service"
	"backend/pkg/logger"
	"backend/pkg/validation"
)

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

func (req *ResetPasswordRequest) Validate(ctx context.Context) error {
	req.Email = helper.CleanEmail(req.Email)
	err := validation.IsEmailValid(req.Email)
	if err != nil {
		return apierrors.NewInvalidInputError("Email", err.Error(), req)
	}
	return nil
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rqData, err := helper.GetRequestData[ResetPasswordRequest](r, w)
	if err != nil {
		response.InvalidInputResponse(w, err)
		return
	}
	userService := service.NewService(ctx)
	logger.InfoCtx(ctx, "Reset password", "email", rqData.Email)
	err = userService.EnqueueResetPassword(ctx, rqData.Email)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to reset password", "error", err)
		response.InternalServerError(w)
		return
	}
	logger.InfoCtx(ctx, "Reset password successfully", "email", rqData.Email)
	dto.SetResetPasswordSuccessResponse(w, ctx)
}
