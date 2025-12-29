package service

import (
	"backend/internal/apperrors"
	"backend/internal/contexthelper"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/pkg/logger"
	"context"
	"time"
)

type UserSettingsService struct {
	usRepo   *repository.UserSettingsRepository
	langRepo *repository.LanguageRepository
}

func NewUserSettingsService(usRepo *repository.UserSettingsRepository, langRepo *repository.LanguageRepository) *UserSettingsService {
	return &UserSettingsService{usRepo: usRepo, langRepo: langRepo}
}

func (s *UserSettingsService) Update(ctx context.Context, data models.UserSettingsData) error {

	userID, ok := contexthelper.GetUserId(ctx)
	if !ok {
		return apperrors.NewGeneralCustomError("Get user error")
	}
	us, err := s.usRepo.GetByUserId(ctx, userID)
	if err != nil {
		return err
	}
	doUpdate := false
	if data.Language != nil {
		lang, err := s.langRepo.GetLangByCode(ctx, *data.Language)
		if err != nil {
			return err
		}
		if lang.Id != us.LangId {
			doUpdate = true
			us.LangId = lang.Id
		}
	}
	if data.AppFlags != nil {
		flags := us.AppFlags
		us.SetAppFlags(*data.AppFlags)
		if flags != us.AppFlags {
			doUpdate = true
		}
	}
	if data.UserFlags != nil {
		flags := us.UserFlags
		us.SetUserFlags(*data.UserFlags)
		if flags != us.UserFlags {
			doUpdate = true
		}
	}
	if data.AppOpt1 != nil && data.AppOpt1 != &us.AppOpt1 {
		us.AppOpt1 = *data.AppOpt1
		doUpdate = true
	}
	if data.AppOpt2 != nil && data.AppOpt2 != &us.AppOpt2 {
		us.AppOpt2 = *data.AppOpt2
		doUpdate = true
	}
	if data.AppOpt3 != nil && data.AppOpt3 != &us.AppOpt3 {
		us.AppOpt3 = *data.AppOpt3
		doUpdate = true
	}
	if doUpdate {
		us.UpdatedAt = time.Now()
		logger.DebugCtx(ctx, "Update user settings: %v", us)
		err = s.usRepo.Update(ctx, us)
	} else {
		logger.DebugCtx(ctx, "No user settings update needed: %v", us)
	}

	return err
}
