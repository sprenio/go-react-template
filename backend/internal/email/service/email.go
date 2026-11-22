package service

import (
	"context"
	"encoding/json"
	"fmt"

	"backend/internal/apierrors"
	"backend/internal/confirmation"
	"backend/internal/contexthelper"
	"backend/internal/email"
	"backend/internal/payload"
	"backend/pkg/apperror"
	"backend/pkg/logger"
)

func (s *EmailService) sendWelcomeEmail(ctx context.Context, ct *confirmation.ConfirmationToken) error {
	var payloadData payload.RegisterQueuePayload
	err := json.Unmarshal([]byte(ct.Payload), &payloadData)
	if err != nil {
		return apperror.Wrap(email.ErrFailedToUnmarshalPayloadData, err.Error())
	}
	if payloadData.Email == "" {
		return apperror.Wrap(email.ErrInvalidEmailData, "email is empty")
	}

	lang, err := s.langRepo.GetLangByCode(ctx, payloadData.LanguageCode)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get language", "code", payloadData.LanguageCode, "error", err)
	}
	cfg := contexthelper.GetConfig(ctx)
	langCode := cfg.DefaultLanguage
	if lang.Code > "" {
		langCode = lang.Code
	}
	sender := email.GetEmailSender(ctx)
	if sender == nil {
		return email.ErrFailedToGetEmailSender
	}
	link := fmt.Sprintf("%s/confirm/%s", cfg.Frontend.BaseURL, ct.Token)
	err = sender.SendWelcomeEmail(ctx, payloadData.Email, payloadData.Name, langCode, link)
	if err != nil {
		return apperror.Wrap(email.ErrFailedToSendEmail, err.Error())
	}
	logger.InfoCtx(ctx, "Email sent", "to", payloadData.Email)
	return nil
}

func (s *EmailService) sendEmailChangeEmail(ctx context.Context, ct *confirmation.ConfirmationToken) error {
	if ct.UserId == 0 {
		return apperror.Wrap(email.ErrInvalidUserId, "empty user id in token")
	}

	var payloadData payload.EmailChangePayload
	err := json.Unmarshal([]byte(ct.Payload), &payloadData)
	if err != nil {
		return apperror.Wrap(email.ErrFailedToUnmarshalPayloadData, err.Error())
	}
	if payloadData.NewEmail == "" {
		return apperror.Wrap(email.ErrInvalidEmailData, "email is empty")
	}

	cfg := contexthelper.GetConfig(ctx)

	user, err := s.userRepo.GetById(ctx, ct.UserId)
	if err != nil {
		return apperror.Wrap(email.ErrInvalidUserId, err.Error())
	}
	if user.Id == 0 {
		return apperror.Wrap(email.ErrInvalidUserId, "user not found")
	}
	sender := email.GetEmailSender(ctx)
	if sender == nil {
		return email.ErrFailedToGetEmailSender
	}
	link := fmt.Sprintf("%s/confirm/%s", cfg.Frontend.BaseURL, ct.Token)
	err = sender.SendEmailChangeEmail(ctx, payloadData.NewEmail, user.Name, link)
	if err != nil {
		return apperror.Wrap(email.ErrFailedToSendEmail, err.Error())
	}
	logger.InfoCtx(ctx, "Email change email sent", "to", payloadData.NewEmail)
	return nil
}

func (s *EmailService) sendPasswordResetEmail(ctx context.Context, ct *confirmation.ConfirmationToken) error {
	if ct.UserId == 0 {
		return apperror.Wrap(email.ErrInvalidUserId, "empty user id in token")
	}

	user, err := s.userRepo.GetById(ctx, ct.UserId)
	if err != nil {
		return apperror.Wrap(email.ErrInvalidUserId, err.Error())
	}
	if user.Id == 0 {
		return apperror.Wrap(email.ErrInvalidUserId, "user not found")
	}
	sender := email.GetEmailSender(ctx)
	if sender == nil {
		return email.ErrFailedToGetEmailSender
	}

	cfg := contexthelper.GetConfig(ctx)

	link := fmt.Sprintf("%s/reset-password/%s", cfg.Frontend.BaseURL, ct.Token)
	err = sender.SendPasswordResetEmail(ctx, user.Email, user.Name, link)
	if err != nil {
		return apperror.Wrap(email.ErrFailedToSendEmail, err.Error())
	}
	logger.InfoCtx(ctx, "Password reset email sent", "to", user.Email)
	return nil
}

func (s *EmailService) SendConfirmationEmail(ctx context.Context, token string) error {
	confirmToken, err := s.confirmationRepo.GetActiveNewToken(ctx, token)
	if err != nil {
		return apperror.Wrap(email.ErrFailedToGetConfirmationToken, err.Error(), "token", token)
	}
	switch confirmToken.Type {
	case confirmation.ConfirmationTokenTypeRegister:
		return s.sendWelcomeEmail(ctx, &confirmToken)
	case confirmation.ConfirmationTokenTypeEmailChange:
		return s.sendEmailChangeEmail(ctx, &confirmToken)
	case confirmation.ConfirmationTokenTypeResetPassword:
		return s.sendPasswordResetEmail(ctx, &confirmToken)
	default:
		return apperror.Wrap(apierrors.ErrUnrecoverableError, "unknown confirmation token type")
	}
}
