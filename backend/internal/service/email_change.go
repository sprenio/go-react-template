package service

import (
	"backend/internal/apperrors"
	"backend/config"
	"backend/internal/contexthelper"
	"backend/internal/queue"
	"backend/internal/payload"
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"encoding/json"
	"strings"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Email struct {
	confirmationTokenRepo *repository.ConfirmationTokenRepository
	userRepo              *repository.UserRepository
	languageRepo          *repository.LanguageRepository
}

func NewEmailService(ctRepo *repository.ConfirmationTokenRepository, uRepo *repository.UserRepository, langRepo *repository.LanguageRepository) *Email {
	return &Email{
		confirmationTokenRepo: ctRepo,
		userRepo:              uRepo,
		languageRepo:          langRepo,
	}
}

func (s *Email) ChangeEmail(ctx context.Context, rabbitConn *amqp.Connection, newEmail string) error {
	userId, ok := contexthelper.GetUserId(ctx)
	if !ok {
		return apperrors.NewGeneralCustomError("User not authenticated")
	}
	lowercaseEmail := strings.ToLower(newEmail)
	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		return err
	}
	if user.Email == lowercaseEmail {
		return apperrors.NewEmailChangeSameEmailError("New email is the same as the current one")
	}
	// Check if email or username exists
	err, exists := s.userRepo.ExistsByEmail(ctx, lowercaseEmail)
	if err != nil {
		return err
	}
	if exists {
		return apperrors.NewEmailChangeEmailAlreadyUsedError("Email already in use")
	}
	err, exists = s.confirmationTokenRepo.ExistsEmailChangeTokenByEmail(ctx, lowercaseEmail)
	if err != nil {
		return err
	}
	if exists {
		logger.Info("Email change token already exists for email: %s", lowercaseEmail)
		return apperrors.NewEmailChangeEmailAlreadyUsedError("Email already exists")
	}

	payload := payload.EmailChangePayload{
		NewEmail: lowercaseEmail,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	confirmationToken, err := s.confirmationTokenRepo.CreateEmailChangeToken(ctx, user.Id, jsonPayload, cfg.EmailChange.ExpirationDays)
	if err != nil {
		return err
	}
	err = queue.PublishEmailChangeTask(ctx, rabbitConn, confirmationToken)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to enqueue email task: %v", err)
		// Nie przerywamy rejestracji – można uznać, że task się nie udał, ale user się zarejestrował
	}

	return nil
}
