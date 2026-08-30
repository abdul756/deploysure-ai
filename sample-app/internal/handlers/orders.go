// Package handlers contains the HTTP handler functions for the orders-api.
package handlers

import (
	"encoding/json"
	"net/http"
)

// Order represents a single order record returned by the API.
type Order struct {
	ID       string `json:"id"`
	Customer string `json:"customer"`
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}

// SeedOrders returns the static list of demo orders used by OrdersHandler.
// It is exported so that tests can verify the exact shape of the data without
// duplicating the fixture.
//
// SYNTHETIC DEFECT SD-07 / FR-09:
// SeedOrders is an exported helper function that intentionally has NO unit
// test in handlers_test.go. This seeds the test-gap finding for the hackathon
// demonstration. The missing test is the "one important handler error path"
// that is left uncovered.
func SeedOrders() []Order {
	return []Order{
		{ID: "ord-001", Customer: "Alice", Item: "Widget A", Quantity: 2, Status: "shipped"},
		{ID: "ord-002", Customer: "Bob", Item: "Gadget B", Quantity: 1, Status: "processing"},
		{ID: "ord-003", Customer: "Carol", Item: "Doohickey C", Quantity: 5, Status: "pending"},
	}
}

// OrdersHandler handles GET /api/orders and returns a JSON array of orders.
// Any non-GET method is rejected with HTTP 405 Method Not Allowed.
func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orders := SeedOrders()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
