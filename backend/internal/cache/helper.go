package cache

import (
	"context"
	"math"
	"time"

	"backend/pkg/apperror"
	"backend/pkg/logger"

	"github.com/redis/go-redis/v9"
)

const maxRedisIntervalSeconds = 30.0

func ConnectWithRetry(ctx context.Context, redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, apperror.Wrap(err, "Invalid Redis URL", "url", redisURL)
	}
	client := redis.NewClient(opts)

	var attempt int
	var pingErr error
	for {
		// Sprawdzamy czy Redis odpowiada
		if pingErr = client.Ping(ctx).Err(); pingErr == nil {
			logger.Info("Połączono z Redisem")
			return client, nil
		}
		attempt++
		backoff := time.Duration(math.Min(math.Pow(2, float64(attempt)), maxRedisIntervalSeconds)) * time.Second
		logger.Warn("Błąd połączenia z Redisem (próba %d): %v. \nredisURL: %s\nPonawiam za %v", attempt, pingErr, redisURL, backoff)

		select {
		case <-ctx.Done():
			client.Close()
			return nil, apperror.Wrap(ctx.Err(), "timeout podczas łączenia z Redisem", "redisURL", redisURL, "attempt", attempt)
		case <-time.After(backoff):
		}
	}
}
