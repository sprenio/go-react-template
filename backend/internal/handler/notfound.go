package handler

import (
	"backend/internal/response"
	"backend/pkg/logger"
	"net/http"
)

func (h *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	logger.ErrorCtx(r.Context(), "404 - Not Found: %v", r.URL.Path)
	response.NotFoundErrorResponse(w)
}
