package middleware

import (
	"backend/internal/contexthelper"
	"database/sql"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func WithServices(db *sql.DB, rabbitConn *amqp.Connection, redisClient *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := contexthelper.SetServices(r.Context(), db, rabbitConn, redisClient)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
