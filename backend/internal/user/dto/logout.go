package dto

import (
	"context"
	"net/http"

	"backend/internal/apicodes"
	"backend/internal/response"
)

func SetLogoutSuccessResponse(w http.ResponseWriter, ctx context.Context) {
	response.SuccessCodeResponse(w, ctx, apicodes.API_Logout_Success)
}

func LogoutFailedResponse(w http.ResponseWriter) {
	response.ErrorCodeResponse(w, http.StatusBadRequest, apicodes.API_Logout_Failed)
}
