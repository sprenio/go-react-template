package models

import "time"

type User struct {
    Id           uint        `db:"id" json:"id"`
    Name         string     `db:"name" json:"name"`
    Email        string     `db:"email" json:"email"`
    Password     string     `db:"password" json:"-"` // nie zwracamy hasła w JSON
    RegisteredAt time.Time  `db:"registered_at" json:"registered_at"`
    ConfirmedAt  time.Time `db:"confirmed_at" json:"confirmed_at,omitempty"`
}

