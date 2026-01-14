package middleware

import (
	"backend/internal/contexthelper"
	"database/sql"
	"net/http"
	amqp "github.com/rabbitmq/amqp091-go"
)
func WithServices(db *sql.DB, rabbitConn *amqp.Connection) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := contexthelper.SetServices(r.Context(), db, rabbitConn)
            next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
