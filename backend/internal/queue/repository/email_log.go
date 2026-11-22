package repository

import (
	"backend/internal/database"
	"backend/internal/queue"
	"backend/pkg/apperror"
	"context"
)

type EmailLogRepository struct {
	db database.DBExecutor
}

const emailLogTableName string = "email_log"

func NewEmailLogRepository(db database.DBExecutor) *EmailLogRepository {
	return &EmailLogRepository{db: db}
}

func (r *EmailLogRepository) Create(ctx context.Context, queueId uint, generation uint, retryNumber uint) (uint, error) {
	sqlQuery := `INSERT INTO ` + emailLogTableName + ` (queue_id, generation, retry_number) VALUES (?, ?, ?)`
	res, err := r.db.ExecContext(ctx, sqlQuery, queueId, generation, retryNumber)
	if err != nil {
		return 0, apperror.Wrap(err, "SQL query: "+sqlQuery)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, apperror.Wrap(err, "Error getting last insert id: "+sqlQuery)
	}
	return uint(id), nil
}

func (r *EmailLogRepository) SetErrorStatus(ctx context.Context, id uint, errorString string) (int64, error) {
	sqlQuery := `UPDATE ` + emailLogTableName + ` SET status = ?, log = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, sqlQuery, queue.EmailLogStatusFailed, errorString, id)
	if err != nil {
		return 0, apperror.Wrap(err, "SQL query: "+sqlQuery)
	}
	return res.RowsAffected()
}

func (r *EmailLogRepository) SetSentStatus(ctx context.Context, id uint) (int64, error) {
	sqlQuery := `UPDATE ` + emailLogTableName + ` SET status = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, sqlQuery, queue.EmailLogStatusSent, id)
	if err != nil {
		return 0, apperror.Wrap(err, "SQL query: "+sqlQuery)
	}
	return res.RowsAffected()
}
