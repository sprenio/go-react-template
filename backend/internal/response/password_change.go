package response

import (
	"net/http"
	"context"
	"backend/internal/apicodes"
)

func PasswordChangeSuccessResponse(w http.ResponseWriter, ctx context.Context) {
    SuccessCodeResponse(w, ctx, apicodes.API_Password_Change_Success)
}