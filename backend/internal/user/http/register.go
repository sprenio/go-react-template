package http

import (
	"context"
	"net/http"
	"strings"

	"backend/internal/apierrors"
	"backend/internal/helper"
	"backend/internal/response"
	"backend/internal/user/dto"
	"backend/internal/user/service"
	"backend/pkg/logger"
	"backend/pkg/validation"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Language string `json:"language"`
}

func (req *RegisterRequest) Validate(ctx context.Context) error {
	// FORMATTING
	req.Email = helper.CleanEmail(req.Email)
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.Language = strings.ToLower(strings.TrimSpace(req.Language))

	// VALIDATION
	if req.Username == "" {
		return apierrors.NewInvalidInputError("Username", "username is required", req)
	}

	err := validation.IsEmailValid(req.Email)
	if err != nil {
		return apierrors.NewInvalidInputError("Email", err.Error(), req)
	}
	err = validation.IsPasswordValid(req.Password)
	if err != nil {
		return apierrors.NewInvalidInputError("Password", err.Error(), req)
	}

	return nil
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rqData, err := helper.GetRequestData[RegisterRequest](r, w)
	if err != nil {
		response.InvalidInputResponse(w, err)
		return
	}
	userService := service.NewService(ctx)
	logger.InfoCtx(ctx, "Register user", "name", rqData.Username, "email", rqData.Email, "language", rqData.Language)
	err = userService.EnqueueRegisterUser(ctx, rqData.Username, rqData.Email, rqData.Password, rqData.Language)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to register user", "error", err)
		if apierrors.IsRegisterUserNameTakenError(err) {
			logger.InfoCtx(ctx, "Username already taken", "username", rqData.Username, "error", err.Error())
			dto.RegisterErrorUserNameTaken(w)
		} else {
			response.InternalServerError(w)
		}
		return
	}
	logger.InfoCtx(ctx, "User registered successfully", "username", rqData.Username, "email", rqData.Email)

	dto.SetRegisterSuccessResponse(w, ctx)
}
