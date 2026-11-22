package dto

import (
	"backend/internal/response"
	"context"
	"net/http"
)

func SetResetPasswordSuccessResponse(w http.ResponseWriter, ctx context.Context) {
	response.SuccessNoContentResponse(w, ctx)
}
