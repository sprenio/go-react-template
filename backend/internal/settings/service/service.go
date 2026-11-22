package service

import (
	"backend/internal/contexthelper"
	languageRepository "backend/internal/language/repository"
	settingsRepository "backend/internal/settings/repository"
	"context"
)

type UserSettingsService struct {
	usRepo   *settingsRepository.UserSettingsRepository
	langRepo *languageRepository.LanguageRepository
}

func NewService(ctx context.Context) *UserSettingsService {
	db := contexthelper.GetDb(ctx)
	usRepo := settingsRepository.NewRepository(db)
	langRepo := languageRepository.NewRepository(db)
	return &UserSettingsService{usRepo: usRepo, langRepo: langRepo}
}
