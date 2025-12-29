package handler

import (
	"backend/internal/repository"
	"backend/internal/response"
	"backend/internal/service"
	"backend/pkg/logger"
	"backend/pkg/validation"

	"encoding/json"
	"net/http"
)

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

func (h *Handler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.InvalidJsonErrorResponse(w)
		return
	}
	logger.DebugCtx(ctx, "ResetPasswordHandler called with data: %v", req)
	if req.Email == "" || !validation.IsEmailValid(req.Email) {
		response.InvalidInputValueErrorResponse(w, "email", "invalid email format")
		return
	}

	uRepo := repository.NewUserRepository(h.db)
	ctRepo := repository.NewConfirmationTokenRepository(h.db)

	service := service.NewPasswordService(ctRepo, uRepo)
	err := service.ResetPassword(ctx, h.rabbitConn, req.Email)

	if err != nil {
		logger.ErrorCtx(ctx, "Reset password failed: %v", err)
		response.InternalServerError(w)
		return
	}

	response.PasswordResetSuccessResponse(w, r.Context())
}
