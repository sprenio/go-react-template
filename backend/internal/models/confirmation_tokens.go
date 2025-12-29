package models

import (
	"time"
)

type ConfirmationToken struct {
	Id              uint      `db:"id" json:"id"`
	Token           string    `db:"token" json:"token"`
	UserId          uint      `db:"user_id" json:"user_id"`
	Type            string    `db:"type" json:"type"`
	Payload         string    `db:"payload" json:"payload"`
	Status          string    `db:"status" json:"status"`
	ExpiresAt       time.Time `db:"expires_at" json:"expires_at"`
	StatusChangedAt time.Time `db:"status_changed_at" json:"status_changed_at"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

const (
	ConfirmationTokenTypeRegister      = "register"
	ConfirmationTokenTypeEmailChange   = "email_change"
	ConfirmationTokenTypePasswordChange = "password_change"

	ConfirmationTokenStatusNew        = "NEW"
	ConfirmationTokenStatusExpired    = "EXPIRED"
	ConfirmationTokenStatusCanceled   = "CANCELED"
	ConfirmationTokenStatusProcessing = "PROCESSING"
	ConfirmationTokenStatusFailed     = "FAILED"
	ConfirmationTokenStatusConsumed   = "CONSUMED"
)
