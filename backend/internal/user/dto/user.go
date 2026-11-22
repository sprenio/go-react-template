package dto

import (
	"backend/internal/apicodes"
	"backend/internal/response"
	"backend/internal/user"
	"context"
	"net/http"
)

type meResponseData struct {
	User user.UserResponseData `json:"user"`
}

func SetMeSuccessResponse(w http.ResponseWriter, ctx context.Context, userData user.UserResponseData) {
	resp := meResponseData{
		User: userData,
	}
	response.SuccessDataResponse(w, ctx, resp)
}

func SetLoginSuccessResponse(w http.ResponseWriter, ctx context.Context, user user.UserResponseData) {
	SetMeSuccessResponse(w, ctx, user)
}

func LoginErrorInvalidCredentials(w http.ResponseWriter) {
	response.ErrorCodeResponse(w, http.StatusUnauthorized, apicodes.API_Login_Invalid_Credentials)
}
