package service_test

import (
	"backend/internal/contexthelper"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/service"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserSettingsService_Update_Success_ChangeLanguage(t *testing.T) {
	ctx := context.Background()
	ctx = contexthelper.SetUserId(ctx, 1)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mock user settings lookup
	mock.ExpectQuery("SELECT.*FROM user_settings").WillReturnRows(
		sqlmock.NewRows([]string{"id", "user_id", "lang_id", "user_flags", "app_flags", "app_opt_1", "app_opt_2", "app_opt_3", "updated_at"}).
			AddRow(1, 1, 1, 0, 0, "", "", "", time.Now()),
	)

	// Mock language lookup
	mock.ExpectQuery("SELECT.*FROM languages WHERE code").WillReturnRows(
		sqlmock.NewRows([]string{"id", "code", "i18n_code", "name"}).
			AddRow(2, "pl", "pl", "Polish"),
	)

	// Mock user settings update
	mock.ExpectExec("UPDATE user_settings").WillReturnResult(sqlmock.NewResult(0, 1))

	usRepo := repository.NewUserSettingsRepository(db)
	langRepo := repository.NewLanguageRepository(db)
	settingsService := service.NewUserSettingsService(usRepo, langRepo)

	lang := "pl"
	data := models.UserSettingsData{
		Language: &lang,
	}
	err = settingsService.Update(ctx, data)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserSettingsService_Update_NoUpdateNeeded(t *testing.T) {
	ctx := context.Background()
	ctx = contexthelper.SetUserId(ctx, 1)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mock user settings lookup
	mock.ExpectQuery("SELECT.*FROM user_settings").WillReturnRows(
		sqlmock.NewRows([]string{"id", "user_id", "lang_id", "user_flags", "app_flags", "app_opt_1", "app_opt_2", "app_opt_3", "updated_at"}).
			AddRow(1, 1, 1, 0, 0, "", "", "", time.Now()),
	)

	// Mock language lookup (same language)
	mock.ExpectQuery("SELECT.*FROM languages WHERE code").WillReturnRows(
		sqlmock.NewRows([]string{"id", "code", "i18n_code", "name"}).
			AddRow(1, "en", "en", "English"),
	)

	// Update should NOT be called
	usRepo := repository.NewUserSettingsRepository(db)
	langRepo := repository.NewLanguageRepository(db)
	settingsService := service.NewUserSettingsService(usRepo, langRepo)

	lang := "en" // Same as current
	data := models.UserSettingsData{
		Language: &lang,
	}
	err = settingsService.Update(ctx, data)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserSettingsService_Update_UserNotAuthenticated(t *testing.T) {
	ctx := context.Background()
	// No user ID in context

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	usRepo := repository.NewUserSettingsRepository(db)
	langRepo := repository.NewLanguageRepository(db)
	settingsService := service.NewUserSettingsService(usRepo, langRepo)

	data := models.UserSettingsData{}
	err = settingsService.Update(ctx, data)

	if err == nil {
		t.Error("expected error when user not authenticated")
	}
}

func TestUserSettingsService_Update_GetByUserIdError(t *testing.T) {
	ctx := context.Background()
	ctx = contexthelper.SetUserId(ctx, 1)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mock user settings not found
	mock.ExpectQuery("SELECT.*FROM user_settings").WillReturnError(sql.ErrNoRows)

	usRepo := repository.NewUserSettingsRepository(db)
	langRepo := repository.NewLanguageRepository(db)
	settingsService := service.NewUserSettingsService(usRepo, langRepo)

	data := models.UserSettingsData{}
	err = settingsService.Update(ctx, data)

	if err == nil {
		t.Error("expected error when GetByUserId fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

