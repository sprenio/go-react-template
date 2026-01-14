package response

import (
	"net/http"
	"context"

	"backend/internal/apicodes"
)

func SetLogoutSuccessResponse(w http.ResponseWriter, ctx context.Context) {
	SuccessCodeResponse(w, ctx,apicodes.API_Logout_Success)
}

func LogoutFailedResponse(w http.ResponseWriter) {
	apiErrorResponse(w, http.StatusBadRequest, apicodes.API_Logout_Failed)
}
