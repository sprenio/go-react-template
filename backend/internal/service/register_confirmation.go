package service

import (
	"backend/internal/models"
	"backend/internal/payload"
	"backend/internal/repository"
	"encoding/json"
	"time"
	"context"

	"github.com/pkg/errors"
)

type RegisterConfirmationService struct {
	uRepo  *repository.UserRepository
	usRepo *repository.UserSettingsRepository
}

func NewRegisterConfirmationService(uRepo *repository.UserRepository, usRepo *repository.UserSettingsRepository) *RegisterConfirmationService {
	return &RegisterConfirmationService{uRepo: uRepo, usRepo: usRepo}
}

func (s *RegisterConfirmationService) ConfirmRegisterToken(ctx context.Context, token models.ConfirmationToken) (uint, error) {
	var rp payload.RegisterPayload
	if err := json.Unmarshal([]byte(token.Payload), &rp); err != nil {
		return 0, errors.Wrap(err, "unmarshal user from confirmation token payload")
	}
	user := rp.ToUser()
	user.ConfirmedAt = time.Now()
	//user.RegisteredAt = token.CreatedAt
	id, err := s.uRepo.Create(ctx, user)
	if err != nil {
		return 0, errors.Wrap(err, "update user after confirmation")
	}
	if id == 0 {
		return 0, errors.New("user ID is zero after creation")
	}
	err = s.usRepo.Create(ctx, id, rp.LanguageId)
	if err != nil {
		return 0, errors.Wrap(err, "create user settings after confirmation")
	}
	
	return id, nil
}
