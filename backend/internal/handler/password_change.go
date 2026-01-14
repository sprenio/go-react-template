package handler

import (
	"backend/internal/contexthelper"
	"backend/internal/repository"
	"backend/internal/response"
	"backend/internal/service"
	"backend/pkg/logger"
	"backend/pkg/validation"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type PasswordChangeRequest struct {
	Password string `json:"password"`
}

func (h *Handler) PasswordChangeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := chi.URLParam(r, "token")
	logger.DebugCtx(ctx, "ConfirmHandler called with token: %s", token)
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
		logger.DebugCtx(ctx, "No active confirmation token found", "token", token, "error", err)
		response.NotFoundErrorResponse(w)
		return
	}

	var req PasswordChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.InvalidJsonErrorResponse(w)
		return
	}
	logger.DebugCtx(ctx, "PasswordChangeHandler called with data: %v", req)

	if !validation.IsPasswordValid(req.Password) {
		response.InvalidInputValueErrorResponse(w, "password", "invalid password format")
		return
	}

	uRepo := repository.NewUserRepository(db)
	service := service.NewPasswordService(ctRepo, uRepo)
	err = service.PasswordChange(ctx, ct.UserId, req.Password)

	if err != nil {
		response.InternalServerError(w)
		return
	}
	err = ctRepo.ConsumeToken(ctx, token)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to consume confirmation token: %v", err)
		return
	}

	response.PasswordChangeSuccessResponse(w, r.Context())
}
