package middleware

import (
	"net/http"
	"time"
)

func Sleep(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		next.ServeHTTP(w, r)
	})
}
