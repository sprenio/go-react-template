package http

import (
	"backend/internal/apierrors"
	"backend/internal/confirmation/dto"
	"backend/internal/confirmation/service"
	"backend/internal/response"
	"backend/pkg/logger"
	"backend/pkg/validation"
	"encoding/json"
	"net/http"
)

type PasswordChangeRequest struct {
	Password string `json:"password"`
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	tokenData, err := h.getRequestTokenData(r)
	ctx := r.Context()
	if err != nil {
		logger.WarnCtx(ctx, "Failed to get request token data", "error", err)
		response.NotFoundErrorResponse(w)
		return
	}

	service := service.NewService(ctx)
	var req PasswordChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.InvalidJsonErrorResponse(w)
		return
	}
	logger.DebugCtx(ctx, "PasswordChangeHandler called with data: %v", req)
	err = validation.IsPasswordValid(req.Password)
	if err != nil {
		respErr := apierrors.NewInvalidInputError("Password", err.Error(), &req)
		response.InvalidInputValueErrorResponse(w, respErr)
		return
	}

	err = service.ChangePassword(ctx, tokenData.UserId, req.Password)

	if err != nil {
		logger.ErrorCtx(ctx, "Failed to change password", "error", err)
		response.InternalServerError(w)
		return
	}
	err = service.ConsumeToken(ctx, tokenData.Token)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to consume confirmation token", "error", err)
		response.InternalServerError(w)
		return
	}

	dto.ChangePasswordSuccessResponse(w, ctx)
}
