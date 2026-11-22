package dto

import (
	"backend/internal/contexthelper"
	"backend/internal/response"
	"context"
	"net/http"
)

type pingResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func PingSuccessResponse(w http.ResponseWriter, ctx context.Context, message string) {
	requestID := contexthelper.GetRequestID(ctx)
	response.SuccessDataResponse(w, ctx, pingResponse{
		Message:   message,
		RequestID: requestID,
	})
}
