// Package api provides the HTTP handlers for the DeploySure backend.
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/reports"
	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/watsonx"
)

// Handler holds the dependencies for all HTTP handlers.
type Handler struct {
	reports  *reports.Service
	watsonx  *watsonx.Client
}

// NewHandler constructs a Handler. Either dependency may be nil; methods check
// before use and return 503 Service Unavailable when unavailable.
func NewHandler(svc *reports.Service, wx *watsonx.Client) *Handler {
	return &Handler{reports: svc, watsonx: wx}
}

// Health handles GET /health – always returns 200 {"status":"ok"}.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready handles GET /ready – returns 200 when the service is ready.
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

// FindingsBefore handles GET /api/v1/findings/before.
func (h *Handler) FindingsBefore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	findings, err := h.reports.FindingsBefore(r.Context())
	if err != nil {
		handleReportError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, findings)
}

// FindingsAfter handles GET /api/v1/findings/after.
func (h *Handler) FindingsAfter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	findings, err := h.reports.FindingsAfter(r.Context())
	if err != nil {
		handleReportError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, findings)
}

// ReportBefore handles GET /api/v1/reports/before.
func (h *Handler) ReportBefore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	content, err := h.reports.ReportBefore(r.Context())
	if err != nil {
		handleReportError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"content": content})
}

// ReportAfter handles GET /api/v1/reports/after.
func (h *Handler) ReportAfter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	content, err := h.reports.ReportAfter(r.Context())
	if err != nil {
		handleReportError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"content": content})
}

// Comparison handles GET /api/v1/comparison.
func (h *Handler) Comparison(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	result, err := h.reports.Comparison(r.Context())
	if err != nil {
		handleReportError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// GraniteAnalyze handles POST /api/v1/granite/analyze.
func (h *Handler) GraniteAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if h.watsonx == nil {
		http.Error(w, "watsonx integration not configured", http.StatusServiceUnavailable)
		return
	}

	// Read and size-limit the request body.
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1 MiB
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	var req reports.GraniteAnalysisRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, fmt.Sprintf("invalid JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Text) == "" {
		http.Error(w, "field 'text' is required and must not be empty", http.StatusUnprocessableEntity)
		return
	}

	analysis, err := h.watsonx.Analyze(r.Context(), req.Text)
	if err != nil {
		log.Printf("level=error msg=\"granite analyze failed\" err=%q", err.Error())
		http.Error(w, "analysis failed", http.StatusBadGateway)
		return
	}

	writeJSON(w, http.StatusOK, reports.GraniteAnalysisResponse{Analysis: analysis})
}

// writeJSON marshals v to JSON and writes it with the given status code.
// If marshalling fails the response is a 500.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(b); err != nil {
		log.Printf("level=error msg=\"response write failed\" err=%q", err.Error())
	}
}

// handleReportError maps service-layer errors to appropriate HTTP responses.
func handleReportError(w http.ResponseWriter, err error) {
	if errors.Is(err, io.EOF) || isNotFound(err) {
		http.Error(w, "report not found", http.StatusNotFound)
		return
	}
	log.Printf("level=error msg=\"report load failed\" err=%q", err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// isNotFound returns true when err's message contains "file not found".
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "file not found")
}
