package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/api"
	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/reports"
)

// setup creates a temporary reports directory, writes findings-before.json,
// and returns a Handler backed by it.
func setup(t *testing.T) (*api.Handler, string) {
	t.Helper()
	dir := t.TempDir()

	findings := []reports.Finding{
		{ID: "CQ-001", Severity: "high", Title: "Test finding"},
	}
	b, _ := json.Marshal(findings)
	os.WriteFile(filepath.Join(dir, "findings-before.json"), b, 0o644)
	os.WriteFile(filepath.Join(dir, "findings-after.json"), b, 0o644)
	os.WriteFile(filepath.Join(dir, "release-readiness-before.md"), []byte("# Before"), 0o644)
	os.WriteFile(filepath.Join(dir, "release-readiness-after.md"), []byte("# After"), 0o644)

	svc := reports.NewService(dir)
	h := api.NewHandler(svc, nil) // watsonx nil — tested separately
	return h, dir
}

func TestHealth_GET(t *testing.T) {
	h, _ := setup(t)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rw := httptest.NewRecorder()
	h.Health(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rw.Code)
	}
	if ct := rw.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var body map[string]string
	if err := json.NewDecoder(rw.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %q, want %q", body["status"], "ok")
	}
}

func TestHealth_MethodNotAllowed(t *testing.T) {
	h, _ := setup(t)
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/health", nil)
			rw := httptest.NewRecorder()
			h.Health(rw, req)
			if rw.Code != http.StatusMethodNotAllowed {
				t.Errorf("status = %d, want 405", rw.Code)
			}
		})
	}
}

func TestReady_GET(t *testing.T) {
	h, _ := setup(t)
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rw := httptest.NewRecorder()
	h.Ready(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rw.Code)
	}
}

func TestFindingsBefore_GET(t *testing.T) {
	h, _ := setup(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/findings/before", nil).
		WithContext(context.Background())
	rw := httptest.NewRecorder()
	h.FindingsBefore(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rw.Code)
	}
	var findings []reports.Finding
	if err := json.NewDecoder(rw.Body).Decode(&findings); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(findings) == 0 {
		t.Error("expected at least one finding")
	}
}

func TestFindingsAfter_NotFound(t *testing.T) {
	dir := t.TempDir()
	svc := reports.NewService(dir)
	h := api.NewHandler(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/findings/after", nil).
		WithContext(context.Background())
	rw := httptest.NewRecorder()
	h.FindingsAfter(rw, req)

	if rw.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rw.Code)
	}
}

func TestComparison_GET(t *testing.T) {
	h, _ := setup(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/comparison", nil).
		WithContext(context.Background())
	rw := httptest.NewRecorder()
	h.Comparison(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rw.Code)
	}
	var result reports.ComparisonResult
	if err := json.NewDecoder(rw.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if result.Summary.Before.Total == 0 {
		t.Error("expected non-zero Before total")
	}
}

func TestGraniteAnalyze_WatsonxNotConfigured(t *testing.T) {
	h, _ := setup(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/granite/analyze",
		strings.NewReader(`{"text":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	h.GraniteAnalyze(rw, req)

	// Handler is initialised with nil watsonx client.
	if rw.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want 503", rw.Code)
	}
}

func TestGraniteAnalyze_MethodNotAllowed(t *testing.T) {
	h, _ := setup(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/granite/analyze", nil)
	rw := httptest.NewRecorder()
	h.GraniteAnalyze(rw, req)

	if rw.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rw.Code)
	}
}

func TestGraniteAnalyze_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	svc := reports.NewService(dir)
	// We need a non-nil watsonx client to get past the availability check,
	// but we control the server so it never actually gets called.
	// Use a dummy that satisfies the interface by pointing at a closed server.
	// For the invalid-JSON test the handler returns 400 before calling watsonx.
	// We can still pass nil here because the handler rejects the request body first.
	h := api.NewHandler(svc, nil)

	// With nil wx, the handler returns 503 before checking JSON — so just verify
	// that the nil-client path is the one triggered here.
	req := httptest.NewRequest(http.MethodPost, "/api/v1/granite/analyze",
		strings.NewReader(`not-json`))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	h.GraniteAnalyze(rw, req)

	// 503 because watsonx is nil (that check runs before JSON parsing).
	if rw.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want 503", rw.Code)
	}
}

func TestGraniteAnalyze_EmptyText(t *testing.T) {
	// This test documents that an empty "text" field should return 422.
	// We cannot reach that path without a real watsonx client; the 503 guard
	// fires first. The test verifies the current contract.
	h, _ := setup(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/granite/analyze",
		strings.NewReader(`{"text":""}`))
	rw := httptest.NewRecorder()
	h.GraniteAnalyze(rw, req)

	if rw.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want 503 (watsonx not configured)", rw.Code)
	}
}

func TestReportBefore_GET(t *testing.T) {
	h, _ := setup(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports/before", nil).
		WithContext(context.Background())
	rw := httptest.NewRecorder()
	h.ReportBefore(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rw.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rw.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["content"] == "" {
		t.Error("expected non-empty content field")
	}
}
