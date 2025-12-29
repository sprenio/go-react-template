package service

import (
	"backend/internal/repository"
	"backend/internal/response"

	"context"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (response.UserResponseData, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil || user.Id == 0 {
		return response.UserResponseData{}, errors.Wrap(err, "user not found")
	}

	// Sprawdzenie hasła bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return response.UserResponseData{}, errors.Wrap(err, "invalid password")
	}
	userData, err := s.repo.GetDataById(ctx, user.Id)
	if err != nil {
		return response.UserResponseData{}, errors.Wrap(err, "user not found")
	}
	return userData, nil
}
