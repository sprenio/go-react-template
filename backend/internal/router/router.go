package router

import (
	"net/http"

	"backend/internal/handler"
	"backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupRouter(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	// Rejestracja middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.AccessLog)
	r.Use(middleware.Recoverer)

	// 🔹 Publiczne endpointy (bez autoryzacji)

	r.Post("/login", h.LoginHandler)
	r.Get("/cfg", h.CfgHandler)
	r.Post("/register", h.RegisterHandler)
	r.Post("/reset-password", h.ResetPasswordHandler)
	r.Post("/password-change/{token}", h.PasswordChangeHandler)
	r.Get("/confirm/{token}", h.ConfirmHandler)

	// 🔒 Endpointy wymagające JWT
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware) // tylko dla tych poniżej
		r.Get("/ping", h.PingHandler)
		// tu możesz dodać inne chronione ścieżki
		r.Get("/me", h.MeHandler)
		r.Post("/settings", h.SettingsHandler)
		r.Post("/email_change", h.EmailChangeHandler)
	})

	// 🔹 Obsługa 404 i 405
	r.NotFound(h.NotFoundHandler)
	r.MethodNotAllowed(h.MethodNotAllowedHandler)
	return r
}
