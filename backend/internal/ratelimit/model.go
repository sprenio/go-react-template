package ratelimit

import "time"

// Result contains the rate limiting check outcome.
type Result struct {
	Allowed    bool
	Remaining  int64
	RetryAfter time.Duration
	Limit      int64
	Window     time.Duration
}
