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
	"github.com/go-chi/chi/v5"
)

func TestPasswordChangeHandler_Success(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	token := "test-token-123"
	newPassword := "NewPassword123!@#"

	// Mock confirmation token lookup
	mock.ExpectQuery("SELECT.*FROM confirmation_tokens").WillReturnRows(
		sqlmock.NewRows([]string{"id", "token", "user_id", "type", "payload", "status", "expires_at", "status_changed_at", "created_at"}).
			AddRow(1, token, 1, "password_change", "{}", "NEW", time.Now().Add(1*time.Hour), time.Now(), time.Now()),
	)

	// Mock password update
	mock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(0, 1))

	// Mock token consume
	mock.ExpectExec("UPDATE confirmation_tokens").WillReturnResult(sqlmock.NewResult(0, 1))

	// Create request
	reqBody := map[string]string{
		"password": newPassword,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/password-change/"+token, bytes.NewBuffer(body))
	ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
	req = req.WithContext(ctx)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("token", token)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	// Execute
	h.PasswordChangeHandler(rr, req)

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

func TestPasswordChangeHandler_EmptyToken(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	reqBody := map[string]string{
		"password": "NewPassword123!@#",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/password-change/", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("token", "")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.PasswordChangeHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestPasswordChangeHandler_InvalidJSON(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	token := "test-token-123"

	// Mock confirmation token lookup (handler checks token before parsing JSON)
	mock.ExpectQuery("SELECT.*FROM confirmation_tokens").WillReturnRows(
		sqlmock.NewRows([]string{"id", "token", "user_id", "type", "payload", "status", "expires_at", "status_changed_at", "created_at"}).
			AddRow(1, token, 1, "password_change", "{}", "NEW", time.Now().Add(1*time.Hour), time.Now(), time.Now()),
	)

	req := httptest.NewRequest(http.MethodPost, "/password-change/"+token, bytes.NewBufferString("invalid json"))
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("token", token)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.PasswordChangeHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPasswordChangeHandler_InvalidPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	token := "test-token-123"

	// Mock confirmation token lookup
	mock.ExpectQuery("SELECT.*FROM confirmation_tokens").WillReturnRows(
		sqlmock.NewRows([]string{"id", "token", "user_id", "type", "payload", "status", "expires_at", "status_changed_at", "created_at"}).
			AddRow(1, token, 1, "password_change", "{}", "NEW", time.Now().Add(1*time.Hour), time.Now(), time.Now()),
	)

	reqBody := map[string]string{
		"password": "weak",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/password-change/"+token, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("token", token)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.PasswordChangeHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPasswordChangeHandler_TokenNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	token := "invalid-token"

	// Mock confirmation token not found
	mock.ExpectQuery("SELECT.*FROM confirmation_tokens").WillReturnError(sql.ErrNoRows)

	reqBody := map[string]string{
		"password": "NewPassword123!@#",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/password-change/"+token, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("token", token)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.PasswordChangeHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

