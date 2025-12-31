package service_test

import (
	"backend/internal/repository"
	"backend/internal/service"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLanguageService_GetLanguages_Success(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mock language lookup
	mock.ExpectQuery("SELECT.*FROM languages").WillReturnRows(
		sqlmock.NewRows([]string{"id", "code", "name"}).
			AddRow(1, "en", "English").
			AddRow(2, "pl", "Polish"),
	)

	langRepo := repository.NewLanguageRepository(db)
	langService := service.NewLanguageService(langRepo)
	languages, err := langService.GetLanguages(ctx)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(languages) != 2 {
		t.Errorf("expected 2 languages, got %d", len(languages))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLanguageService_GetLanguages_EmptyList(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mock empty language list - rows.Next() returns false immediately
	rows := sqlmock.NewRows([]string{"id", "code", "name"})
	mock.ExpectQuery("SELECT.*FROM languages").WillReturnRows(rows)

	langRepo := repository.NewLanguageRepository(db)
	langService := service.NewLanguageService(langRepo)
	languages, err := langService.GetLanguages(ctx)

	// Service returns error when len(languages) == 0
	if err == nil {
		t.Error("expected error for empty language list")
	}
	if len(languages) != 0 {
		t.Errorf("expected empty list, got %d languages", len(languages))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLanguageService_GetLanguages_Error(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mock repository error
	mock.ExpectQuery("SELECT.*FROM languages").WillReturnError(sql.ErrNoRows)

	langRepo := repository.NewLanguageRepository(db)
	langService := service.NewLanguageService(langRepo)
	languages, err := langService.GetLanguages(ctx)

	if err == nil {
		t.Error("expected error when repository returns error")
	}
	if len(languages) != 0 {
		t.Errorf("expected empty list, got %d languages", len(languages))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

