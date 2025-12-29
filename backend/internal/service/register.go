package service

import (
	"backend/config"
	"backend/internal/apperrors"
	"backend/internal/models"
	"backend/internal/payload"
	"backend/internal/queue"
	"backend/internal/repository"
	"backend/pkg/logger"
	"backend/pkg/validation"
	"time"

	"context"
	"encoding/json"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/crypto/bcrypt"
)

type RegisterService struct {
	confirmationTokenRepo *repository.ConfirmationTokenRepository
	userRepo              *repository.UserRepository
	languageRepo          *repository.LanguageRepository
}

func NewRegisterService(ctRepo *repository.ConfirmationTokenRepository, uRepo *repository.UserRepository, langRepo *repository.LanguageRepository) *RegisterService {
	return &RegisterService{
		confirmationTokenRepo: ctRepo,
		userRepo:              uRepo,
		languageRepo:          langRepo,
	}
}

func (s *RegisterService) RegisterUser(ctx context.Context, rabbitConn *amqp.Connection, userName, email, password, langCode string) error {
	// Implement the user registration logic here
	if !validation.IsEmailValid(email) {
		return apperrors.NewInvalidInputError("Email", "invalid email format")
	}

	if !validation.IsPasswordValid(password) {
		return apperrors.NewInvalidInputError("Password", "Password must be at least 6 characters and contain a letter, number, and special character")
	}

	lowercaseEmail := strings.ToLower(email)

	// Check if email or username exists
	err, exists := s.userRepo.ExistsByEmailOrName(ctx, lowercaseEmail, userName)
	if err != nil {
		return err
	}
	if exists {
		return apperrors.NewRegisterUserNameOrEmailTakenError("Username or email already taken")
	}
	err, exists = s.confirmationTokenRepo.ExistsRegisterTokenByEmailOrName(ctx, userName, lowercaseEmail)
	if err != nil {
		return err
	}
	if exists {
		logger.Info("Register token already exists for user: %s or email: %s", userName, lowercaseEmail)
		return apperrors.NewRegisterUserNameOrEmailTakenError("Register token already exists")
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	lang, err := s.languageRepo.GetLangByCode(ctx, langCode)
	if err != nil {
		return err
	}
	if lang.Id == 0 {
		lang, err = s.languageRepo.GetLangByCode(ctx, cfg.DefaultLanguage)
		if err != nil {
			return err
		}
	}

	strPassword := string(hashedPassword)
	payload := payload.RegisterPayload{
		User: models.User{
			Name:         userName,
			Email:        lowercaseEmail,
			Password:     strPassword,
			RegisteredAt: time.Now(),
		},
		NewPassword: strPassword,
		LanguageId:  lang.Id,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	var registerToken string
	registerToken, err = s.confirmationTokenRepo.CreateRegisterToken(ctx, jsonPayload, cfg.Register.ExpirationDays)
	if err != nil {
		return err
	}
	err = queue.PublishRegisterEmailTask(ctx, rabbitConn, registerToken)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to enqueue email task: %v", err)
		// Nie przerywamy rejestracji – można uznać, że task się nie udał, ale user się zarejestrował
	}

	return nil
}
