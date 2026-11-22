package http

import (
	"net/http"

	"backend/internal/confirmation"
	"backend/internal/confirmation/dto"
	"backend/internal/confirmation/service"
	"backend/internal/response"
	"backend/pkg/logger"
)

func (h *Handler) Confirm(w http.ResponseWriter, r *http.Request) {
	tokenData, err := h.getRequestTokenData(r)
	ctx := r.Context()
	if err != nil {
		logger.WarnCtx(ctx, "Failed to get request token data", "error", err)
		response.NotFoundErrorResponse(w)
		return
	}
	confirmationService := service.NewService(ctx)

	switch tokenData.Type {
	case confirmation.ConfirmationTokenTypeRegister:
		id, err := confirmationService.ConfirmRegisterToken(ctx, *tokenData)
		if err != nil || id == 0 {
			logger.ErrorCtx(ctx, "Failed to confirm registration token", "error", err)
			response.InternalServerError(w)
			return
		}
	case confirmation.ConfirmationTokenTypeEmailChange:
		err := confirmationService.ConfirmEmailChangeToken(ctx, *tokenData)
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to confirm email change token", "error", err)
			response.InternalServerError(w)
			return
		}
	default:
		logger.WarnCtx(ctx, "Unknown confirmation token type", "type", tokenData.Type)
		response.BadRequestErrorResponse(w)
		return
	}
	err = confirmationService.ConsumeToken(ctx, tokenData.Token)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to consume confirmation token", "error", err)
		return
	}
	dto.SetConfirmSuccessResponse(w, ctx, tokenData.Type)
}
