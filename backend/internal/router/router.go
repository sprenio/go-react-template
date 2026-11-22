package router

import (
	"database/sql"
	"net/http"

	"backend/config"
	"backend/internal/middleware"
	"backend/internal/ratelimit"

	cfgHttp "backend/internal/cfg/http"
	confirmationHttp "backend/internal/confirmation/http"
	healthHttp "backend/internal/health/http"
	pingHttp "backend/internal/ping/http"
	settingsHttp "backend/internal/settings/http"
	userHttp "backend/internal/user/http"
	versionHttp "backend/internal/version/http"

	errorHttp "backend/internal/error/http"

	"github.com/go-chi/chi/v5"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func SetupRouter(cfg *config.Config, db *sql.DB, rabbitConn *amqp.Connection, redisClient *redis.Client) http.Handler {
	r := chi.NewRouter()

	// Rejestracja middleware
	r.Use(middleware.IP)
	r.Use(middleware.Config(cfg))
	r.Use(middleware.WithServices(db, rabbitConn, redisClient))
	r.Use(middleware.RequestID)
	r.Use(middleware.Version)
	r.Use(middleware.AccessLog)
	r.Use(middleware.Recoverer)

	rateLimiterConfig := ratelimit.FromAppConfig(cfg.RateLimit)
	rateLimiterCore := ratelimit.NewRedisRateLimiter(redisClient)
	rateLimiterService := ratelimit.NewService(rateLimiterCore, &rateLimiterConfig)

	r.Use(middleware.WithRateLimiter(rateLimiterService))
	r.Use(middleware.JWTAuth)
	r.Use(middleware.GlobalRateLimit(rateLimiterService))

	// 🔹 Publiczne endpointy (bez autoryzacji)
	healthHttp.Register(r)

	userHttp.RegisterUnauthorized(r)
	confirmationHttp.RegisterUnauthorized(r)
	cfgHttp.Register(r)
	versionHttp.Register(r)

	r.Group(func(r chi.Router) {
		r.Use(middleware.RefreshSession)
		r.Use(middleware.AuthOnly)

		userHttp.RegisterAuthorized(r)
		pingHttp.Register(r)
		settingsHttp.Register(r)
	})

	// 🔹 Obsługa 404 i 405
	errorHttp.Register(r)

	return r
}
