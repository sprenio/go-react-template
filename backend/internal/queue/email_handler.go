package queue

import (
	"backend/config"
	"backend/internal/email"
	"backend/internal/models"
	"backend/internal/payload"
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	emailMainQueue  = "email_tasks"
	emailRetryQueue = "email_tasks_retry"
	emailDLQQueue   = "email_tasks_dlq"

	emailMaxRetries = 3
	emailRetryDelay = 5000 // ms

	registerEmailTask      = "send_register_email"
	emailChangeEmailTask   = "send_email_change_email"
	passwordResetEmailTask = "send_password_reset_email"
)

type WelcomeEmailData struct {
	RegisterToken string `json:"register_token"`
}
type EmailChangeEmailData struct {
	EmailChangeToken string `json:"email_change_token"`
}
type PasswordResetEmailData struct {
	PasswordResetToken string `json:"password_reset_token"`
}

func (c *Consumer) HandleEmailTask(ctx context.Context, task string, rawMessage json.RawMessage) error {
	logger.Info("📧 starting handling email task: %s", task)

	switch task {
	case registerEmailTask:
		err := c.sendWelcomeEmail(ctx, rawMessage)
		if err != nil {
			return err
		}
	case emailChangeEmailTask:
		err := c.sendEmailChangeEmail(ctx, rawMessage)
		if err != nil {
			return err
		}
	case passwordResetEmailTask:
		err := c.sendPasswordResetEmail(ctx, rawMessage)
		if err != nil {
			return err
		}
	default:
		logger.Error("❌ Unknown email task: %s", task)
		return errors.New("unknown email task")
	}
	logger.Info("✅ Email task handled successfully")
	return nil
}

func setupEmailQueues(ch *amqp.Channel) error {
	return setupQueues(ch, emailMainQueue, emailRetryQueue, emailDLQQueue, int32(emailRetryDelay))
}

func (c *Consumer) sendWelcomeEmail(ctx context.Context, rawMessage json.RawMessage) error {
	var data WelcomeEmailData
	if err := json.Unmarshal(rawMessage, &data); err != nil {
		return err
	}

	if data.RegisterToken == "" {
		return errors.New("empty register token")
	}
	confirmationRepo := repository.NewConfirmationTokenRepository(c.db)
	ct, err := confirmationRepo.GetActiveNewTokenWithType(ctx, data.RegisterToken, models.ConfirmationTokenTypeRegister)
	if err != nil {
		return err
	}

	sender := email.GetEmailSender()
	if sender == nil {
		return errors.New("failed to get email sender")
	}
	var payloadData payload.RegisterPayload
	err = json.Unmarshal([]byte(ct.Payload), &payloadData)
	if err != nil {
		return err
	}
	if payloadData.Email == "" {
		logger.Error("Invalid email data")
		return errors.New("invalid email data")
	}
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}
	langRepo := repository.NewLanguageRepository(c.db)
	lang, err := langRepo.GetById(ctx, payloadData.LanguageId)
	if err != nil {
		logger.Error("Failed to get language by id: %d, error: %v", payloadData.LanguageId, err)
	}
	langCode := cfg.DefaultLanguage
	if lang.I18nCode > "" {
		langCode = lang.I18nCode
	}

	link := fmt.Sprintf("%s/confirm/%s", cfg.Frontend.BaseURL, data.RegisterToken)
	err = sender.SendWelcomeEmail(payloadData.Email, payloadData.Name, langCode, link)
	if err != nil {
		return err
	}
	logger.Info("Email sent to %s", payloadData.Email)
	return nil
}

func (c *Consumer) sendEmailChangeEmail(ctx context.Context, rawMessage json.RawMessage) error {
	var data EmailChangeEmailData
	if err := json.Unmarshal(rawMessage, &data); err != nil {
		return err
	}

	if data.EmailChangeToken == "" {
		return errors.New("empty email change token")
	}

	confirmationRepo := repository.NewConfirmationTokenRepository(c.db)
	ct, err := confirmationRepo.GetActiveNewTokenWithType(ctx, data.EmailChangeToken, models.ConfirmationTokenTypeEmailChange)
	if err != nil {
		return err
	}
	if ct.UserId == 0 {
		return errors.New("invalid user id in token")
	}

	sender := email.GetEmailSender()
	if sender == nil {
		return errors.New("failed to get email sender")
	}
	var payloadData payload.EmailChangePayload
	err = json.Unmarshal([]byte(ct.Payload), &payloadData)
	if err != nil {
		return err
	}
	if payloadData.NewEmail == "" {
		logger.Error("Invalid email data")
		return errors.New("invalid email data")
	}
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	userRepository := repository.NewUserRepository(c.db)
	user, err := userRepository.GetById(ctx, ct.UserId)
	if err != nil {
		return err
	}
	if user.Id == 0 {
		return errors.New("user not found")
	}

	link := fmt.Sprintf("%s/confirm/%s", cfg.Frontend.BaseURL, data.EmailChangeToken)
	err = sender.SendEmailChangeEmail(payloadData.NewEmail, user.Name, link)
	if err != nil {
		return err
	}
	logger.Info("Email change email sent to %s", payloadData.NewEmail)
	return nil
}

func (c *Consumer) sendPasswordResetEmail(ctx context.Context, rawMessage json.RawMessage) error {
	var data PasswordResetEmailData
	if err := json.Unmarshal(rawMessage, &data); err != nil {
		return err
	}

	if data.PasswordResetToken == "" {
		return errors.New("empty password reset token")
	}

	confirmationRepo := repository.NewConfirmationTokenRepository(c.db)
	ct, err := confirmationRepo.GetActiveNewTokenWithType(ctx, data.PasswordResetToken, models.ConfirmationTokenTypePasswordChange)
	if err != nil {
		return err
	}

	userRepository := repository.NewUserRepository(c.db)
	user, err := userRepository.GetById(ctx, ct.UserId)
	if err != nil {
		return err
	}
	sender := email.GetEmailSender()
	if sender == nil {
		return errors.New("failed to get email sender")
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	link := fmt.Sprintf("%s/reset-password/%s", cfg.Frontend.BaseURL, data.PasswordResetToken)
	err = sender.SendPasswordResetEmail(user.Email, user.Name, link)
	if err != nil {
		return err
	}
	logger.Info("Password reset email sent to %s", user.Email)
	return nil
}
