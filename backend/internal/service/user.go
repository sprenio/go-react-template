package service

import (
	"backend/internal/models"
	"backend/internal/payload"
	"backend/internal/repository"
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserResponseData(ctx context.Context, userId uint) (models.UserResponseData, error) {
	userData, settings, err := s.repo.GetDataById(ctx, userId)
	if err != nil || userData.Id == 0 {
		return models.UserResponseData{}, errors.Wrap(err, "user not found")
	}
	userResponseData := models.UserResponseData{
		Id:           userData.Id,
		Name:         userData.Name,
		Email:        userData.Email,
		RegisteredAt: userData.RegisteredAt.Format("2006-01-02 15:04:05"),
		ConfirmedAt:  userData.ConfirmedAt.Format("2006-01-02 15:04:05"),
		Settings:     settings,
	}

	return userResponseData, nil
}

func (s *UserService) ConfirmEmailChangeToken(ctx context.Context, token models.ConfirmationToken) error {
	var rp payload.EmailChangePayload
	if err := json.Unmarshal([]byte(token.Payload), &rp); err != nil {
		return errors.Wrap(err, "unmarshal user from confirmation token payload")
	}
	err := s.repo.ChangeEmail(ctx, token.UserId, rp.NewEmail)
	if err != nil {
		return errors.Wrap(err, "update user after confirmation")
	}
	return nil
}
