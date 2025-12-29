package handler

import (
	"backend/internal/response"
	"backend/pkg/logger"
	"net/http"
)

func (h *Handler) MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	logger.ErrorCtx(r.Context(), "405 - Method Not Allowed: path %v, method %v", r.URL.Path, r.Method)
	response.MethodNotAllowedErrorResponse(w)
}
