package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ---------------------------------------------------------------------------
// OrdersHandler tests
// ---------------------------------------------------------------------------

func TestOrdersHandler_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
	w := httptest.NewRecorder()

	OrdersHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	var orders []Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if len(orders) == 0 {
		t.Fatal("expected at least one order in response")
	}
}

func TestOrdersHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/orders", nil)
	w := httptest.NewRecorder()

	OrdersHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

// SYNTHETIC DEFECT SD-07 note:
// SeedOrders() is an exported helper that intentionally has NO test below.
// This missing coverage is the seeded test-gap finding (runbook FR-09).
// A test for SeedOrders would look like:
//
//   func TestSeedOrders(t *testing.T) { ... }
//
// but it is deliberately omitted to be detected by the test-gap subagent.

// ---------------------------------------------------------------------------
// ReadinessHandler tests
// ---------------------------------------------------------------------------

func TestReadinessHandler_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	ReadinessHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body statusResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if body.Status != "ready" {
		t.Fatalf("expected status=ready, got %s", body.Status)
	}
}

func TestReadinessHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/ready", nil)
	w := httptest.NewRecorder()

	ReadinessHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// loggingMiddleware / NewRouter smoke test
// ---------------------------------------------------------------------------

func TestNewRouter_OrdersRoute(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from router, got %d", w.Code)
	}
}

func TestNewRouter_ReadyRoute(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from router, got %d", w.Code)
	}
}

// TestNewRouter_MethodNotAllowedViaRouter exercises the logging middleware's
// WriteHeader capture path by routing a non-GET request through the full
// middleware stack, which calls WriteHeader with a non-200 status.
func TestNewRouter_MethodNotAllowedViaRouter(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodDelete, "/api/orders", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 from router, got %d", w.Code)
	}
}
