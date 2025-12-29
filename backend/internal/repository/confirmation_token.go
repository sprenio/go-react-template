package repository

import (
	"backend/internal/models"
	"backend/pkg/uuidstr"
	"database/sql"

	"context"
	"fmt"

	"github.com/pkg/errors"
)

const (
	ConfirmationTokenTable = "confirmation_tokens"

	confirmationTokenColumns = `id, token, user_id, type, payload, status, expires_at, status_changed_at, created_at`

	getActiveNewTokenSql = `SELECT ` + confirmationTokenColumns + ` FROM ` + ConfirmationTokenTable + ` WHERE token = ? AND status = "` + models.ConfirmationTokenStatusNew + `" AND expires_at > NOW()`

	// Check for active tokens with status 'new' or 'processing'
	activeStatusWhereCondition = `status IN ("` + models.ConfirmationTokenStatusNew + `","` + models.ConfirmationTokenStatusProcessing + `")`
	checkRegisterEmailSql      = `SELECT 1 FROM ` + ConfirmationTokenTable + ` WHERE type ='` + models.ConfirmationTokenTypeRegister + `' AND JSON_UNQUOTE(JSON_EXTRACT(payload, '$.email')) = ? AND ` + activeStatusWhereCondition
	checkNewEmailSql           = `SELECT 1 FROM ` + ConfirmationTokenTable + ` WHERE type ='` + models.ConfirmationTokenTypeEmailChange + `' AND JSON_UNQUOTE(JSON_EXTRACT(payload, '$.new_email')) = ? AND ` + activeStatusWhereCondition
	createTokenSqlPattern      = `INSERT INTO ` + ConfirmationTokenTable + ` (token, user_id, type, payload, status, expires_at, status_changed_at) VALUES (?, ?, '%s', ?, "` + models.ConfirmationTokenStatusNew + `", DATE_ADD(NOW(), INTERVAL %d DAY), NOW())`
)

type ConfirmationTokenRepository struct {
	db DBExecutor
}

func NewConfirmationTokenRepository(db DBExecutor) *ConfirmationTokenRepository {
	return &ConfirmationTokenRepository{db: db}
}

func (r *ConfirmationTokenRepository) GetActiveNewToken(ctx context.Context, token string) (models.ConfirmationToken, error) {
	var ct models.ConfirmationToken

	err := r.db.QueryRowContext(ctx, getActiveNewTokenSql, token).Scan(&ct.Id, &ct.Token, &ct.UserId, &ct.Type, &ct.Payload, &ct.Status, &ct.ExpiresAt, &ct.StatusChangedAt, &ct.CreatedAt)
	if err != nil || ct.Id == 0 {
		return models.ConfirmationToken{}, errors.Wrap(err, "Failed to retrieve active confirmation token: "+token)
	}
	return ct, nil
}

func (r *ConfirmationTokenRepository) ExistsRegisterTokenByEmailOrName(ctx context.Context, username, email string) (error, bool) {
	var exists bool
	row := r.db.QueryRowContext(ctx, `SELECT 1 FROM `+ConfirmationTokenTable+` WHERE type = '`+models.ConfirmationTokenTypeRegister+`' AND JSON_UNQUOTE(JSON_EXTRACT(payload, '$.username')) = ? AND `+activeStatusWhereCondition+`
		UNION `+checkRegisterEmailSql+` UNION `+checkNewEmailSql+`
		LIMIT 1`,
		username,
		email,
		email,
	)
	err := row.Scan(&exists)
	if err != sql.ErrNoRows && err != nil {
		return err, false
	}
	return nil, exists
}

func (r *ConfirmationTokenRepository) ExistsRegisterTokenByEmail(ctx context.Context, email string) (error, bool) {
	return r.existsTokenByEmail(ctx, email)
}
func (r *ConfirmationTokenRepository) ExistsEmailChangeTokenByEmail(ctx context.Context, email string) (error, bool) {
	return r.existsTokenByEmail(ctx, email)
}

func (r *ConfirmationTokenRepository) CreateRegisterToken(ctx context.Context, payload []byte, days int) (string, error) {
	return r.createToken(ctx, 0, models.ConfirmationTokenTypeRegister, payload, days)
}

func (r *ConfirmationTokenRepository) CreateEmailChangeToken(ctx context.Context, userId uint, payload []byte, days int) (string, error) {
	return r.createToken(ctx, userId, models.ConfirmationTokenTypeEmailChange, payload, days)
}
func (r *ConfirmationTokenRepository) CreatePasswordChangeToken(ctx context.Context, userId uint, days int) (string, error) {
	return r.createToken(ctx, userId, models.ConfirmationTokenTypePasswordChange, []byte("{}"), days)
}

func (r *ConfirmationTokenRepository) ConsumeToken(ctx context.Context, token string) error {
	return r.updateStatus(ctx, token, models.ConfirmationTokenStatusConsumed)
}

func (r *ConfirmationTokenRepository) GetActiveNewTokenWithType(ctx context.Context, token string, tokenType string) (models.ConfirmationToken, error) {
	var ct models.ConfirmationToken
	err := r.db.QueryRowContext(ctx, getActiveNewTokenSql+` AND type = ?`, token, tokenType).Scan(&ct.Id, &ct.Token, &ct.UserId, &ct.Type, &ct.Payload, &ct.Status, &ct.ExpiresAt, &ct.StatusChangedAt, &ct.CreatedAt)
	if err != nil {
		return models.ConfirmationToken{}, err
	}
	return ct, nil
}

func (r *ConfirmationTokenRepository) updateStatus(ctx context.Context, token string, status string) error {
	sql := `UPDATE ` + ConfirmationTokenTable + ` SET status = ?, status_changed_at = NOW() WHERE token = ? AND status != ?`
	_, err := r.db.ExecContext(ctx, sql, status, token, status)
	if err != nil {
		return errors.Wrap(err, "Failed to update confirmation token status")
	}
	return nil
}

func (r *ConfirmationTokenRepository) existsTokenByEmail(ctx context.Context, email string) (error, bool) {
	var exists bool
	row := r.db.QueryRowContext(ctx, checkRegisterEmailSql+` UNION `+checkNewEmailSql+` LIMIT 1`, email, email)
	err := row.Scan(&exists)
	if err != sql.ErrNoRows && err != nil {
		return err, false
	}
	return nil, exists
}
func (r *ConfirmationTokenRepository) createToken(ctx context.Context, userId uint, tokenType string, payload []byte, days int) (string, error) {
	confirmationToken := uuidstr.GetUniqBase36(32)
	sql := fmt.Sprintf(createTokenSqlPattern, tokenType, days)
	_, err := r.db.ExecContext(ctx, sql, confirmationToken, userId, payload)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create reset token: "+sql)
	}
	return confirmationToken, nil
}
