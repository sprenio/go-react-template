package response

import (
	"net/http"
	"context"

	"backend/internal/apicodes"
)

func SetConfirmSuccessResponse(w http.ResponseWriter, ctx context.Context, tokenType string) {
	data := map[string]string{
		"token_type": tokenType,
	}
	SuccessDataCodeResponse(w, ctx, data, apicodes.API_Confirm_Success)
}

func SetConfirmFailureResponse(w http.ResponseWriter) {
	apiErrorResponse(w, http.StatusBadRequest, apicodes.API_Confirm_Failure)
}

func SetConfirmInvalidTokenResponse(w http.ResponseWriter) {
	apiErrorResponse(w, http.StatusBadRequest, apicodes.API_Confirm_Invalid_Token)
}
