package service

import (
	"backend/internal/apierrors"
	"backend/internal/confirmation"
	"backend/internal/contexthelper"
	"backend/internal/helper"
	languageService "backend/internal/language/service"
	"backend/internal/payload"
	"backend/pkg/apperror"
	"backend/pkg/validation"
	"context"
	"encoding/json"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *ConfirmationService) CreateRegisterToken(ctx context.Context, payload payload.RegisterQueuePayload) (string, error) {
	cfg := contexthelper.GetConfig(ctx)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return s.confirmationRepo.CreateRegisterToken(ctx, jsonPayload, cfg.Register.ExpirationDays)
}
func (s *ConfirmationService) ConfirmRegisterToken(ctx context.Context, token confirmation.ConfirmationToken) (uint, error) {
	var rp payload.RegisterQueuePayload
	if err := json.Unmarshal([]byte(token.Payload), &rp); err != nil {
		return 0, apperror.Wrap(apierrors.ErrUnrecoverableError, "unmarshal user from confirmation token payload", "error", err)
	}
	lngService := languageService.NewService(ctx)
	lang, err := lngService.GetLangByCode(ctx, rp.LanguageCode)
	if err != nil {
		return 0, err
	}
	if lang.Id == 0 {
		config := contexthelper.GetConfig(ctx)
		lang, err = lngService.GetLangByCode(ctx, config.DefaultLanguage)
		if err != nil {
			return 0, err
		}
	}
	if lang.Id == 0 {
		return 0, apperror.New("language not found")
	}
	var createdID uint
	err = s.execInTx(ctx, func(txCtx *txContext) error {
		user := rp.ToUser()
		user.ConfirmedAt = time.Now()
		//user.RegisteredAt = token.CreatedAt
		id, err := txCtx.UserRepo.Create(ctx, user)
		if err != nil {
			return apperror.Wrap(err, "update user after confirmation")
		}
		if id == 0 {
			return apperror.New("user ID is zero after creation")
		}
		err = txCtx.UserSettingsRepo.Create(ctx, id, lang.Id)
		if err != nil {
			return apperror.Wrap(err, "create user settings after confirmation")
		}
		createdID = id
		return nil
	})
	if err != nil {
		return 0, err
	}

	return createdID, nil
}

func (s *ConfirmationService) CreateResetPasswordToken(ctx context.Context, userId uint) (string, error) {
	cfg := contexthelper.GetConfig(ctx)
	return s.confirmationRepo.CreateResetPasswordToken(ctx, userId, cfg.ResetPassword.ExpirationDays)

}

func (s *ConfirmationService) CreateEmailChangeToken(ctx context.Context, payload payload.ChangeEmailQueuePayload) (string, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	cfg := contexthelper.GetConfig(ctx)
	return s.confirmationRepo.CreateEmailChangeToken(ctx, payload.UserId, jsonPayload, cfg.EmailChange.ExpirationDays)
}

func (s *ConfirmationService) ChangePassword(ctx context.Context, userId uint, newPassword string) error {
	err := validation.IsPasswordValid(newPassword)
	if err != nil {
		return apierrors.NewInvalidInputError("Password", err.Error(), &newPassword)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	strPassword := string(hashedPassword)
	err = s.userRepo.UpdatePasswordById(ctx, userId, strPassword)
	if err != nil {
		return err
	}
	return nil
}
func (s *ConfirmationService) ConfirmEmailChangeToken(ctx context.Context, token confirmation.ConfirmationToken) error {
	var rp payload.EmailChangePayload
	if err := json.Unmarshal([]byte(token.Payload), &rp); err != nil {
		return apperror.Wrap(err, "unmarshal user from confirmation token payload")
	}
	cleanEmail := helper.CleanEmail(rp.NewEmail)
	normalizedEmail := helper.NormalizeEmail(cleanEmail)
	err := s.userRepo.ChangeEmail(ctx, token.UserId, cleanEmail, normalizedEmail)
	if err != nil {
		return apperror.Wrap(err, "update user after confirmation")
	}
	return nil
}

func (s *ConfirmationService) GetActiveNewToken(ctx context.Context, token string) (confirmation.ConfirmationToken, error) {
	return s.confirmationRepo.GetActiveNewToken(ctx, token)
}

func (s *ConfirmationService) ConsumeToken(ctx context.Context, token string) error {
	return s.confirmationRepo.ConsumeToken(ctx, token)
}
