package middleware

import (
	"backend/config"
	"backend/internal/response"
	"backend/pkg/logger"
	"backend/internal/contexthelper"

	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)



func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			logger.ErrorCtx(ctx, "Auth header is missing or malformed")
			response.UnauthorizedErrorResponse(w, "Auth header is missing or malformed")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		logger.DebugCtx(ctx, "Received token: %s", tokenStr)
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// sprawdź typ algorytmu
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unsupported algorithm: %v", t.Header["alg"])
			}
			return config.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			logger.ErrorCtx(ctx, "Invalid JWT token: %v", err)
			response.UnauthorizedErrorResponse(w, "Invalid JWT token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.ErrorCtx(ctx, "Invalid claims in JWT token")
			response.UnauthorizedErrorResponse(w, "Invalid claims in JWT token")
			return
		}

		userIDFloat, ok := claims["user_id"].(float64) // bo JSON parsuje liczby jako float64
		if !ok {
			logger.ErrorCtx(ctx, "User ID not found in JWT claims")
			response.UnauthorizedErrorResponse(w, "Invalid JWT token")
			return
		}
		userID := uint(userIDFloat)
		logger.InfoCtx(ctx, "Authenticated user ID: %d", userID)
		
		// dodaj user_id do kontekstu
		ctx = contexthelper.SetUserId(ctx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
