package handler_test

import (
	"backend/internal/handler"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

type PingResponse struct {
	Data    map[string]interface{} `json:"data"`
	Message string                 `json:"message"`
	Code    int                    `json:"code"`
}

func TestPingHandler(t *testing.T) {
	req, rr := NewTestRequest(
		http.MethodGet,
		"/ping",
		nil,
		TestDeps{
			RequestID: "test-id-123",
		},
	)
	h := handler.NewHandler()
	h.PingHandler(rr, req)

	// --- asserts ---
	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	// Check body
	var body PingResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("error decoding response: %v", err)
	}

	if body.Data["message"] != "pong strong" {
		t.Errorf("unexpected message: %v", body.Data["message"])
	}

	if body.Data["request_id"] != "test-id-123" {
		t.Errorf("unexpected request_id: %v", body.Data["request_id"])
	}
	if body.Code != 1000 {
		t.Errorf("unexpected code: %v", body.Code)
	}
}
