package database
import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"backend/pkg/logger"
)

const maxDBIntervalSeconds = 30

// connectWithRetry próbuje połączyć się z bazą danych z exponential backoff.
// Przerywa działanie, jeśli minie timeout z kontekstu.
func ConnectWithRetry(ctx context.Context, db *sql.DB) error {
	var attempt int
	var err error

	for {
		if err = db.PingContext(ctx); err == nil {
			logger.Info("Połączono z bazą danych")
			return nil
		}

		attempt++
		// backoff: min(2^attempt, 30) sekund
		backoff := time.Duration(math.Min(math.Pow(2, float64(attempt)), maxDBIntervalSeconds)) * time.Second
		logger.Warn("Błąd połączenia z bazą (próba %d): %v. Ponawiam za %v", attempt, err, backoff)

		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout podczas łączenia z bazą: %w", err)
		case <-time.After(backoff):
		}
	}
}