package dto

import (
	"context"
	"net/http"

	"backend/internal/response"
)

func SetHealthResponse(w http.ResponseWriter, ctx context.Context) {
	response.SuccessNoContentResponse(w, ctx)
}
