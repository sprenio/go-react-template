package handler

import (
	"backend/internal/contexthelper"
	"backend/internal/repository"
	"backend/internal/response"
	"backend/internal/service"
	"backend/pkg/logger"

	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorCtx(ctx, "Failed to decode login request", "error", err)
		response.InvalidJsonErrorResponse(w)
		return
	}
	logger.DebugCtx(ctx, "Login request: %v", req)
	if req.Email == "" || req.Password == "" {
		logger.ErrorCtx(ctx, "Invalid login request", "error", "email or password is empty")
		response.LoginErrorInvalidCredentials(w)
		return
	}

	db := contexthelper.GetDb(ctx)
	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(userRepo)

	user, err := authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		logger.ErrorCtx(ctx, "Login failed: %v", err)
		response.LoginErrorInvalidCredentials(w)
		return
	}
	logger.InfoCtx(ctx, "User %d logged in successfully", user.Id)
	accessTokenData, ctx := contexthelper.GetAccessTokenData(ctx)
	accessTokenData.SetCookies = true
	accessTokenData.UserId = user.Id
	logger.DebugCtx(ctx, "Login request 2: %v", req)
	if req.RememberMe {
		logger.DebugCtx(ctx, "inside if statement")
		sessionRepo := repository.NewUserSessionsRepository(db)
		sessionService := service.NewSessionService(sessionRepo, w)
		refreshToken, err := sessionService.CreateRefreshToken(ctx, user.Id, r.UserAgent())
		if err != nil {
			logger.ErrorCtx(ctx, "Create refresh token failed: %v", err)
		} else {
			logger.DebugCtx(ctx, "set refresh token as true")
			accessTokenData.RefreshToken = refreshToken
		}

	}
	logger.DebugCtx(ctx, "accessTokenData in handler: %v", accessTokenData)
	response.SetLoginSuccessResponse(w, ctx, user)
}
