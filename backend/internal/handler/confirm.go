package handler

import (
	"net/http"

	"backend/internal/contexthelper"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/response"
	"backend/internal/service"
	"backend/pkg/logger"
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)


func (h *Handler) ConfirmHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := chi.URLParam(r, "token")
	logger.InfoCtx(ctx, "ConfirmHandler called with token: %s", token)
	if token == "" {
		logger.WarnCtx(ctx, "Empty confirmation token provided")
		response.NotFoundErrorResponse(w)
		return
	}
	db := contexthelper.GetDb(ctx)

	ctRepo := repository.NewConfirmationTokenRepository(db)
	if ctRepo == nil {
		response.InternalServerError(w)
		return
	}
	ct, err := ctRepo.GetActiveNewToken(ctx, token)
	if err != nil || ct.Id == 0 {
		logger.WarnCtx(ctx, "No active confirmation token found", "token", token, "error", err)
		response.NotFoundErrorResponse(w)
		return
	}
	switch ct.Type {
	case models.ConfirmationTokenTypeRegister:
		err := h.confirmRegisterHandler(ctx, ct)
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to confirm registration token: %v", err)
			response.InternalServerError(w)
			return
		}
	case models.ConfirmationTokenTypeEmailChange:
		err := h.confirmEmailChangeHandler(ctx, ct)
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to confirm email change token: %v", err)
			response.InternalServerError(w)
			return
		}
	default:
		logger.WarnCtx(ctx, "Unknown confirmation token type: %s", ct.Type)
		response.BadRequestErrorResponse(w)
		return
	}
	err = ctRepo.ConsumeToken(ctx, token)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to consume confirmation token: %v", err)
		return
	}
	response.SetConfirmSuccessResponse(w, ctx, ct.Type)
}

func (h *Handler) confirmRegisterHandler(ctx context.Context, ct models.ConfirmationToken) error {
	db := contexthelper.GetDb(ctx)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	uRepo := repository.NewUserRepository(tx)
	usRepo := repository.NewUserSettingsRepository(tx)
	service := service.NewRegisterConfirmationService(uRepo, usRepo)
	id, err := service.ConfirmRegisterToken(ctx, ct)
	if err != nil || id == 0 {
		if rbErr := tx.Rollback(); rbErr != nil {
			logger.ErrorCtx(ctx, "Failed to rollback transaction: %v", rbErr)
		}
		return errors.Wrap(err, "failed to confirm registration token")
	}
	if err := tx.Commit(); err != nil {
		logger.ErrorCtx(ctx, "Failed to commit transaction: %v", err)
		if rbErr := tx.Rollback(); rbErr != nil {
			logger.ErrorCtx(ctx, "Failed to rollback transaction: %v", rbErr)
		}
		return errors.Wrap(err, "failed to commit transaction")
	}
	return nil
}
func (h *Handler) confirmEmailChangeHandler(ctx context.Context, ct models.ConfirmationToken) error {
	db := contexthelper.GetDb(ctx)
	uRepo := repository.NewUserRepository(db)
	service := service.NewUserService(uRepo)
	err := service.ConfirmEmailChangeToken(ctx, ct)
	if err != nil {
		return errors.Wrap(err, "failed to confirm email change token")
	}
	return nil
}
func (h *Handler) passwordChangeHandler(ctx context.Context, ct models.ConfirmationToken, newPassword string) error {
	db := contexthelper.GetDb(ctx)
	uRepo := repository.NewUserRepository(db)
	ctRepo := repository.NewConfirmationTokenRepository(db)
	service := service.NewPasswordService(ctRepo, uRepo)
	err := service.PasswordChange(ctx, ct.UserId, newPassword)
	if err != nil {
		return errors.Wrap(err, "failed to confirm password change token")
	}
	return nil
}
