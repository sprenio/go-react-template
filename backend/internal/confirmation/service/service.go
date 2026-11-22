package service

import (
	confirmationRepository "backend/internal/confirmation/repository"
	"backend/internal/contexthelper"
	settingsRepository "backend/internal/settings/repository"
	userRepository "backend/internal/user/repository"
	"backend/pkg/apperror"
	"context"
)

type ConfirmationService struct {
	confirmationRepo *confirmationRepository.ConfirmationTokenRepository
	userRepo         *userRepository.UserRepository
	userSettingsRepo *settingsRepository.UserSettingsRepository
}

type txContext struct {
	ConfirmationRepo *confirmationRepository.ConfirmationTokenRepository
	UserRepo         *userRepository.UserRepository
	UserSettingsRepo *settingsRepository.UserSettingsRepository
}

func NewService(ctx context.Context) *ConfirmationService {
	db := contexthelper.GetDb(ctx)

	return &ConfirmationService{
		confirmationRepo: confirmationRepository.NewRepository(db),
		userRepo:         userRepository.NewRepository(db),
		userSettingsRepo: settingsRepository.NewRepository(db),
	}
}

func (s *ConfirmationService) execInTx(ctx context.Context, fn func(txCtx *txContext) error) error {
	db := contexthelper.GetDb(ctx)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Wrap(err, "failed to begin transaction")
	}

	txCtx := &txContext{
		ConfirmationRepo: confirmationRepository.NewRepository(tx),
		UserRepo:         userRepository.NewRepository(tx),
		UserSettingsRepo: settingsRepository.NewRepository(tx),
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return apperror.Wrap(err, "failed to commit transaction")
	}

	return nil
}
