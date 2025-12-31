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

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCfgHandler_Success(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	// Mock language lookup
	mock.ExpectQuery("SELECT.*FROM languages").WillReturnRows(
		sqlmock.NewRows([]string{"id", "code", "name"}).
			AddRow(1, "en", "English").
			AddRow(2, "pl", "Polish"),
	)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/cfg", nil)
	ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Execute
	h.CfgHandler(rr, req)

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

func TestCfgHandler_LanguagesError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler(db, nil)

	// Mock language lookup error
	mock.ExpectQuery("SELECT.*FROM languages").WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/cfg", nil)
	ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	h.CfgHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

