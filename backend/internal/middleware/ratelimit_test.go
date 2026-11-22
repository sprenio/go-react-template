package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend/internal/contexthelper"
	"backend/internal/ratelimit"
)

type mockLimiterForMiddleware struct {
	allowed    bool
	retryAfter time.Duration
}

func (m *mockLimiterForMiddleware) Check(ctx context.Context, keyPrefix string, rules []ratelimit.LimitRule) (*ratelimit.Result, error) {
	return &ratelimit.Result{
		Allowed:    m.allowed,
		RetryAfter: m.retryAfter,
	}, nil
}

func TestGlobalRateLimit_Allowed(t *testing.T) {
	mock := &mockLimiterForMiddleware{allowed: true}
	cfg := ratelimit.DefaultConfig()
	svc := ratelimit.NewService(mock, &cfg)

	handler := GlobalRateLimit(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestGlobalRateLimit_Blocked(t *testing.T) {
	mock := &mockLimiterForMiddleware{allowed: false, retryAfter: 2 * time.Second}
	cfg := ratelimit.DefaultConfig()
	svc := ratelimit.NewService(mock, &cfg)

	handler := GlobalRateLimit(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", rr.Code)
	}
	if retryAfter := rr.Header().Get("Retry-After"); retryAfter != "2" {
		t.Errorf("Expected Retry-After header '2', got '%s'", retryAfter)
	}
}

func TestEndpointLimit_Middleware(t *testing.T) {
	mock := &mockLimiterForMiddleware{allowed: false, retryAfter: 5 * time.Second}
	cfg := ratelimit.DefaultConfig()
	svc := ratelimit.NewService(mock, &cfg)

	handler := EndpointLimit("login")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/login", nil)
	ctx := contexthelper.SetRateLimiter(req.Context(), svc)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", rr.Code)
	}
}
