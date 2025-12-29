package handler

import (
	"backend/internal/contexthelper"
	"backend/internal/response"
	"backend/pkg/logger"
	"net/http"
)

type pingResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func (h *Handler) PingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := contexthelper.GetRequestID(ctx)

	logger.InfoCtx(ctx, "ping request; id: %s", requestID)
	response.SuccessDataResponse(w, ctx,pingResponse{
		Message:   "pong strong",
		RequestID: requestID,
	})
}
