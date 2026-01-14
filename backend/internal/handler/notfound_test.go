package handler_test

import (
	"backend/internal/handler"
	"encoding/json"
	"net/http"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestNotFoundHandler(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler()
	req, rr := NewTestRequest(
		http.MethodGet,
		"/nonexistent",
		nil,
		TestDeps{},
	)

	h.NotFoundHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("error decoding response: %v", err)
	}

	if response["code"] == nil {
		t.Error("expected error code in response")
	}
}

