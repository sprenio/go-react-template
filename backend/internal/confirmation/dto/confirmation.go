package dto

import (
	"backend/internal/apicodes"
	"backend/internal/response"
	"context"
	"net/http"
)

func ResetPasswordSuccessResponse(w http.ResponseWriter, ctx context.Context) {
	response.SuccessNoContentResponse(w, ctx)
}
func ChangePasswordSuccessResponse(w http.ResponseWriter, ctx context.Context) {
	response.SuccessNoContentResponse(w, ctx)
}
func SetConfirmSuccessResponse(w http.ResponseWriter, ctx context.Context, tokenType string) {
	data := map[string]string{
		"token_type": tokenType,
	}
	response.SuccessDataCodeResponse(w, ctx, data, apicodes.API_Confirm_Success)
}
func EmailChangeSuccessResponse(w http.ResponseWriter, ctx context.Context) {
	response.SuccessCodeResponse(w, ctx, apicodes.API_Email_Change_Success)
}
