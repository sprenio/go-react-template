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
	db := contexthelper.GetDb(ctx)
	uRepo := repository.NewUserRepository(db)
	ctRepo := repository.NewConfirmationTokenRepository(db)
	langRepo := repository.NewLanguageRepository(db)

	service := service.NewEmailService(ctRepo, uRepo, langRepo)
	err := service.ChangeEmail(ctx, req.Email)

	if err != nil {
		response.InternalServerError(w)
		return
	}

	response.EmailChangeSuccessResponse(w, r.Context())
}