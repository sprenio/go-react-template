package cache

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	DefaultTimeout  = 5 * time.Second
	DefaultInterval = 100 * time.Millisecond
)

// RedisLock struktura przechowująca stan locka
type RedisLock struct {
	client *redis.Client
	key    string
	value  string // Unikalny token identyfikujący właściciela blokady
}

func NewLock(client *redis.Client, key string) *RedisLock {
	return &RedisLock{
		client: client,
		key:    "lock:" + key,
		value:  uuid.New().String(), // Losowe UUID per blokada
	}
}

func WaitUntilLockReleased(ctx context.Context, client *redis.Client, key string, timeout time.Duration, interval time.Duration) error {
	lock := NewLock(client, key)

	return lock.WaitUntilUnlocked(ctx, timeout, interval)
}

// Acquire próbuje zdobyć blokadę na określony czas (ttl)
func (l *RedisLock) Acquire(ctx context.Context, ttl time.Duration) (bool, error) {
	return l.client.SetNX(ctx, l.key, l.value, ttl).Result()
}

// AcquireWithRetry próbuje zdobyć locka przez czas `timeout`, sprawdzając co `retryInterval`
func (l *RedisLock) AcquireWithRetry(ctx context.Context, ttl time.Duration, timeout time.Duration, interval time.Duration) (bool, error) {

	retryCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(interval)

	ok, err := l.Acquire(retryCtx, ttl)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	for {
		select {
		case <-retryCtx.Done():
			// Przekroczono czas oczekiwania (timeout) lub kontekst nadrzędny został anulowany
			return false, nil
		case <-ticker.C:
			ok, err := l.Acquire(retryCtx, ttl)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
	}
}

// Skrypt Lua zapewniający atomość zwalniania:
// Sk kasuj klucz TYLKO jeśli jego wartość zrównuje się z naszym UUID
var releaseScript = redis.NewScript(`
    if redis.call("get", KEYS[1]) == ARGV[1] then
        return redis.call("del", KEYS[1])
    else
        return 0
    end
`)

// Release bezpiecznie zwalnia blokadę
func (l *RedisLock) Release(ctx context.Context) error {
	res, err := releaseScript.Run(ctx, l.client, []string{l.key}, l.value).Int64()
	if err != nil {
		return err
	}
	if res == 0 {
		// Blokada albo wygasła i została przejęta przez inny proces, albo nigdy do nas nie należała
		return errors.New("lock stale lub wygasł - nie usunięto")
	}
	return nil
}

func (l *RedisLock) IsLocked(ctx context.Context) (bool, error) {
	exists, err := l.client.Exists(ctx, l.key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

func (l *RedisLock) WaitUntilUnlocked(
	ctx context.Context,
	timeout time.Duration,
	interval time.Duration,
) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		locked, err := l.IsLocked(ctx)
		if err != nil {
			return err
		}

		if !locked {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-timer.C:
			return errors.New("timeout waiting for queue lock")

		case <-ticker.C:
		}
	}
}
