package handler

import (
	"backend/internal/repository"
	"backend/internal/response"
	"backend/internal/service"
	"backend/pkg/logger"
	"backend/internal/contexthelper"

	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorCtx(ctx, "Failed to decode login request", "error", err)
		response.InvalidJsonErrorResponse(w)
		return
	}

	if req.Email == "" || req.Password == "" {
		logger.ErrorCtx(ctx, "Invalid login request", "error", "email or password is empty")
		response.LoginErrorInvalidCredentials(w)
		return
	}

	userRepo := repository.NewUserRepository(h.db)
	authService := service.NewAuthService(userRepo)

	user, err := authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		logger.ErrorCtx(ctx, "Login failed: %v", err)
		response.LoginErrorInvalidCredentials(w)
		return
	}
	logger.InfoCtx(ctx, "User %d logged in successfully", user.Id)
	ctx = contexthelper.SetUserId(ctx, user.Id)
	response.SetLoginSuccessResponse(w, ctx, user)
}
