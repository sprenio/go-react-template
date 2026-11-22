package queue

import "time"

type QueueTask string
type QueueStatus string

const (
	QueueTaskRegisterUser  QueueTask = "register_user"
	QueueTaskResetPassword QueueTask = "reset_password"
	QueueTaskEmail         QueueTask = "email"
)
const (
	EmailLogStatusStarted QueueStatus = "STARTED"
	EmailLogStatusSent    QueueStatus = "SENT"
	EmailLogStatusFailed  QueueStatus = "FAILED"
)

func (t QueueTask) String() string {
	return string(t)
}

const (
	QueueStatusPending      QueueStatus = "PENDING"
	QueueStatusQueued       QueueStatus = "QUEUED"
	QueueStatusProcessing   QueueStatus = "PROCESSING"
	QueueStatusDone         QueueStatus = "DONE"
	QueueStatusFailed       QueueStatus = "FAILED"
	QueueStatusRetryPending QueueStatus = "RETRY_PENDING"
)

func (s QueueStatus) String() string {
	return string(s)
}

const InitialGeneration uint = 1

type Queue struct {
	Id           uint        `db:"id" json:"id"`
	Task         QueueTask   `db:"task" json:"task"`
	Payload      []byte      `db:"payload" json:"-"`
	Status       QueueStatus `db:"status" json:"status"`
	MaxRetries   uint        `db:"max_retries" json:"max_retries"`
	Generation   uint        `db:"generation" json:"generation"`
	RetriesCount uint        `db:"retries_count" json:"retries_count"`
	LastError    *string     `db:"last_error" json:"last_error"`
	CreatedAt    time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time   `db:"updated_at" json:"updated_at"`
	NextRetryAt  *time.Time  `db:"next_retry_at" json:"next_retry_at"`
}
