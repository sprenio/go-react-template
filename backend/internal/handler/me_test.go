package handler_test

import (
	"backend/internal/contexthelper"
	"backend/internal/handler"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMeHandler_Success(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	// Mock user data lookup - GetDataById returns 11 columns
	mock.ExpectQuery("SELECT.*FROM users AS u").WillReturnRows(
		sqlmock.NewRows([]string{"u.id", "u.name", "email", "registered_at", "confirmed_at", "user_flags", "app_flags", "app_opt_1", "app_opt_2", "app_opt_3", "l.code"}).
			AddRow(1, "testuser", "test@example.com", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), uint64(0), uint64(0), "", "", "", "en"),
	)

	// Create request with user ID in context
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
	ctx = contexthelper.SetUserId(ctx, 1)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Execute
	h.MeHandler(rr, req)

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

func TestMeHandler_Unauthorized(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	// Create request without user ID in context
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	h.MeHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestMeHandler_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	// Mock user not found
	mock.ExpectQuery("SELECT.*FROM users").WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
	ctx = contexthelper.SetUserId(ctx, 1)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	h.MeHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

