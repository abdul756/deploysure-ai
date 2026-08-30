// Package handlers wires all HTTP routes onto a single ServeMux.
package handlers

import (
	"log"
	"net/http"
	"time"
)

// NewRouter constructs and returns the application ServeMux with all routes
// registered and a request-logging middleware applied.
//
// Routes registered:
//   - GET /api/orders
//   - GET /ready
//
// SYNTHETIC DEFECT SD-01 / runbook check DC-01:
// GET /health is required by docs/release-runbook.md section 3.1 and FR-32
// but is deliberately NOT registered here. This is a seeded hackathon
// demonstration defect to be detected by the docs-to-code consistency
// subagent.
func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/orders", OrdersHandler)
	mux.HandleFunc("/ready", ReadinessHandler)

	// NOTE: /health is intentionally omitted — see SYNTHETIC DEFECT comment above.

	return loggingMiddleware(mux)
}

// loggingMiddleware wraps an http.Handler and logs each request method, path,
// status code, and elapsed time using the standard log package.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		log.Printf("method=%s path=%s status=%d duration=%s",
			r.Method, r.URL.Path, lrw.statusCode, time.Since(start))
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture the status code
// written by downstream handlers so the logging middleware can record it.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before delegating to the underlying
// ResponseWriter.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
