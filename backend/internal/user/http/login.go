package http

import (
	"backend/internal/contexthelper"
	"backend/internal/response"
	sessionservice "backend/internal/session/service"
	"backend/internal/user/dto"
	"backend/internal/user/service"
	"backend/pkg/logger"

	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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
		dto.LoginErrorInvalidCredentials(w)
		return
	}

	userService := service.NewService(ctx)

	user, err := userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		logger.ErrorCtx(ctx, "Login failed", "error", err)
		dto.LoginErrorInvalidCredentials(w)
		return
	}
	logger.InfoCtx(ctx, "User logged in successfully", "userId", user.Id)
	accessTokenData, ctx := contexthelper.GetAccessTokenData(ctx)
	accessTokenData.SetCookies = true
	accessTokenData.UserId = user.Id
	logger.DebugCtx(ctx, "Login request 2: %v", req)
	if req.RememberMe {
		logger.DebugCtx(ctx, "inside if statement")
		sessionService := sessionservice.NewService(ctx, w)
		refreshToken, err := sessionService.CreateRefreshToken(ctx, user.Id, r.UserAgent())
		if err != nil {
			logger.ErrorCtx(ctx, "Create refresh token failed", "error", err)
		} else {
			logger.DebugCtx(ctx, "set refresh token as true")
			accessTokenData.RefreshToken = refreshToken
		}

	}
	logger.DebugCtx(ctx, "accessTokenData in handler: %v", accessTokenData)
	dto.SetLoginSuccessResponse(w, ctx, user)
}
