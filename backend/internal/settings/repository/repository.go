package repository

import "backend/internal/database"

type UserSettingsRepository struct {
	db database.DBExecutor
}

func NewRepository(db database.DBExecutor) *UserSettingsRepository {
	return &UserSettingsRepository{db: db}
}
