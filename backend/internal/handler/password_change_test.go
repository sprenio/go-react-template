package handler_test

import (
	"backend/internal/handler"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
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

	h := handler.NewHandler()

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
	req, rr := NewTestRequest(
		http.MethodPost,
		"/password-change/"+token,
		bytes.NewBuffer(body),
		TestDeps{DB: db},
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("token", token)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Execute
	h.PasswordChangeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPasswordChangeHandler_EmptyToken(t *testing.T) {
	h := handler.NewHandler()

	reqBody := map[string]string{
		"password": "NewPassword123!@#",
	}
	body, _ := json.Marshal(reqBody)
	req, rr := NewTestRequest(
		http.MethodPost,
		"/password-change/",
		bytes.NewBuffer(body),
		TestDeps{},
	)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("token", "")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.PasswordChangeHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

func TestPasswordChangeHandler_InvalidJSON(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler()

	token := "test-token-123"

	// Mock confirmation token lookup (handler checks token before parsing JSON)
	mock.ExpectQuery("SELECT.*FROM confirmation_tokens").WillReturnRows(
		sqlmock.NewRows([]string{"id", "token", "user_id", "type", "payload", "status", "expires_at", "status_changed_at", "created_at"}).
			AddRow(1, token, 1, "password_change", "{}", "NEW", time.Now().Add(1*time.Hour), time.Now(), time.Now()),
	)
	req, rr := NewTestRequest(
		http.MethodPost,
		"/password-change/"+token,
		bytes.NewBufferString("invalid json"),
		TestDeps{DB:db},
	)

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

	h := handler.NewHandler()

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
	req, rr := NewTestRequest(
		http.MethodPost,
		"/password-change/"+token,
		bytes.NewBuffer(body),
		TestDeps{DB:db},
	)

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

	h := handler.NewHandler()

	token := "invalid-token"

	// Mock confirmation token not found
	mock.ExpectQuery("SELECT.*FROM confirmation_tokens").WillReturnError(sql.ErrNoRows)

	reqBody := map[string]string{
		"password": "NewPassword123!@#",
	}
	body, _ := json.Marshal(reqBody)

	req, rr := NewTestRequest(
		http.MethodPost,
		"/password-change/"+token,
		bytes.NewBuffer(body),
		TestDeps{DB:db},
	)
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

