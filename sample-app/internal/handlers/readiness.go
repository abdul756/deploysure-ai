package handlers

import (
	"encoding/json"
	"net/http"
)

// statusResponse is the JSON body shape used by probe endpoints.
type statusResponse struct {
	Status string `json:"status"`
}

// ReadinessHandler handles GET /ready and returns {"status":"ready"} with
// HTTP 200 when the service is able to accept traffic.
func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statusResponse{Status: "ready"}) //nolint:errcheck
}
