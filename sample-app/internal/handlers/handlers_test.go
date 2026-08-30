package handlers

import (
	"encoding/json"
	"fmt"
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
	if len(orders) != 3 {
		t.Fatalf("expected 3 orders in response, got %d", len(orders))
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

// failingWriter is a test helper that wraps httptest.ResponseRecorder but
// makes Write() always return an error, exercising the json.Encoder error path.
type failingWriter struct {
	http.ResponseWriter
}

func (f *failingWriter) Write([]byte) (int, error) {
	return 0, fmt.Errorf("disk full")
}

func TestOrdersHandler_EncoderError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
	rec := httptest.NewRecorder()
	fw := &failingWriter{ResponseWriter: rec}

	OrdersHandler(fw, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// ---------------------------------------------------------------------------
// SeedOrders tests  (REM-010 · TC-001 / DC-010)
// ---------------------------------------------------------------------------

func TestSeedOrders(t *testing.T) {
	orders := SeedOrders()
	if got, want := len(orders), 3; got != want {
		t.Fatalf("SeedOrders() len = %d, want %d", got, want)
	}
	if orders[0].ID != "ord-001" {
		t.Errorf("orders[0].ID = %q, want %q", orders[0].ID, "ord-001")
	}
	if orders[1].ID != "ord-002" {
		t.Errorf("orders[1].ID = %q, want %q", orders[1].ID, "ord-002")
	}
	if orders[2].ID != "ord-003" {
		t.Errorf("orders[2].ID = %q, want %q", orders[2].ID, "ord-003")
	}
}

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

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
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
// HealthHandler tests  (REM-024 · TC-003)
// ---------------------------------------------------------------------------

func TestHealthHandler_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	HealthHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status=ok, got %s", body["status"])
	}
}

func TestHealthHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	w := httptest.NewRecorder()

	HealthHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

// ---------------------------------------------------------------------------
// loggingMiddleware / NewRouter smoke tests
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

func TestNewRouter_HealthRoute(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from router /health, got %d", w.Code)
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
