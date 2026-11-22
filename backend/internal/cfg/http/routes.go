package http

import (
	"github.com/go-chi/chi/v5"
)

func Register(r chi.Router) {
	h := NewHandler()

	r.Get("/cfg", h.Get)
}
