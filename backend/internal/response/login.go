package response

import (
	"net/http"
	"context"

	"backend/internal/apicodes"
)

func SetLoginSuccessResponse(w http.ResponseWriter, ctx context.Context, user UserResponseData) {
	SetMeSuccessResponse(w, ctx, user)
}

func LoginErrorInvalidCredentials(w http.ResponseWriter) {
	apiErrorResponse(w, http.StatusUnauthorized, apicodes.API_Login_Invalid_Credentials)
}
