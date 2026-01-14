package handler_test

import (
	"backend/internal/handler"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMethodNotAllowedHandler(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	h := handler.NewHandler()

	req, rr := NewTestRequest(
		http.MethodDelete,
		"/some-endpoint",
		nil,
		TestDeps{},
	)

	h.MethodNotAllowedHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("error decoding response: %v", err)
	}

	if response["code"] == nil {
		t.Error("expected error code in response")
	}
}

