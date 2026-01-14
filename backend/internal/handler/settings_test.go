package handler_test

import (
	"backend/internal/handler"
	"backend/internal/models"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSettingsHandler_Success(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler()

	// Mock user settings lookup (current language is 1 = en)
	mock.ExpectQuery("SELECT.*FROM user_settings").WillReturnRows(
		sqlmock.NewRows([]string{"id", "user_id", "lang_id", "user_flags", "app_flags", "app_opt_1", "app_opt_2", "app_opt_3", "updated_at"}).
			AddRow(1, 1, 1, 0, 0, "", "", "", time.Now()),
	)

	// Mock language lookup (changing to pl, which has different ID)
	mock.ExpectQuery("SELECT.*FROM languages WHERE code").WillReturnRows(
		sqlmock.NewRows([]string{"id", "code", "i18n_code", "name"}).
			AddRow(2, "pl", "pl", "Polish"),
	)

	// Mock user settings update (language changed, so update is needed)
	mock.ExpectExec("UPDATE user_settings").WillReturnResult(sqlmock.NewResult(0, 1))

	// Create request
	lang := "pl"
	reqBody := models.UserSettingsData{
		Language: &lang,
	}
	body, _ := json.Marshal(reqBody)
	req, rr := NewTestRequest(
		http.MethodPost,
		"/settings",
		bytes.NewBuffer(body),
		TestDeps{DB: db, UserID: 1},
	)

	// Execute
	h.SettingsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSettingsHandler_InvalidJSON(t *testing.T) {
	h := handler.NewHandler()
	req, rr := NewTestRequest(
		http.MethodPost,
		"/settings",
		bytes.NewBufferString("invalid json"),
		TestDeps{},
	)
	h.SettingsHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestSettingsHandler_InvalidAppOpt2(t *testing.T) {

	h := handler.NewHandler()

	invalidOpt := models.AppOption2("INVALID_OPTION")
	reqBody := models.UserSettingsData{
		AppOpt2: &invalidOpt,
	}
	body, _ := json.Marshal(reqBody)
	req, rr := NewTestRequest(
		http.MethodPost,
		"/settings",
		bytes.NewBuffer(body),
		TestDeps{UserID: 1},
	)

	h.SettingsHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestSettingsHandler_ValidAppOpt2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler()

	// Mock user settings lookup
	mock.ExpectQuery("SELECT.*FROM user_settings").WillReturnRows(
		sqlmock.NewRows([]string{"id", "user_id", "lang_id", "user_flags", "app_flags", "app_opt_1", "app_opt_2", "app_opt_3", "updated_at"}).
			AddRow(1, 1, 1, 0, 0, "", "", "", time.Now()),
	)

	// Mock user settings update
	mock.ExpectExec("UPDATE user_settings").WillReturnResult(sqlmock.NewResult(0, 1))

	optA := models.APP_OPTION_A
	reqBody := models.UserSettingsData{
		AppOpt2: &optA,
	}
	body, _ := json.Marshal(reqBody)
	req, rr := NewTestRequest(
		http.MethodPost,
		"/settings",
		bytes.NewBuffer(body),
		TestDeps{DB: db, UserID: 1},
	)
	h.SettingsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSettingsHandler_ServiceError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler()

	// Mock user settings lookup error
	mock.ExpectQuery("SELECT.*FROM user_settings").WillReturnError(sql.ErrNoRows)

	lang := "en"
	reqBody := models.UserSettingsData{
		Language: &lang,
	}
	body, _ := json.Marshal(reqBody)
	req, rr := NewTestRequest(
		http.MethodPost,
		"/settings",
		bytes.NewBuffer(body),
		TestDeps{DB: db, UserID: 1},
	)
	h.SettingsHandler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

