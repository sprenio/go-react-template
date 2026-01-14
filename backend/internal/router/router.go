package router

import (
	"database/sql"
	"net/http"

	"backend/config"
	"backend/internal/handler"
	"backend/internal/middleware"

	"github.com/go-chi/chi/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SetupRouter(h *handler.Handler, cfg *config.Config, db *sql.DB, rabbitConn *amqp.Connection) http.Handler {
	r := chi.NewRouter()

	// Rejestracja middleware
	r.Use(middleware.IP)
	r.Use(middleware.Config(cfg))
	r.Use(middleware.WithServices(db, rabbitConn))
	r.Use(middleware.RequestID)
	r.Use(middleware.AccessLog)
	r.Use(middleware.Recoverer)

	r.Use(middleware.JWTAuth)

	// 🔹 Publiczne endpointy (bez autoryzacji)

	r.Post("/login", h.LoginHandler)
	r.Get("/cfg", h.CfgHandler)
	r.Post("/register", h.RegisterHandler)
	r.Post("/reset-password", h.ResetPasswordHandler)
	r.Post("/password-change/{token}", h.PasswordChangeHandler)
	r.Get("/confirm/{token}", h.ConfirmHandler)
	r.Get("/logout", h.LogoutHandler)


	r.Group(func(r chi.Router) {
		r.Use(middleware.RefreshSession)
		r.Use(middleware.AuthOnly)
		r.Get("/me", h.MeHandler)
		r.Get("/ping", h.PingHandler)
		// tu możesz dodać inne chronione ścieżki
		//r.Get("/me", h.MeHandler)
		r.Post("/settings", h.SettingsHandler)
		r.Post("/email_change", h.EmailChangeHandler)
	})

	// 🔹 Obsługa 404 i 405
	r.NotFound(h.NotFoundHandler)
	r.MethodNotAllowed(h.MethodNotAllowedHandler)
	return r
}
