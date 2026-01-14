package handler_test

import (
	"backend/internal/handler"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	amqp "github.com/rabbitmq/amqp091-go"
)

// TestResetPasswordHandler_Success is skipped for unit tests because it requires RabbitMQ connection
// For full testing, use integration tests with a real RabbitMQ connection or proper mocking
func TestResetPasswordHandler_Success(t *testing.T) {
	t.Skip("Skipping test that requires RabbitMQ connection - use integration tests")
	
	// Setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler()

	// Mock user lookup by email
	mock.ExpectQuery("SELECT.*FROM users").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "email", "password", "registered_at", "confirmed_at"}).
			AddRow(1, "testuser", "test@example.com", "hashed", "2024-01-01 00:00:00", "2024-01-01 00:00:00"),
	)

	// Mock confirmation token insert
	mock.ExpectExec("INSERT INTO confirmation_tokens").WillReturnResult(sqlmock.NewResult(1, 1))

	// Create request
	reqBody := map[string]string{
		"email": "test@example.com",
	}
	body, _ := json.Marshal(reqBody)
	req, rr := NewTestRequest(
		http.MethodPost,
		"/reset-password",
		bytes.NewBuffer(body),
		TestDeps{DB: db},
	)

	// Execute
	h.ResetPasswordHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestResetPasswordHandler_InvalidJSON(t *testing.T) {
	h := handler.NewHandler()
	req, rr := NewTestRequest(
		http.MethodPost,
		"/reset-password",
		bytes.NewBufferString("invalid json"),
		TestDeps{},
	)

	h.ResetPasswordHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestResetPasswordHandler_EmptyEmail(t *testing.T) {
	h := handler.NewHandler()

	reqBody := map[string]string{
		"email": "",
	}
	body, _ := json.Marshal(reqBody)
	req, rr := NewTestRequest(
		http.MethodPost,
		"/reset-password",
		bytes.NewBuffer(body),
		TestDeps{},
	)
	h.ResetPasswordHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestResetPasswordHandler_InvalidEmail(t *testing.T) {
	h := handler.NewHandler()

	reqBody := map[string]string{
		"email": "invalid-email",
	}
	body, _ := json.Marshal(reqBody)
	req, rr := NewTestRequest(
		http.MethodPost,
		"/reset-password",
		bytes.NewBuffer(body),
		TestDeps{},
	)

	h.ResetPasswordHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestResetPasswordHandler_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	rabbitConn := &amqp.Connection{}
	h := handler.NewHandler()

	// Mock user not found
	mock.ExpectQuery("SELECT.*FROM users").WillReturnError(sql.ErrNoRows)

	reqBody := map[string]string{
		"email": "nonexistent@example.com",
	}
	body, _ := json.Marshal(reqBody)
	req, rr := NewTestRequest(
		http.MethodPost,
		"/reset-password",
		bytes.NewBuffer(body),
		TestDeps{DB: db, RabbitConn:rabbitConn},
	)

	h.ResetPasswordHandler(rr, req)

	// Should still return 200 to prevent email enumeration
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

