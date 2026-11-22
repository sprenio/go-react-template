package middleware

import (
	"net/http"
)

// EndpointLimit creates an HTTP middleware that extracts the rate limiter from context and enforces the named limit.
func EndpointLimit(endpointIdentifier string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter := GetRateLimiterFromContext(r.Context())
			if limiter != nil {
				EndpointRateLimit(limiter, endpointIdentifier)(next).ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
