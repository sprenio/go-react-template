package handler_test

import (
	"backend/internal/handler"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestConfigHandler(t *testing.T) {
	tests := []struct {
		name       string
		mock       func(sqlmock.Sqlmock)
		wantStatus int
	}{
		{
			name: "success",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT.*FROM languages").WillReturnRows(
					sqlmock.NewRows([]string{"id", "code", "name"}).
						AddRow(1, "en", "English").
						AddRow(2, "pl", "Polish"),
				)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "lang errors",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT.*FROM languages").WillReturnError(sql.ErrNoRows)
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			tt.mock(mock)

			req, rr := NewTestRequest(
				http.MethodGet,
				"/cfg",
				nil,
				TestDeps{DB: db},
			)

			handler.NewHandler().CfgHandler(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("expected %d, got %d", tt.wantStatus, rr.Code)
			}
			if rr.Code == http.StatusOK {
				resp := rr.Result()
				defer resp.Body.Close()
				var response map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					t.Fatalf("error decoding response: %v", err)
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
