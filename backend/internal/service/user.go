package service

import (
	"backend/internal/repository"
	"backend/internal/response"
	"backend/internal/models"
	"backend/internal/payload"

	"context"
	"github.com/pkg/errors"
	"encoding/json"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserResponseData(ctx context.Context, userId uint) (response.UserResponseData, error) {
	user, err := s.repo.GetDataById(ctx, userId)
	if err != nil || user.Id == 0 {
		return response.UserResponseData{}, errors.Wrap(err, "user not found")
	}
	return user, nil
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