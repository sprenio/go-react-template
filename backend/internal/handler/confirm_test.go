package handler_test

import (
	"backend/internal/contexthelper"
	"backend/internal/handler"
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

func TestConfirmHandler_RegisterToken_Success(t *testing.T) {
	// Setup
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
			AddRow(1, token, 1, "register", "{}", "NEW", time.Now().Add(1*time.Hour), time.Now(), time.Now()),
	)

	// Mock transaction begin
	mock.ExpectBegin()

	// Mock user insert (confirm registration creates new user)
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock user settings insert
	mock.ExpectExec("INSERT INTO user_settings").WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock transaction commit
	mock.ExpectCommit()

	// Mock token consume
	mock.ExpectExec("UPDATE confirmation_tokens").WillReturnResult(sqlmock.NewResult(0, 1))

	// Create request with token in URL
	req := httptest.NewRequest(http.MethodGet, "/confirm/"+token, nil)
	ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
	req = req.WithContext(ctx)

	// Set up chi router context with token
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("token", token)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	// Execute
	h.ConfirmHandler(rr, req)

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

func TestConfirmHandler_EmptyToken(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	req := httptest.NewRequest(http.MethodGet, "/confirm/", nil)
	ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
	req = req.WithContext(ctx)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("token", "")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	h.ConfirmHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestConfirmHandler_TokenNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	token := "invalid-token"

	// Mock confirmation token not found
	mock.ExpectQuery("SELECT.*FROM confirmation_tokens").WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/confirm/"+token, nil)
	ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
	req = req.WithContext(ctx)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("token", token)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	h.ConfirmHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestConfirmHandler_InvalidTokenType(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	token := "test-token-123"

	// Mock confirmation token with invalid type
	mock.ExpectQuery("SELECT.*FROM confirmation_tokens").WillReturnRows(
		sqlmock.NewRows([]string{"id", "token", "user_id", "type", "payload", "status", "expires_at", "status_changed_at", "created_at"}).
			AddRow(1, token, 1, "invalid_type", "{}", "NEW", time.Now().Add(1*time.Hour), time.Now(), time.Now()),
	)

	req := httptest.NewRequest(http.MethodGet, "/confirm/"+token, nil)
	ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
	req = req.WithContext(ctx)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("token", token)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	h.ConfirmHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

