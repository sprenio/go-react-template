package service_test

import (
	"backend/internal/repository"
	"backend/internal/service"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserService_GetUserResponseData_Success(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	regTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	confTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Mock user data lookup
	mock.ExpectQuery("SELECT.*FROM users AS u").WillReturnRows(
		sqlmock.NewRows([]string{"u.id", "u.name", "email", "registered_at", "confirmed_at", "user_flags", "app_flags", "app_opt_1", "app_opt_2", "app_opt_3", "l.code"}).
			AddRow(1, "testuser", "test@example.com", regTime, confTime, uint64(0), uint64(0), "", "", "", "en"),
	)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	user, err := userService.GetUserResponseData(ctx, 1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if user.Id != 1 {
		t.Errorf("expected user ID 1, got %d", user.Id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserService_GetUserResponseData_UserNotFound(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mock user not found
	mock.ExpectQuery("SELECT.*FROM users AS u").WillReturnError(sql.ErrNoRows)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	_, err = userService.GetUserResponseData(ctx, 999)

	if err == nil {
		t.Error("expected error for non-existent user")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserService_GetUserResponseData_EmptyId(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mock returns error (simulating user not found which results in ID = 0)
	mock.ExpectQuery("SELECT.*FROM users AS u").WillReturnError(sql.ErrNoRows)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	_, err = userService.GetUserResponseData(ctx, 0)

	if err == nil {
		t.Error("expected error for user with ID 0")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUserService_ConfirmEmailChangeToken_Success would require constructing
// a models.ConfirmationToken with proper payload structure.
// This is better tested as an integration test with actual token creation.

