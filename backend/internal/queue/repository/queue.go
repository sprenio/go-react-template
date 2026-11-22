package repository

import (
	"backend/internal/database"
	"backend/internal/queue"
	"backend/pkg/apperror"
	"backend/pkg/logger"
	"context"
	"time"
)

type QueueRepository struct {
	db database.DBExecutor
}

const queueTableName string = "queues"

func NewQueueRepository(db database.DBExecutor) *QueueRepository {
	return &QueueRepository{db: db}
}

func (r *QueueRepository) Create(ctx context.Context, task string, payload []byte, maxRetries uint) (uint, error) {
	sqlQuery := `INSERT INTO ` + queueTableName + ` (task, payload, max_retries) VALUES (?, ?, ?)`
	res, err := r.db.ExecContext(ctx, sqlQuery, task, string(payload), maxRetries)
	if err != nil {
		return 0, apperror.Wrap(err, "SQL query: "+sqlQuery)
	}
	tmp, err := res.LastInsertId()
	if err != nil {
		return 0, apperror.Wrap(err, "Error getting last insert id: "+sqlQuery)
	}
	return uint(tmp), nil
}

func (r *QueueRepository) GetById(ctx context.Context, id uint) (queue.Queue, error) {
	sqlQuery := `SELECT id, task, payload, status, retries_count, max_retries, generation, last_error, created_at, updated_at, next_retry_at FROM ` + queueTableName + ` WHERE id = ?`
	row := r.db.QueryRowContext(ctx, sqlQuery, id)
	var q queue.Queue
	if err := row.Scan(&q.Id, &q.Task, &q.Payload, &q.Status, &q.RetriesCount, &q.MaxRetries, &q.Generation, &q.LastError, &q.CreatedAt, &q.UpdatedAt, &q.NextRetryAt); err != nil {
		return queue.Queue{}, apperror.Wrap(err, "SQL query: "+sqlQuery)
	}
	return q, nil
}

func (r *QueueRepository) SetErrorStatus(ctx context.Context, id uint, errorMessage string, nextRetryAt *time.Time) error {
	logger.InfoCtx(ctx, "Setting error status for queue", "id", id, "errorMessage", errorMessage, "nextRetryAt", nextRetryAt)
	sqlQuery := `UPDATE ` + queueTableName + ` SET status = ?, last_error = ?, next_retry_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, sqlQuery, queue.QueueStatusFailed, errorMessage, nextRetryAt, id)
	if err != nil {
		logger.DebugCtx(ctx, "Error setting error status for queue", "id", id, "error", err)
		return apperror.Wrap(err, "SQL query: "+sqlQuery)
	}
	logger.DebugCtx(ctx, "Queue successfully set to error", "id", id)
	return nil
}

func (r *QueueRepository) SetQueuedStatus(ctx context.Context, id uint) (uint, error) {
	sqlQuery := `UPDATE ` + queueTableName + ` SET status = ?, retries_count = 0, last_error = NULL, next_retry_at = NULL, generation = generation + 1 WHERE id = ? AND status IN (?, ?) RETURNING generation`
	row := r.db.QueryRowContext(ctx, sqlQuery, queue.QueueStatusQueued, id, queue.QueueStatusRetryPending, queue.QueueStatusPending)
	var generation uint
	if err := row.Scan(&generation); err != nil {
		return 0, apperror.Wrap(err, "SQL query: "+sqlQuery)
	}
	return generation, nil
}

func (r *QueueRepository) GetRetryPendingQueues(ctx context.Context, limit int) ([]queue.Queue, error) {
	sqlQuery := `SELECT id, task, payload, status, retries_count, max_retries, generation, last_error, created_at, updated_at, next_retry_at FROM ` + queueTableName + ` WHERE status = ? ORDER BY id ASC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, sqlQuery, queue.QueueStatusRetryPending, limit)
	if err != nil {
		return nil, apperror.Wrap(err, "SQL query: "+sqlQuery)
	}
	defer rows.Close()

	var pendingQueues []queue.Queue
	for rows.Next() {
		var queueRecord queue.Queue
		if err := rows.Scan(
			&queueRecord.Id,
			&queueRecord.Task,
			&queueRecord.Payload,
			&queueRecord.Status,
			&queueRecord.RetriesCount,
			&queueRecord.MaxRetries,
			&queueRecord.Generation,
			&queueRecord.LastError,
			&queueRecord.CreatedAt,
			&queueRecord.UpdatedAt,
			&queueRecord.NextRetryAt,
		); err != nil {
			return nil, apperror.Wrap(err, "SQL query: "+sqlQuery)
		}
		pendingQueues = append(pendingQueues, queueRecord)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.Wrap(err, "SQL query: "+sqlQuery)
	}

	return pendingQueues, nil
}

