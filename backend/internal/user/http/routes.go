package http

import (
	"backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterUnauthorized(r chi.Router) {
	h := NewHandler()

	r.With(middleware.EndpointLimit("login")).Post("/login", h.Login)
	r.With(middleware.EndpointLimit("register")).Post("/register", h.Register)
	r.Get("/logout", h.Logout)
	r.With(middleware.EndpointLimit("reset")).Post("/reset-password", h.ResetPassword)
}

func RegisterAuthorized(r chi.Router) {
	h := NewHandler()

	r.Get("/me", h.Me)
	r.Post("/email_change", h.ChangeEmail)
}
