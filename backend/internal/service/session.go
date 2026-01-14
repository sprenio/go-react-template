package service

import (
	"backend/internal/contexthelper"
	"backend/internal/cookie"
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

type SessionService struct {
	sessionRepo *repository.UserSessionsRepository
	respWritter http.ResponseWriter
}

func NewSessionService(sRepo *repository.UserSessionsRepository, w http.ResponseWriter) *SessionService {
	return &SessionService{sessionRepo: sRepo, respWritter: w}
}

func (s *SessionService) CreateRefreshToken(ctx context.Context, userId uint, userAgent string) (string, error) {
	refreshToken, err := generateRefreshToken()
	if err != nil {
		logger.ErrorCtx(ctx, "Generate refresh token failed: %v", err)
		return "",err
	}
	ttlSeconds, expiresAt  := getRefreshTokenTtlData(ctx)
	logger.DebugCtx(ctx, "ttl: %d, expires: %v", ttlSeconds, expiresAt)
	ip := contexthelper.GetClientIp(ctx)
	err = s.sessionRepo.Create(ctx, userId, refreshToken, expiresAt, userAgent, ip)
	if err != nil {
		logger.ErrorCtx(ctx, "Create refresh token failed: %v", err)
		return "", err
	}
	return refreshToken, nil
}
func (s *SessionService) RefreshToken(ctx context.Context, token string) error {
	if token == "" {
		logger.DebugCtx(ctx, "Empty refresh token")
		return nil
	}
	userId, ok := contexthelper.GetUserId(ctx)
	if !ok {
		err := fmt.Errorf("No user ID in context")
		logger.ErrorCtx(ctx, "Get user id failed: %v", err)
		return err
	}
	ttlSeconds, expiresAt  := getRefreshTokenTtlData(ctx)
	logger.DebugCtx(ctx, "ttl: %d, expires: %v", ttlSeconds, expiresAt)
	err := s.sessionRepo.RefreshToken(ctx, userId, token, expiresAt)
	if err != nil {
		logger.ErrorCtx(ctx, "Refresh token failed: %v", err)
		cookie.DeleteRefreshToken(ctx, s.respWritter)
		return err
	}
	logger.DebugCtx(ctx, "ping")
	cookie.SetRefreshToken(s.respWritter, ctx, token, ttlSeconds)
	return nil
}

func (s *SessionService) Logout(ctx context.Context, token string)error {
	userId, ok := contexthelper.GetUserId(ctx)
	if !ok {
		return fmt.Errorf("User is not authenticated")
	}
	cookie.DeleteAccessToken(ctx, s.respWritter)
	if token == "" {
		return nil
	}
	cookie.DeleteRefreshToken(ctx, s.respWritter)
	return s.sessionRepo.Revoke(ctx, userId, token)
}
func (s *SessionService) Login(ctx context.Context, token string) (uint, error) {
	return s.sessionRepo.GetActiveUserId(ctx, token)
}
func generateRefreshToken() (string, error) {
    b := make([]byte, 32) // 256 bitów
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return base64.RawURLEncoding.EncodeToString(b), nil
}
func getRefreshTokenTtlData(ctx context.Context) (int, time.Time){
	cfg := contexthelper.GetConfig(ctx)
	ttl := time.Hour * 24 * time.Duration(int64(cfg.Token.RefreshTokenTtlDays))
	expiresAt := time.Now().Add(ttl);
	ttlSeconds := int(ttl / time.Second)
	return ttlSeconds, expiresAt
}