package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// NewRouter constructs and returns the application's ServeMux.
// frontendDir is the path to the directory containing index.html, styles.css,
// and app.js. reportsDir is the directory from which report files are read.
func NewRouter(h *Handler, frontendDir string) http.Handler {
	mux := http.NewServeMux()

	// Health / readiness endpoints.
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/ready", h.Ready)

	// API routes.
	mux.HandleFunc("/api/v1/findings/before", h.FindingsBefore)
	mux.HandleFunc("/api/v1/findings/after", h.FindingsAfter)
	mux.HandleFunc("/api/v1/reports/before", h.ReportBefore)
	mux.HandleFunc("/api/v1/reports/after", h.ReportAfter)
	mux.HandleFunc("/api/v1/comparison", h.Comparison)
	mux.HandleFunc("/api/v1/granite/analyze", h.GraniteAnalyze)

	// Frontend – directory listing disabled; only known files are served.
	mux.HandleFunc("/styles.css", serveStaticFile(frontendDir, "styles.css", "text/css"))
	mux.HandleFunc("/app.js", serveStaticFile(frontendDir, "app.js", "application/javascript"))
	mux.HandleFunc("/", serveIndex(frontendDir))

	return loggingMiddleware(mux)
}

// serveIndex serves frontend/index.html for GET /.
// Any path other than "/" gets a 404; directory listings are never served.
func serveIndex(frontendDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		// Only serve the root path; anything else under "/" that hasn't matched
		// an explicit route is a 404.
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		indexPath := filepath.Join(frontendDir, "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, indexPath)
	}
}

// serveStaticFile serves a single known static file from the frontend directory.
// The filename is fixed at registration time — no user-supplied path is used.
func serveStaticFile(frontendDir, filename, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		// Use filepath.Base to strip any traversal attempts from the registered
		// filename constant (a belt-and-suspenders guard).
		safe := filepath.Base(filename)
		if safe != filename || strings.ContainsAny(filename, "/\\") {
			http.NotFound(w, r)
			return
		}
		fullPath := filepath.Join(frontendDir, safe)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", contentType)
		http.ServeFile(w, r, fullPath)
	}
}

// loggingMiddleware logs every request in a simple structured format.
// It does not log query parameters to avoid leaking sensitive data.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf(
			"level=info msg=\"http request\" method=%s path=%s status=%d duration=%s",
			r.Method,
			r.URL.Path,
			rw.status,
			fmt.Sprintf("%.3fms", float64(time.Since(start).Microseconds())/1000),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture the written status code.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
