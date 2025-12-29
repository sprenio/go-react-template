package handler

import (
	"backend/internal/response"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/logger"
	"backend/pkg/validation"
	"encoding/json"
	"net/http"
)

type EmailChangeRequest struct {
	Email string `json:"email"`
}


func(h *Handler) EmailChangeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req EmailChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.InvalidJsonErrorResponse(w)
		return
	}
	logger.DebugCtx(ctx, "EmailChangeHandler called with data: %v", req)

	if req.Email == "" || !validation.IsEmailValid(req.Email) {
		response.InvalidInputValueErrorResponse(w, "email", "invalid email format")
		return
	}

	uRepo := repository.NewUserRepository(h.db)
	ctRepo := repository.NewConfirmationTokenRepository(h.db)
	langRepo := repository.NewLanguageRepository(h.db)

	service := service.NewEmailService(ctRepo, uRepo, langRepo)
	err := service.ChangeEmail(ctx, h.rabbitConn, req.Email)

	if err != nil {
		response.InternalServerError(w)
		return
	}

	response.EmailChangeSuccessResponse(w, r.Context())
}