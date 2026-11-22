package http

import (
	"backend/internal/response"
	"backend/pkg/logger"
	"net/http"
)

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	logger.ErrorCtx(r.Context(), "404 - Not Found", "path", r.URL.Path)
	response.NotFoundErrorResponse(w)
}

func (h *Handler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	logger.ErrorCtx(r.Context(), "405 - Method Not Allowed", "path", r.URL.Path, "method", r.Method)
	response.MethodNotAllowedErrorResponse(w)
}
