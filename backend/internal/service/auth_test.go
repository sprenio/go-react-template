package service_test

import (
	"backend/internal/repository"
	"backend/internal/service"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login_Success(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Hash password for test
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	regTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	confTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Mock user lookup by email
	mock.ExpectQuery("SELECT.*FROM users WHERE email").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "email", "password", "registered_at", "confirmed_at"}).
			AddRow(1, "testuser", "test@example.com", string(hashedPassword), regTime, confTime),
	)

	// Mock user data lookup
	mock.ExpectQuery("SELECT.*FROM users AS u").WillReturnRows(
		sqlmock.NewRows([]string{"u.id", "u.name", "email", "registered_at", "confirmed_at", "user_flags", "app_flags", "app_opt_1", "app_opt_2", "app_opt_3", "l.code"}).
			AddRow(1, "testuser", "test@example.com", regTime, confTime, uint64(0), uint64(0), "", "", "", "en"),
	)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	user, err := authService.Login(ctx, "test@example.com", "password123")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if user.Id != 1 {
		t.Errorf("expected user ID 1, got %d", user.Id)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Mock user not found
	mock.ExpectQuery("SELECT.*FROM users WHERE email").WillReturnError(sql.ErrNoRows)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	_, err = authService.Login(ctx, "nonexistent@example.com", "password123")

	if err == nil {
		t.Error("expected error for non-existent user")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Hash password for test
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	regTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	confTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Mock user lookup by email
	mock.ExpectQuery("SELECT.*FROM users WHERE email").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "email", "password", "registered_at", "confirmed_at"}).
			AddRow(1, "testuser", "test@example.com", string(hashedPassword), regTime, confTime),
	)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	_, err = authService.Login(ctx, "test@example.com", "wrongpassword")

	if err == nil {
		t.Error("expected error for invalid password")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAuthService_Login_GetDataByIdError(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Hash password for test
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	regTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	confTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Mock user lookup by email
	mock.ExpectQuery("SELECT.*FROM users WHERE email").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "email", "password", "registered_at", "confirmed_at"}).
			AddRow(1, "testuser", "test@example.com", string(hashedPassword), regTime, confTime),
	)

	// Mock GetDataById error
	mock.ExpectQuery("SELECT.*FROM users AS u").WillReturnError(sql.ErrNoRows)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	_, err = authService.Login(ctx, "test@example.com", "password123")

	if err == nil {
		t.Error("expected error when GetDataById fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

