package http

import (
	"backend/internal/confirmation"
	"backend/internal/confirmation/service"
	"backend/pkg/apperror"
	"backend/pkg/logger"
	"net/http"

	"errors"

	"github.com/go-chi/chi/v5"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) getRequestTokenData(r *http.Request) (*confirmation.ConfirmationToken, error) {
	ctx := r.Context()
	token := chi.URLParam(r, "token")
	logger.InfoCtx(ctx, "ConfirmHandler called", "token", token)
	if token == "" {
		logger.WarnCtx(ctx, "Empty confirmation token provided")
		return nil, errors.New("empty confirmation token")
	}
	confirmationService := service.NewService(ctx)
	data, err := confirmationService.GetActiveNewToken(ctx, token)
	if err != nil || data.Id == 0 {
		logger.WarnCtx(ctx, "No active confirmation token found", "token", token, "error", err)
		return nil, apperror.Wrap(err, "no active confirmation token found")
	}
	return &data, nil
}
