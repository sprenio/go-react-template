package handler_test

import (
	"backend/internal/handler"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestRegisterHandler_Success is skipped for unit tests because it requires RabbitMQ connection
// For full testing, use integration tests with a real RabbitMQ connection or proper mocking
func TestRegisterHandler_Success(t *testing.T) {
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
	h := handler.NewHandler()

	// Mock database queries
	// Check if user exists by email or name (should return no rows - sql.ErrNoRows)
	mock.ExpectQuery("SELECT id FROM users WHERE email").WillReturnError(sql.ErrNoRows)

	// Check if confirmation token exists (should return no rows)
	mock.ExpectQuery("SELECT 1 FROM confirmation_tokens").WillReturnError(sql.ErrNoRows)

	// Get language by code (4 columns: id, code, i18n_code, name)
	mock.ExpectQuery("SELECT.*FROM languages WHERE code").WillReturnRows(
		sqlmock.NewRows([]string{"id", "code", "i18n_code", "name"}).
			AddRow(1, "en", "en", "English"),
	)

	// Insert confirmation token (register doesn't insert user, only token)
	mock.ExpectExec("INSERT INTO confirmation_tokens").WillReturnResult(sqlmock.NewResult(1, 1))

	// Create request
	reqBody := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "Test123!@#",
		"language": "en",
	}
	body, _ := json.Marshal(reqBody)
	req, rr := NewTestRequest(
		http.MethodPost,
		"/register",
		bytes.NewBuffer(body),
		TestDeps{DB: db},
	)

	// Execute
	h.RegisterHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRegisterHandler_InvalidMethod(t *testing.T) {
	h := handler.NewHandler()

	req, rr := NewTestRequest(
		http.MethodGet,
		"/register",
		nil,
		TestDeps{},
	)

	h.RegisterHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	h := handler.NewHandler()
	req, rr := NewTestRequest(
		http.MethodPost,
		"/register",
		bytes.NewBufferString("invalid json"),
		TestDeps{},
	)

	h.RegisterHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestRegisterHandler_MissingUsername(t *testing.T) {
	h := handler.NewHandler()

	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "Test123!@#",
	}
	body, _ := json.Marshal(reqBody)
	req, rr := NewTestRequest(
		http.MethodPost,
		"/register",
		bytes.NewBuffer(body),
		TestDeps{},
	)

	h.RegisterHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	if response["code"] == nil {
		t.Error("expected error code in response")
	}
}

func TestRegisterHandler_MissingEmail(t *testing.T) {
	h := handler.NewHandler()

	reqBody := map[string]string{
		"username": "testuser",
		"password": "Test123!@#",
	}
	body, _ := json.Marshal(reqBody)
	req, rr := NewTestRequest(
		http.MethodPost,
		"/register",
		bytes.NewBuffer(body),
		TestDeps{},
	)

	h.RegisterHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestRegisterHandler_MissingPassword(t *testing.T) {
	h := handler.NewHandler()

	reqBody := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
	}
	body, _ := json.Marshal(reqBody)
	req, rr := NewTestRequest(
		http.MethodPost,
		"/register",
		bytes.NewBuffer(body),
		TestDeps{},
	)

	h.RegisterHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}
