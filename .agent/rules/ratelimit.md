---
trigger: model_decision
description: Use this rule when working with API rate limiting, protecting endpoints from abuse, adding rate limits to routes, or implementing custom business counters in Go backend.
---

# Rate Limiting in Backend API

The application has a built-in rate limiter powered by Redis using the Sliding Window Log algorithm (atomically executed via Lua script).

## 1. Architecture Overview
- **Core Limiter (`internal/ratelimit`)**: Evaluates requests using Redis sorted sets (ZSET) per configured time window. Calculates accurate `RetryAfter` durations.
- **Global Rate Limiting (`middleware.GlobalRateLimit`)**: Automatically runs for all requests after `JWTAuth`.
  - For unauthenticated users: tracked by client IP (`ip:<client_ip>:global:<window>`).
  - For authenticated users: tracked by User ID (`user:<user_id>:global:<window>`).
- **Endpoint-specific Rate Limiting (`middleware.EndpointLimit`)**: Middleware applied directly to routes requiring specific thresholds.
- **Custom / Business Logic Rate Limiting (`Service.CheckCustom` / `Service.CheckCustomSingle`)**: Allows ad-hoc or domain-level rate limiting in Go services.

## 2. Configuration (`config/config.yaml`)

Rate limits are configured under the `rate_limit` section. Time windows use Go duration format (e.g., `1s`, `1m`, `10m`, `1h`).
Endpoint limits are defined in a single unified `endpoints` map – tracking is done automatically per User ID (for logged-in users) or per Client IP (for anonymous visitors).

```yaml
rate_limit:
  enabled: true
  anonymous_global:
    - rate: 10
      window: 1s
    - rate: 100
      window: 1m
    - rate: 800
      window: 10m
  authorized_global:
    - rate: 15
      window: 1s
    - rate: 150
      window: 1m
    - rate: 1200
      window: 10m
  endpoints:
    register:
      - rate: 1
        window: 1s
      - rate: 3
        window: 1m
      - rate: 10
        window: 10m
    login:
      - rate: 1
        window: 1s
      - rate: 5
        window: 1m
      - rate: 20
        window: 10m
    my_sensitive_endpoint:
      - rate: 3
        window: 1s
      - rate: 100
        window: 1m
  custom_limits:
    export_report:
      - rate: 2
        window: 1m
```

## 3. How to Apply to Routes

Use `middleware.EndpointLimit("configured_identifier")` in route registration:

```go
package http

import (
    "backend/internal/middleware"
    "github.com/go-chi/chi/v5"
)

func Register(r chi.Router) {
    h := NewHandler()

    // Protected endpoint with specific rate limits:
    r.With(middleware.EndpointLimit("my_sensitive_endpoint")).Post("/expensive-action", h.Create)
}
```

The middleware automatically evaluates limits against the `endpoints` configuration. If the user is authenticated, tracking uses User ID (`user:<id>:endpoint:<name>`), otherwise it uses Client IP (`ip:<ip>:endpoint:<name>`).

## 4. How to Use in Business Logic (Custom Limits)

When rate limiting must happen inside domain services or handlers (e.g. per email, per token, per resource ID):

### Option A: Using Configured Limits (`CheckCustom`)
```go
limiter := middleware.GetRateLimiterFromContext(ctx)
if limiter != nil {
    subjectIdentifier := strconv.Itoa(int(userId)) // or email / resource ID
    result, err := limiter.CheckCustom(ctx, "export_report", subjectIdentifier, nil)
    if err == nil && !result.Allowed {
        // Return 429 or domain error with retry information
        response.TooManyRequestsErrorResponse(w, int(math.Ceil(result.RetryAfter.Seconds())))
        return
    }
}
```

### Option B: Using Ad-hoc Code Limits (`CheckCustomSingle`)
```go
limiter := middleware.GetRateLimiterFromContext(ctx)
if limiter != nil {
    // E.g. limit to 5 attempts per 10 minutes for a specific target
    result, err := limiter.CheckCustomSingle(ctx, "sms_code_resend", targetPhone, 5, 10*time.Minute)
    if err == nil && !result.Allowed {
        return apperror.New("Rate limit reached")
    }
}
```

## 5. HTTP Response on Limit Exceeded
When a limit is reached:
- HTTP Status: `429 Too Many Requests`
- Headers: `Retry-After: <seconds>`
- Error body:
```json
{
  "code": 1004,
  "message": "Too many requests"
}
```
