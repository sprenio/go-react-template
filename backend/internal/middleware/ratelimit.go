package middleware

import (
	"context"
	"math"
	"net/http"

	"backend/internal/contexthelper"
	"backend/internal/ratelimit"
	"backend/internal/response"
	"backend/pkg/logger"
)

// WithRateLimiter stores the ratelimit.Service in the request context.
func WithRateLimiter(rateLimiter *ratelimit.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := contexthelper.SetRateLimiter(r.Context(), rateLimiter)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GlobalRateLimit inspects global limits for anonymous (by IP) or logged-in users (by User ID).
func GlobalRateLimit(rateLimiter *ratelimit.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			result, err := rateLimiter.CheckGlobal(ctx)
			if err != nil {
				logger.ErrorCtx(ctx, "Global rate limit check encountered error", "error", err)
				// In case of rate limiter backend error, do not block user traffic unexpectedly
				next.ServeHTTP(w, r)
				return
			}

			if !result.Allowed {
				logger.WarnCtx(ctx, "Global rate limit exceeded",
					"retryAfter", result.RetryAfter,
					"limit", result.Limit,
					"window", result.Window,
				)
				retryAfterSeconds := int(math.Ceil(result.RetryAfter.Seconds()))
				if retryAfterSeconds <= 0 {
					retryAfterSeconds = 1
				}
				response.TooManyRequestsErrorResponse(w, retryAfterSeconds)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// EndpointRateLimit enforces rate limiting on a specific endpoint identifier.
func EndpointRateLimit(rateLimiter *ratelimit.Service, endpointIdentifier string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			result, err := rateLimiter.CheckEndpoint(ctx, endpointIdentifier)
			if err != nil {
				logger.ErrorCtx(ctx, "Endpoint rate limit check encountered error", "endpoint", endpointIdentifier, "error", err)
				next.ServeHTTP(w, r)
				return
			}

			if !result.Allowed {
				logger.WarnCtx(ctx, "Endpoint rate limit exceeded",
					"endpoint", endpointIdentifier,
					"retryAfter", result.RetryAfter,
					"limit", result.Limit,
					"window", result.Window,
				)
				retryAfterSeconds := int(math.Ceil(result.RetryAfter.Seconds()))
				if retryAfterSeconds <= 0 {
					retryAfterSeconds = 1
				}
				response.TooManyRequestsErrorResponse(w, retryAfterSeconds)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetRateLimiterFromContext extracts ratelimit.Service from context.
func GetRateLimiterFromContext(ctx context.Context) *ratelimit.Service {
	val := contexthelper.GetRateLimiter(ctx)
	if limiter, ok := val.(*ratelimit.Service); ok {
		return limiter
	}
	return nil
}
