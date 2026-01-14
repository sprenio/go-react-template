package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(uRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: uRepo}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (models.UserResponseData, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user.Id == 0 {
		return models.UserResponseData{}, errors.Wrap(err, "user not found")
	}

	// Sprawdzenie hasła bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.UserResponseData{}, errors.Wrap(err, "invalid password")
	}
	userData, settings, err := s.userRepo.GetDataById(ctx, user.Id)
	if err != nil {
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
