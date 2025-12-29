package handler_test

import (
	"backend/internal/contexthelper"
	"backend/internal/handler"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type PingResponse struct {
	Data    map[string]interface{} `json:"data"`
	Message string                 `json:"message"`
	Code    int                    `json:"code"`
}

func TestPingHandler(t *testing.T) {
	// --- przygotowanie ---
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)

	// dodaj request ID do contextu, jak robi middleware
	ctx := context.WithValue(req.Context(), contexthelper.RequestIDKey, "test-id-123")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// --- wykonanie handlera ---
	h := handler.NewHandler(nil, nil) // przekazujemy nil, bo nie używamy bazy ani innych zależności w PingHandler
	h.PingHandler(rr, req)

	// --- asercje ---
	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	// Sprawdź body
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
