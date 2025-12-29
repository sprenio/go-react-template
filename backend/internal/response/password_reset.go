package response

import (
	"net/http"
	"context"
	"backend/internal/apicodes"
)

func PasswordResetSuccessResponse(w http.ResponseWriter, ctx context.Context) {
    SuccessCodeResponse(w, ctx, apicodes.API_Password_Reset_Success)
}