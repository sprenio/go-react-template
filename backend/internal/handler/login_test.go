package handler_test

import (
	"backend/internal/contexthelper"
	"backend/internal/handler"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func TestLoginHandler_Success(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	// Hash password for test
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Mock user lookup by email
	regTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	confTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	mock.ExpectQuery("SELECT.*FROM users").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "email", "password", "registered_at", "confirmed_at"}).
			AddRow(1, "testuser", "test@example.com", string(hashedPassword), regTime, confTime),
	)

	// Mock user data lookup - GetDataById returns 11 columns
	mock.ExpectQuery("SELECT.*FROM users AS u").WillReturnRows(
		sqlmock.NewRows([]string{"u.id", "u.name", "email", "registered_at", "confirmed_at", "user_flags", "app_flags", "app_opt_1", "app_opt_2", "app_opt_3", "l.code"}).
			AddRow(1, "testuser", "test@example.com", regTime, confTime, uint64(0), uint64(0), "", "", "", "en"),
	)

	// Create request
	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Execute
	h.LoginHandler(rr, req)

	// Assert
	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("error decoding response: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("invalid json"))
	rr := httptest.NewRecorder()

	h.LoginHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestLoginHandler_EmptyEmail(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	reqBody := map[string]string{
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.LoginHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestLoginHandler_EmptyPassword(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	reqBody := map[string]string{
		"email": "test@example.com",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.LoginHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	// Mock user not found
	mock.ExpectQuery("SELECT.*FROM users").WillReturnError(sql.ErrNoRows)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.LoginHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	// Hash password for test
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	// Mock user lookup by email
	regTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	confTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	mock.ExpectQuery("SELECT.*FROM users").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "email", "password", "registered_at", "confirmed_at"}).
			AddRow(1, "testuser", "test@example.com", string(hashedPassword), regTime, confTime),
	)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.LoginHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

