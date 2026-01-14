package response

import (
	"context"
	"net/http"

	"backend/internal/apicodes"
	"backend/internal/models"
)

func SetLoginSuccessResponse(w http.ResponseWriter, ctx context.Context, user models.UserResponseData) {
	SetMeSuccessResponse(w, ctx, user)
}

func LoginErrorInvalidCredentials(w http.ResponseWriter) {
	apiErrorResponse(w, http.StatusUnauthorized, apicodes.API_Login_Invalid_Credentials)
}
