package repository

import "backend/internal/database"

type UserRepository struct {
	db database.DBExecutor
}

func NewRepository(db database.DBExecutor) *UserRepository {
	return &UserRepository{db: db}
}
