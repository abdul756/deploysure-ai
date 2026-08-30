package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/api"
	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/reports"
)

// setupRouter creates a temporary reports dir and frontend dir, wires up the
// full router, and returns the router plus the dirs.
func setupRouter(t *testing.T) http.Handler {
	t.Helper()

	reportsDir := t.TempDir()
	frontendDir := t.TempDir()

	findings := []reports.Finding{{ID: "CQ-001", Severity: "high", Title: "Test"}}
	b, _ := json.Marshal(findings)
	os.WriteFile(filepath.Join(reportsDir, "findings-before.json"), b, 0o644)

	// Minimal frontend files.
	os.WriteFile(filepath.Join(frontendDir, "index.html"), []byte("<html></html>"), 0o644)
	os.WriteFile(filepath.Join(frontendDir, "styles.css"), []byte("body{}"), 0o644)
	os.WriteFile(filepath.Join(frontendDir, "app.js"), []byte("var x=1;"), 0o644)

	svc := reports.NewService(reportsDir)
	h := api.NewHandler(svc, nil)
	return api.NewRouter(h, frontendDir)
}

func TestRouter_Health(t *testing.T) {
	router := setupRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("/health: status = %d, want 200", rw.Code)
	}
}

func TestRouter_Ready(t *testing.T) {
	router := setupRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("/ready: status = %d, want 200", rw.Code)
	}
}

func TestRouter_FindingsBefore(t *testing.T) {
	router := setupRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/findings/before", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("/api/v1/findings/before: status = %d, want 200", rw.Code)
	}
}

func TestRouter_FindingsAfter_NotFound(t *testing.T) {
	router := setupRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/findings/after", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)
	if rw.Code != http.StatusNotFound {
		t.Errorf("/api/v1/findings/after: status = %d, want 404", rw.Code)
	}
}

func TestRouter_Comparison(t *testing.T) {
	router := setupRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/comparison", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("/api/v1/comparison: status = %d, want 200", rw.Code)
	}
}

func TestRouter_IndexHTML(t *testing.T) {
	router := setupRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("GET /: status = %d, want 200", rw.Code)
	}
	ct := rw.Header().Get("Content-Type")
	if ct == "" {
		t.Error("GET /: missing Content-Type")
	}
}

func TestRouter_StylesCSS(t *testing.T) {
	router := setupRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/styles.css", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("GET /styles.css: status = %d, want 200", rw.Code)
	}
}

func TestRouter_AppJS(t *testing.T) {
	router := setupRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("GET /app.js: status = %d, want 200", rw.Code)
	}
}

func TestRouter_UnknownPath_Returns404(t *testing.T) {
	router := setupRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/unknown-file.txt", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)
	if rw.Code != http.StatusNotFound {
		t.Errorf("unknown path: status = %d, want 404", rw.Code)
	}
}

func TestRouter_DirectoryListing_Disabled(t *testing.T) {
	router := setupRouter(t)
	// The catch-all "/" handler returns 404 for any path that is not "/".
	for _, path := range []string{"/frontend/", "/reports/", "/backend/"} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)
			if rw.Code == http.StatusOK {
				t.Errorf("%s: got 200, want non-200 (directory listing disabled)", path)
			}
		})
	}
}

func TestRouter_GraniteAnalyze_Registered(t *testing.T) {
	router := setupRouter(t)
	// GET /api/v1/granite/analyze should return 405 (only POST is allowed).
	req := httptest.NewRequest(http.MethodGet, "/api/v1/granite/analyze", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Errorf("/api/v1/granite/analyze GET: status = %d, want 405", rw.Code)
	}
}
