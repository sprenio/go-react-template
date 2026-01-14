package handler_test

import (
	"backend/internal/handler"
	"context"
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

func TestConfirmHandler(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		mock       func(sqlmock.Sqlmock)
		wantStatus int
	}{
		{
			name:  "success",
			token: "test-token",
			mock: func(m sqlmock.Sqlmock) {
				// Mock confirmation token lookup
				m.ExpectQuery("SELECT.*FROM confirmation_tokens").WillReturnRows(
					sqlmock.NewRows([]string{"id", "token", "user_id", "type", "payload", "status", "expires_at", "status_changed_at", "created_at"}).
						AddRow(1, "test-token", 1, "register", "{}", "NEW", time.Now().Add(1*time.Hour), time.Now(), time.Now()),
				)
				// Mock transaction begin
				m.ExpectBegin()
				// Mock user insert (confirm registration creates new user)
				m.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
				// Mock user settings insert
				m.ExpectExec("INSERT INTO user_settings").WillReturnResult(sqlmock.NewResult(1, 1))
				// Mock transaction commit
				m.ExpectCommit()
				// Mock token consume
				m.ExpectExec("UPDATE confirmation_tokens").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "token not found",
			token: "invalid-token",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT.*FROM confirmation_tokens").WillReturnError(sql.ErrNoRows)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:  "invalid token type",
			token: "invalid-token",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT.*FROM confirmation_tokens").WillReturnRows(
					sqlmock.NewRows([]string{"id", "token", "user_id", "type", "payload", "status", "expires_at", "status_changed_at", "created_at"}).
						AddRow(1, "invalid-token", 1, "invalid_type", "{}", "NEW", time.Now().Add(1*time.Hour), time.Now(), time.Now()),
					)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty token",
			token:      "",
			mock:       func(m sqlmock.Sqlmock) {},
			wantStatus: http.StatusNotFound,
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
				"/confirm/"+tt.token,
				nil,
				TestDeps{DB: db},
			)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("token", tt.token)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			handler.NewHandler().ConfirmHandler(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("expected %d, got %d", tt.wantStatus, rr.Code)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
