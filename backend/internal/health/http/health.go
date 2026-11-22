package http

import (
	"backend/internal/health/dto"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func Register(r chi.Router) {
	h := NewHandler()

	r.Get("/health", h.Health)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto.SetHealthResponse(w, ctx)

}
