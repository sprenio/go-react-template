package service

import (
	confirmationRepository "backend/internal/confirmation/repository"
	"backend/internal/contexthelper"
	languageRepository "backend/internal/language/repository"
	userRepository "backend/internal/user/repository"
	"context"
)

type EmailService struct {
	confirmationRepo *confirmationRepository.ConfirmationTokenRepository
	langRepo         *languageRepository.LanguageRepository
	userRepo         *userRepository.UserRepository
}

func NewService(ctx context.Context) *EmailService {
	db := contexthelper.GetDb(ctx)
	confirmationRepo := confirmationRepository.NewRepository(db)
	langRepo := languageRepository.NewRepository(db)
	userRepo := userRepository.NewRepository(db)
	return &EmailService{confirmationRepo: confirmationRepo, langRepo: langRepo, userRepo: userRepo}
}
