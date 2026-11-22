package repository

import "backend/internal/database"

type UserSessionsRepository struct {
	db database.DBExecutor
}

func NewRepository(db database.DBExecutor) *UserSessionsRepository {
	return &UserSessionsRepository{db: db}
}
