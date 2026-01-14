package service

import (
	"backend/internal/apperrors"
	"backend/internal/contexthelper"
	"backend/internal/queue"
	"backend/internal/repository"
	"backend/pkg/logger"
	"backend/pkg/validation"
	"context"
	"database/sql"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
	confirmationTokenRepo *repository.ConfirmationTokenRepository
	userRepo              *repository.UserRepository
}

func NewPasswordService(confirmationTokenRepo *repository.ConfirmationTokenRepository, userRepo *repository.UserRepository) *PasswordService {
	return &PasswordService{
		confirmationTokenRepo: confirmationTokenRepo,
		userRepo:              userRepo,
	}
}
func (s *PasswordService) ResetPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get user by email: %v", err)
		if err == sql.ErrNoRows {
			r := rand.Intn(100) + 50
			// Sleep for a random time between 50 and 150 microseconds to mitigate timing attacks
			time.Sleep(time.Duration(r) * time.Microsecond)
			return nil // Do not reveal whether the email exists
		}
		return err
	}

	cfg := contexthelper.GetConfig(ctx)

	confirmationToken, err := s.confirmationTokenRepo.CreatePasswordChangeToken(ctx, user.Id, cfg.ResetPassword.ExpirationDays)
	if err != nil {
		return err
	}
	rabbitConn := contexthelper.GetRabbitConn(ctx)
	err = queue.PublishPasswordResetTask(ctx, rabbitConn, confirmationToken)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to enqueue password reset task: %v", err)
	}

	return nil
}
func (s *PasswordService) PasswordChange(ctx context.Context, UserId uint, newPassword string) error {
	if !validation.IsPasswordValid(newPassword) {
		return apperrors.NewInvalidInputError("Password", "Password must be at least 6 characters and contain a letter, number, and special character")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	strPassword := string(hashedPassword)
	err = s.userRepo.UpdatePasswordById(ctx, UserId, strPassword)
	if err != nil {
		return err
	}
	return nil
}
