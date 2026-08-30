package watsonx_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/watsonx"
)

// iamToken is a minimal valid IAM token response used by all test servers.
const testToken = "test-iam-token"

// iamHandler writes a synthetic IAM token response.
func iamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	resp := map[string]interface{}{
		"access_token": testToken,
		"expires_in":   3600,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp) //nolint:errcheck
}

// newTestClient returns a Client whose IAM endpoint and watsonx base URL both
// point at the provided test server URL.
func newTestClient(serverURL string) *watsonx.Client {
	return watsonx.NewClient("test-api-key", "test-project-id", serverURL, "").
		WithIAMURL(serverURL + "/identity/token")
}

// mux returns a ServeMux that routes IAM token requests and generation requests
// to the supplied generation handler.
func mux(genHandler http.HandlerFunc) *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/identity/token", iamHandler)
	m.HandleFunc("/ml/v1/text/generation", genHandler)
	return m
}

func TestAnalyze_Success(t *testing.T) {
	srv := httptest.NewServer(mux(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		// Verify Authorization header is present (do not log its value).
		if r.Header.Get("Authorization") == "" {
			t.Error("Authorization header missing")
		}
		// Return a minimal valid response.
		resp := map[string]interface{}{
			"results": []map[string]string{
				{"generated_text": "This is the analysis result."},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp) //nolint:errcheck
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	got, err := client.Analyze(context.Background(), "deploy: latest")
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if got != "This is the analysis result." {
		t.Errorf("analysis = %q, unexpected", got)
	}
}

func TestAnalyze_HTTPError(t *testing.T) {
	srv := httptest.NewServer(mux(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.Analyze(context.Background(), "some text")
	if err == nil {
		t.Fatal("expected error for 401, got nil")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("error %q should mention 401", err.Error())
	}
}

func TestAnalyze_EmptyResponse(t *testing.T) {
	srv := httptest.NewServer(mux(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{"results": []map[string]string{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp) //nolint:errcheck
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.Analyze(context.Background(), "some text")
	if err == nil {
		t.Fatal("expected error for empty results")
	}
}

func TestAnalyze_ContextCancelled(t *testing.T) {
	srv := httptest.NewServer(mux(func(w http.ResponseWriter, r *http.Request) {
		// Block until the client cancels.
		<-r.Context().Done()
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before the request is made

	client := newTestClient(srv.URL)
	_, err := client.Analyze(ctx, "some text")
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestAnalyze_BadJSON(t *testing.T) {
	srv := httptest.NewServer(mux(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json")) //nolint:errcheck
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.Analyze(context.Background(), "some text")
	if err == nil {
		t.Fatal("expected error for bad JSON response")
	}
}

func TestAnalyze_IAMError(t *testing.T) {
	// Server returns 401 for all requests, including IAM.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "forbidden", http.StatusForbidden)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.Analyze(context.Background(), "some text")
	if err == nil {
		t.Fatal("expected error when IAM returns non-200")
	}
	if !strings.Contains(err.Error(), "iam") {
		t.Errorf("error %q should mention iam", err.Error())
	}
}

func TestModelID_Default(t *testing.T) {
	c := watsonx.NewClient("key", "proj", "http://example.com", "")
	if c.ModelID() != watsonx.DefaultModelID {
		t.Errorf("ModelID() = %q, want %q", c.ModelID(), watsonx.DefaultModelID)
	}
}

func TestModelID_Custom(t *testing.T) {
	const custom = "ibm/granite-3b-code-instruct"
	c := watsonx.NewClient("key", "proj", "http://example.com", custom)
	if c.ModelID() != custom {
		t.Errorf("ModelID() = %q, want %q", c.ModelID(), custom)
	}
}
