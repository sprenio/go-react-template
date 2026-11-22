package http

import (
	"backend/internal/ping/dto"
	"backend/pkg/logger"
	"net/http"
)

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.InfoCtx(ctx, "ping request")
	dto.PingSuccessResponse(w, ctx, "pong")
}
