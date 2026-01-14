package models

import "time"

type UserSessions struct {
	Id          uint      `db:"id" json:"id"`
	UserId      uint      `db:"user_id" json:"user_id"`
	TokenHash   string    `db:"token_hash" json:"token_hash"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	ExpiresAt   time.Time `db:"expires_at" json:"expires_at"`
	RefreshedAt time.Time `db:"refreshed_at" json:"refreshed_at"`
	RevokedAt   time.Time `db:"revoked_at" json:"revoked_at"`
	UserAgent   string    `db:"user_agent" json:"user_agent"`
	Ip          string    `db:"ip" json:"ip"`
}
