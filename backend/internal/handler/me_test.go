package handler_test

import (
	"backend/internal/handler"
	"backend/internal/contexthelper"
	"database/sql"
	"encoding/json"
	"net/http"
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

	h := handler.NewHandler()

	// Mock user data lookup - GetDataById returns 11 columns
	mock.ExpectQuery("SELECT.*FROM users AS u").WillReturnRows(
		sqlmock.NewRows([]string{"u.id", "u.name", "email", "registered_at", "confirmed_at", "user_flags", "app_flags", "app_opt_1", "app_opt_2", "app_opt_3", "l.code"}).
			AddRow(1, "testuser", "test@example.com", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), uint64(0), uint64(0), "", "", "", "en"),
	)

	req, rr := NewTestRequest(
		http.MethodGet,
		"/me",
		nil,
		TestDeps{
			DB: db,
			UserID: 1,
			AccessTokenData: &contexthelper.AccessTokenData{
				UserId: 1,
				SetCookies: true,
			},
		},
	)

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
	if(!IsAccessCookieSet(resp)){
		t.Fatalf("expected access cookie to be set")
	}
}

func TestMeHandler_Unauthorized(t *testing.T) {

	h := handler.NewHandler()
	req, rr := NewTestRequest(
		http.MethodGet,
		"/me",
		nil,
		TestDeps{},
	)

	h.MeHandler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestMeHandler_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler()

	// Mock user not found
	mock.ExpectQuery("SELECT.*FROM users").WillReturnError(sql.ErrNoRows)

	req, rr := NewTestRequest(
		http.MethodGet,
		"/me",
		nil,
		TestDeps{DB: db, UserID: 1},
	)

	h.MeHandler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	resp := rr.Result()
	defer resp.Body.Close()
	if(IsAccessCookieSet(resp)){
		t.Fatalf("expected access cookie not to be set")
	}
}

