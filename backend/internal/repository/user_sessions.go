package repository

import (
	"backend/pkg/logger"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

const (
	UserSessionsTable = "user_sessions"
)

type UserSessionsRepository struct {
	db DBExecutor
}

func NewUserSessionsRepository(db DBExecutor) *UserSessionsRepository {
	return &UserSessionsRepository{db: db}
}
func (r *UserSessionsRepository) Create(ctx context.Context, userId uint, token string, expiresAt time.Time, userAgent string, ip string) error {
	hash := hashToken(token)
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO `+UserSessionsTable+` (user_id, token_hash, created_at, expires_at, user_agent, ip)
		VALUES (?, ?, NOW(), ?, ?, ?)`,
		userId, hash, expiresAt, userAgent, ip)
	return err
}
func (r *UserSessionsRepository) GetActiveUserId(ctx context.Context, token string) (uint, error) {
	var userId uint = 0
	hash := hashToken(token)
	row := r.db.QueryRowContext(ctx, `SELECT user_id FROM `+UserSessionsTable+` WHERE token_hash = ? AND revoked_at IS NULL AND expires_at > NOW()`,
		hash,
	)
	logger.DebugCtx(ctx, "hash: %v", hash)
	err := row.Scan(&userId)
	return userId, err
}
func (r *UserSessionsRepository) RefreshToken(ctx context.Context, userId uint, token string, expiresAt time.Time) error {
	hash := hashToken(token)
	sql := `UPDATE ` + UserSessionsTable + ` SET expires_at=?, refreshed_at=NOW() WHERE user_id=? AND token_hash =? AND revoked_at IS NULL`
	result, err := r.db.ExecContext(ctx, sql, expiresAt, userId, hash)
	if err != nil {
		return err
	}
	refreshed, _ := result.RowsAffected()
	logger.InfoCtx(ctx, "Sessions refreshed: %d", refreshed)
	return nil
}
func (r *UserSessionsRepository) Revoke(ctx context.Context, userId uint, token string) error {
	hash := hashToken(token)
	sql := `UPDATE ` + UserSessionsTable + ` SET revoked_at=NOW() WHERE user_id=? AND token_hash =? AND revoked_at IS NULL`
	logger.DebugCtx(ctx, "sql: %s;userId: %v; hash: %v", sql, userId, hash)
	result, err := r.db.ExecContext(ctx, sql, userId, hash)
	if err != nil {
		return err

	}
	revoked, _ := result.RowsAffected()
	logger.InfoCtx(ctx, "Sessions revoked: %d", revoked)
	return nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:]) // 64 znaki
}
