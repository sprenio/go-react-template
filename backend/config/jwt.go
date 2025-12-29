package config

import (
	"os"
	"time"
)

var (
	// JWTSecret to globalny klucz do podpisywania JWT
	JWTSecret = []byte(getEnv("JWT_SECRET", "supersekretnyklucz"))

	// TokenTTL defines how long tokens are valid (e.g. 15 min)
	TokenTTL = time.Minute * 15

	// RefreshTokenTTL defines how long refresh tokens are valid (e.g. 30 dni)
	RefreshTokenTTL = time.Hour * 24 * 30
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