func (r *QueueRepository) IncrementRetries(ctx context.Context, id uint) (uint, error) {
	sqlQuery := `UPDATE ` + queueTableName + ` SET retries_count = retries_count + 1 WHERE id = ? RETURNING retries_count`
	row := r.db.QueryRowContext(ctx, sqlQuery, id)
	var retriesCount uint
	if err := row.Scan(&retriesCount); err != nil {
		return 0, apperror.Wrap(err, "SQL query: "+sqlQuery)
	}
	return retriesCount, nil
}

func (r *QueueRepository) SetDoneStatus(ctx context.Context, id uint, retriesCount uint) error {
	sqlQuery := `UPDATE ` + queueTableName + ` SET status = ?, next_retry_at = NULL, retries_count = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, sqlQuery, queue.QueueStatusDone, retriesCount, id)
	if err != nil {
		return apperror.Wrap(err, "SQL query: "+sqlQuery)
	}
	return nil
}

func (r *QueueRepository) SetProcessingStatus(ctx context.Context, id uint) (bool, error) {
	logger.InfoCtx(ctx, "Setting processing status for queue", "id", id)
	sqlQuery := `UPDATE ` + queueTableName + ` SET status = ? WHERE id = ? AND status IN (?, ?)`
	res, err := r.db.ExecContext(ctx, sqlQuery, queue.QueueStatusProcessing, id, queue.QueueStatusQueued, queue.QueueStatusFailed)
	if err != nil {
		logger.DebugCtx(ctx, "Error setting processing status for queue", "id", id, "error", err)
		return false, apperror.Wrap(err, "SQL query: "+sqlQuery)
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		logger.DebugCtx(ctx, "Queue not found or status is not Queued or Failed", "id", id)
		return false, nil
	}
	logger.DebugCtx(ctx, "Queue successfully set to processing", "id", id)
	return true, nil
}

func (r *QueueRepository) UpdateNextRetryAt(ctx context.Context, id uint, nextRetryAt *time.Time) error {
	logger.InfoCtx(ctx, "Updating next retry at for queue", "id", id, "nextRetryAt", nextRetryAt)
	sqlQuery := `UPDATE ` + queueTableName + ` SET next_retry_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, sqlQuery, nextRetryAt, id)
	if err != nil {
		return apperror.Wrap(err, "SQL query: "+sqlQuery)
	}
	return nil
}

func (r *QueueRepository) ExistsRegisterQueueByUserName(ctx context.Context, userName string) (bool, error) {
	sqlQuery := `SELECT EXISTS(SELECT 1 FROM ` + queueTableName + ` WHERE task = ? AND LOWER(JSON_UNQUOTE(JSON_EXTRACT(payload, '$.name'))) = LOWER(?) AND status IN(?, ?, ?))`
	row := r.db.QueryRowContext(ctx, sqlQuery, queue.QueueTaskRegisterUser, userName, queue.QueueStatusPending, queue.QueueStatusQueued, queue.QueueStatusProcessing)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, apperror.Wrap(err, "SQL query: "+sqlQuery)
	}
	return exists, nil
}

func (r *QueueRepository) GetStatus(ctx context.Context, id uint) (queue.QueueStatus, error) {
	sqlQuery := `SELECT status FROM ` + queueTableName + ` WHERE id = ?`
	row := r.db.QueryRowContext(ctx, sqlQuery, id)
	var status queue.QueueStatus
	if err := row.Scan(&status); err != nil {
		return "", apperror.Wrap(err, "SQL query: "+sqlQuery)
	}
	return status, nil
}
