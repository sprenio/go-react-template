package cookie

import (
	"backend/internal/contexthelper"
	"backend/pkg/logger"
	"fmt"

	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func getCookieValue(r *http.Request, cookieKey string) string {
	c, err := r.Cookie(cookieKey)
	if err != nil {
		return ""
	}
	return c.Value
}
func deleteCookie(ctx context.Context, w http.ResponseWriter, cookieKey string) {
	logger.DebugCtx(ctx, "Delete cookie %v", cookieKey)
	cfg := contexthelper.GetConfig(ctx)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieKey,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   !cfg.IsDevEnv(),
		SameSite: http.SameSiteLaxMode,
	})
}
func GetUserIdFromJwtToken(r *http.Request) (uint, error) {
	ctx := r.Context()
	cookieToken := GetAccessToken(r)
	logger.DebugCtx(ctx, "Received token: %s", cookieToken)
	if cookieToken == "" {
		err := fmt.Errorf("Access token is missing")
		return 0, err
	}
	cfg := contexthelper.GetConfig(ctx)

	token, err := jwt.Parse(cookieToken, func(t *jwt.Token) (interface{}, error) {
		// sprawdź typ algorytmu
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unsupported algorithm: %v", t.Header["alg"])
		}
		return []byte(cfg.Token.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		if err == nil {
			err = fmt.Errorf("Invalid JWT token")
		}
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err := fmt.Errorf("Invalid claims in JWT token")
		return 0, err
	}
	userIDFloat, ok := claims["user_id"].(float64) // bo JSON parsuje liczby jako float64
	if !ok {
		err := fmt.Errorf("User ID not found in JWT claims")
		return 0, err
	}
	return uint(userIDFloat), nil
}
func GetAccessToken(r *http.Request) string {
	return getCookieValue(r, AccessTokenKey)
}
func GetRefreshToken(r *http.Request) string {
	return getCookieValue(r, RefreshTokenKey)
}
func DeleteAccessToken(ctx context.Context, w http.ResponseWriter) {
	deleteCookie(ctx, w, AccessTokenKey)
}
func DeleteRefreshToken(ctx context.Context, w http.ResponseWriter) {
	deleteCookie(ctx, w, RefreshTokenKey)
}

func SetAccessToken(ctx context.Context, w http.ResponseWriter, userId uint) {
	token, err := generateJWT(ctx, userId)
	if err != nil {
		logger.WarnCtx(ctx, "Failed to generate JWT token: %v", err)
		return
	}
	cfg := contexthelper.GetConfig(ctx)
	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenKey,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   !cfg.IsDevEnv(), // true na HTTPS
		SameSite: http.SameSiteLaxMode,
		//MaxAge:   300, // 5 minut
	})
	logger.DebugCtx(ctx, "Set access token cookie value", token)
}
func SetRefreshToken(w http.ResponseWriter, ctx context.Context, token string, maxAge int) {
	cfg := contexthelper.GetConfig(ctx)
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenKey,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   !cfg.IsDevEnv(), // true na HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge, // 5 minut
	})
	logger.DebugCtx(ctx, "Set refresh token cookie value", token)
}

func generateJWT(ctx context.Context, userID uint) (string, error) {
	cfg := contexthelper.GetConfig(ctx)
	ttl := time.Minute * time.Duration(int64(cfg.Token.AccessTokenTtlMinutes))
	logger.DebugCtx(ctx, "minutes:", cfg.Token.AccessTokenTtlMinutes, "ttl:", ttl, "expires:", jwt.NewNumericDate(time.Now().Add(ttl)))
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Token.JwtSecret))
}
