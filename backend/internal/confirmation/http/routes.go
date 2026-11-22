package http

import (
	"github.com/go-chi/chi/v5"
)

func RegisterUnauthorized(r chi.Router) {
	h := NewHandler()

	r.Get("/confirm/{token}", h.Confirm)
	r.Post("/password-change/{token}", h.ChangePassword)
}
