package service

import (
	"backend/internal/contexthelper"
	"backend/internal/session/repository"
	"context"
	"net/http"
)

type SessionService struct {
	sessionRepo *repository.UserSessionsRepository
	respWritter http.ResponseWriter
}

func NewService(ctx context.Context, w http.ResponseWriter) *SessionService {
	db := contexthelper.GetDb(ctx)
	sessionRepo := repository.NewRepository(db)
	return &SessionService{sessionRepo: sessionRepo, respWritter: w}
}
