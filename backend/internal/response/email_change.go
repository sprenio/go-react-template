package response

import (
	"net/http"
	"context"
	"backend/internal/apicodes"
)

func EmailChangeSuccessResponse(w http.ResponseWriter, ctx context.Context) {
    SuccessCodeResponse(w, ctx, apicodes.API_Email_Change_Success)
}