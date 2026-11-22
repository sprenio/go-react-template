package repository

import (
	"backend/internal/confirmation"
	"backend/internal/database"
	"backend/pkg/apperror"
	"backend/pkg/logger"
	"backend/pkg/uuidstr"
	"context"
	"database/sql"
	"fmt"
)

const (
	ConfirmationTokenTable = "confirmation_tokens"

	confirmationTokenColumns = `id, token, user_id, type, payload, status, expires_at, status_changed_at, created_at`

	getActiveNewTokenSql = `SELECT ` + confirmationTokenColumns + ` FROM ` + ConfirmationTokenTable + ` WHERE token = ? AND status = "` + confirmation.ConfirmationTokenStatusNew + `" AND expires_at > NOW()`

	// Check for active tokens with status 'new' or 'processing'
	activeStatusWhereCondition = `status IN ("` + confirmation.ConfirmationTokenStatusNew + `","` + confirmation.ConfirmationTokenStatusProcessing + `")`
	checkRegisterEmailSql      = `SELECT 1 FROM ` + ConfirmationTokenTable + ` WHERE type ='` + confirmation.ConfirmationTokenTypeRegister + `' AND JSON_UNQUOTE(JSON_EXTRACT(payload, '$.email')) = ? AND ` + activeStatusWhereCondition
	checkNewEmailSql           = `SELECT 1 FROM ` + ConfirmationTokenTable + ` WHERE type ='` + confirmation.ConfirmationTokenTypeEmailChange + `' AND JSON_UNQUOTE(JSON_EXTRACT(payload, '$.new_email')) = ? AND ` + activeStatusWhereCondition
	createTokenSqlPattern      = `INSERT INTO ` + ConfirmationTokenTable + ` (token, user_id, type, payload, status, expires_at, status_changed_at) VALUES (?, ?, '%s', ?, "` + confirmation.ConfirmationTokenStatusNew + `", DATE_ADD(NOW(), INTERVAL %d DAY), NOW())`
)

type ConfirmationTokenRepository struct {
	db database.DBExecutor
}

func NewRepository(db database.DBExecutor) *ConfirmationTokenRepository {
	return &ConfirmationTokenRepository{db: db}
}

func (r *ConfirmationTokenRepository) GetActiveNewToken(ctx context.Context, token string) (confirmation.ConfirmationToken, error) {
	var ct confirmation.ConfirmationToken

	err := r.db.QueryRowContext(ctx, getActiveNewTokenSql, token).Scan(&ct.Id, &ct.Token, &ct.UserId, &ct.Type, &ct.Payload, &ct.Status, &ct.ExpiresAt, &ct.StatusChangedAt, &ct.CreatedAt)
	if err != nil || ct.Id == 0 {
		return confirmation.ConfirmationToken{}, apperror.Wrap(err, "Failed to retrieve active confirmation token: "+token)
	}
	return ct, nil
}

func (r *ConfirmationTokenRepository) ExistsRegisterTokenByName(ctx context.Context, username string) (bool, error) {
	var exists bool
	row := r.db.QueryRowContext(ctx, `SELECT 1 FROM `+ConfirmationTokenTable+` WHERE type = '`+confirmation.ConfirmationTokenTypeRegister+`' AND LOWER(JSON_UNQUOTE(JSON_EXTRACT(payload, '$.name'))) = LOWER(?) AND `+activeStatusWhereCondition+` LIMIT 1`,
		username,
	)
	err := row.Scan(&exists)
	if err != sql.ErrNoRows && err != nil {
		return false, err
	}
	return exists, nil
}

func (r *ConfirmationTokenRepository) CreateRegisterToken(ctx context.Context, payload []byte, days int) (string, error) {
	return r.createToken(ctx, 0, confirmation.ConfirmationTokenTypeRegister, payload, days)
}

func (r *ConfirmationTokenRepository) CreateEmailChangeToken(ctx context.Context, userId uint, payload []byte, days int) (string, error) {
	return r.createToken(ctx, userId, confirmation.ConfirmationTokenTypeEmailChange, payload, days)
}
func (r *ConfirmationTokenRepository) CreateResetPasswordToken(ctx context.Context, userId uint, days int) (string, error) {
	return r.createToken(ctx, userId, confirmation.ConfirmationTokenTypeResetPassword, []byte("{}"), days)
}

func (r *ConfirmationTokenRepository) ConsumeToken(ctx context.Context, token string) error {
	return r.updateStatus(ctx, token, confirmation.ConfirmationTokenStatusConsumed)
}

func (r *ConfirmationTokenRepository) GetActiveNewTokenWithType(ctx context.Context, token string, tokenType string) (confirmation.ConfirmationToken, error) {
	var ct confirmation.ConfirmationToken
	err := r.db.QueryRowContext(ctx, getActiveNewTokenSql+` AND type = ?`, token, tokenType).Scan(&ct.Id, &ct.Token, &ct.UserId, &ct.Type, &ct.Payload, &ct.Status, &ct.ExpiresAt, &ct.StatusChangedAt, &ct.CreatedAt)
	if err != nil {
		return confirmation.ConfirmationToken{}, err
	}
	return ct, nil
}

func (r *ConfirmationTokenRepository) updateStatus(ctx context.Context, token string, status string) error {
	sql := `UPDATE ` + ConfirmationTokenTable + ` SET status = ?, status_changed_at = NOW() WHERE token = ? AND status != ?`
	_, err := r.db.ExecContext(ctx, sql, status, token, status)
	if err != nil {
		return apperror.Wrap(err, "Failed to update confirmation token status")
	}
	return nil
}

func (r *ConfirmationTokenRepository) IsNewEmailAvailable(ctx context.Context, email string) (bool, error) {
	var exists bool
	logger.DebugCtx(ctx, checkRegisterEmailSql+` UNION `+checkNewEmailSql+` LIMIT 1`, "email", email)
	row := r.db.QueryRowContext(ctx, checkRegisterEmailSql+` UNION `+checkNewEmailSql+` LIMIT 1`, email, email)
	err := row.Scan(&exists)
	if err != sql.ErrNoRows && err != nil {
		return false, err
	}
	return !exists, nil
}

func (r *ConfirmationTokenRepository) createToken(ctx context.Context, userId uint, tokenType string, payload []byte, days int) (string, error) {
	confirmation := uuidstr.GetUniqBase36(32)
	sql := fmt.Sprintf(createTokenSqlPattern, tokenType, days)
	_, err := r.db.ExecContext(ctx, sql, confirmation, userId, payload)
	if err != nil {
		return "", apperror.Wrap(err, "Failed to create reset token: "+sql)
	}
	return confirmation, nil
}
