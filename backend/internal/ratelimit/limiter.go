package ratelimit

import (
	"context"
	"fmt"
	"time"

	"backend/pkg/apperror"
	"backend/pkg/logger"

	"github.com/redis/go-redis/v9"
)

// Redis Lua script implementing sliding window rate limiting via sorted set.
// KEYS[1]: sorted set key
// ARGV[1]: current timestamp in milliseconds
// ARGV[2]: window size in milliseconds
// ARGV[3]: max allowed requests
// ARGV[4]: unique request member identifier
// Returns:
// [1] allowed (1 = allowed, 0 = rejected)
// [2] current count after (or before if rejected)
// [3] earliest timestamp in current window (in milliseconds)
var slidingWindowScript = redis.NewScript(`
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local member = ARGV[4]

local clearBefore = now - window
redis.call('ZREMRANGEBYSCORE', key, 0, clearBefore)

local currentCount = redis.call('ZCARD', key)
if currentCount < limit then
    redis.call('ZADD', key, now, member)
    local ttlSeconds = math.ceil(window / 1000) + 1
    redis.call('EXPIRE', key, ttlSeconds)
    return {1, currentCount + 1, 0}
else
    local oldest = redis.call('ZRANGE', key, 0, 0, 'WITHSCORES')
    local oldestScore = 0
    if #oldest >= 2 then
        oldestScore = tonumber(oldest[2])
    end
    return {0, currentCount, oldestScore}
end
`)

type Limiter interface {
	Check(ctx context.Context, keyPrefix string, rules []LimitRule) (*Result, error)
}

type RedisRateLimiter struct {
	redisClient *redis.Client
}

func NewRedisRateLimiter(redisClient *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{
		redisClient: redisClient,
	}
}

// Check tests the request against all provided rules.
// If any rule exceeds limit, allowed is false, and RetryAfter reflects the time until this counter clears.
func (r *RedisRateLimiter) Check(ctx context.Context, keyPrefix string, rules []LimitRule) (*Result, error) {
	if len(rules) == 0 {
		return &Result{Allowed: true}, nil
	}

	now := time.Now()
	nowMs := now.UnixNano() / int64(time.Millisecond)

	var worstResult *Result

	for ruleIndex, rule := range rules {
		windowMs := rule.Window.Milliseconds()
		if windowMs <= 0 {
			windowMs = 1000
		}

		key := fmt.Sprintf("ratelimit:%s:%dms", keyPrefix, windowMs)
		member := fmt.Sprintf("%d_%d", now.UnixNano(), ruleIndex)

		evalResult, err := slidingWindowScript.Run(ctx, r.redisClient, []string{key}, nowMs, windowMs, rule.Rate, member).Result()
		if err != nil {
			logger.ErrorCtx(ctx, "Redis rate limit script execution failed", "error", err, "key", key)
			return nil, apperror.Wrap(err, "Rate limiter evaluation failed", "key", key)
		}

		resultSlice, ok := evalResult.([]any)
		if !ok || len(resultSlice) < 3 {
			return nil, apperror.New("Unexpected result format from rate limit script")
		}

		allowedFlag, _ := resultSlice[0].(int64)
		count, _ := resultSlice[1].(int64)
		oldestScoreMs, _ := resultSlice[2].(int64)

		isAllowed := allowedFlag == 1
		remaining := rule.Rate - count
		if remaining < 0 {
			remaining = 0
		}

		if !isAllowed {
			retryAfterMs := (oldestScoreMs + windowMs) - nowMs
			if retryAfterMs <= 0 {
				retryAfterMs = 1000
			}
			retryAfter := time.Duration(retryAfterMs) * time.Millisecond

			if worstResult == nil || retryAfter > worstResult.RetryAfter {
				worstResult = &Result{
					Allowed:    false,
					Remaining:  0,
					RetryAfter: retryAfter,
					Limit:      rule.Rate,
					Window:     rule.Window,
				}
			}
		} else if worstResult == nil {
			// Track lowest remaining capacity among rules that passed so far
			worstResult = &Result{
				Allowed:    true,
				Remaining:  remaining,
				RetryAfter: 0,
				Limit:      rule.Rate,
				Window:     rule.Window,
			}
		}
	}

	return worstResult, nil
}
