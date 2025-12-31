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
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestEmailChangeHandler_Success is skipped for unit tests because it requires RabbitMQ connection
// For full testing, use integration tests with a real RabbitMQ connection or proper mocking
func TestEmailChangeHandler_Success(t *testing.T) {
	t.Skip("Skipping test that requires RabbitMQ connection - use integration tests")
	
	// Setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Note: RabbitMQ connection is nil for unit tests
	// The service will log an error but continue (as per service code)
	// For full integration testing, a real RabbitMQ connection would be needed
	h := handler.NewHandler(db, nil)

	// Mock user lookup by ID (current user)
	regTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	confTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	mock.ExpectQuery("SELECT.*FROM users WHERE id").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "email", "password", "registered_at", "confirmed_at"}).
			AddRow(1, "testuser", "old@example.com", "hashed", regTime, confTime),
	)

	// Mock check if new email exists (should return no rows)
	mock.ExpectQuery("SELECT id FROM users WHERE email").WillReturnError(sql.ErrNoRows)

	// Mock check if email change token exists (should return no rows)
	mock.ExpectQuery("SELECT 1 FROM confirmation_tokens").WillReturnError(sql.ErrNoRows)

	// Mock confirmation token insert
	mock.ExpectExec("INSERT INTO confirmation_tokens").WillReturnResult(sqlmock.NewResult(1, 1))

	// Create request with user ID in context
	reqBody := map[string]string{
		"email": "newemail@example.com",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/email-change", bytes.NewBuffer(body))
	ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
	ctx = contexthelper.SetUserId(ctx, 1)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Execute
	h.EmailChangeHandler(rr, req)

	// Assert
	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEmailChangeHandler_InvalidJSON(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	req := httptest.NewRequest(http.MethodPost, "/email-change", bytes.NewBufferString("invalid json"))
	rr := httptest.NewRecorder()

	h.EmailChangeHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestEmailChangeHandler_EmptyEmail(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	reqBody := map[string]string{
		"email": "",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/email-change", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.EmailChangeHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestEmailChangeHandler_InvalidEmail(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	reqBody := map[string]string{
		"email": "invalid-email",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/email-change", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.EmailChangeHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

