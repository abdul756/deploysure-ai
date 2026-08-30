// Package watsonx provides a client for the IBM watsonx.ai text generation API.
// IBM Cloud API keys and IAM tokens are never logged.
package watsonx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultModelID is the Granite model used when WATSONX_MODEL_ID is not set.
	DefaultModelID = "ibm/granite-13b-instruct-v2"
	// defaultMaxNewTokens is the upper bound on generated tokens.
	defaultMaxNewTokens = 1024
	// iamTokenURL is the IBM Cloud IAM token endpoint.
	iamTokenURL = "https://iam.cloud.ibm.com/identity/token"
)

// GraniteModelID is exported for callers that want to log the active model name.
// It is set to the default; the actual model sent is determined by NewClient.
const GraniteModelID = DefaultModelID

// iamToken caches the current IBM IAM access token and its expiry.
type iamToken struct {
	mu          sync.Mutex
	accessToken string    // never logged
	expiry      time.Time // refresh 60 s before expiry
}

// Client calls the watsonx.ai text-generation endpoint.
type Client struct {
	apiKey     string // IBM Cloud API key — never exposed in logs
	projectID  string
	modelID    string
	baseURL    string
	iamURL     string // injectable for tests
	httpClient *http.Client
	token      iamToken
}

// NewClient constructs a Client with the provided credentials and model.
// apiKey is the IBM Cloud API key read from IBM_CLOUD_API_KEY; it is never logged.
// modelID selects the Granite model; pass "" to use DefaultModelID.
func NewClient(apiKey, projectID, baseURL, modelID string) *Client {
	if modelID == "" {
		modelID = DefaultModelID
	}
	return &Client{
		apiKey:    apiKey,
		projectID: projectID,
		modelID:   modelID,
		baseURL:   baseURL,
		iamURL:    iamTokenURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ModelID returns the model identifier this client sends to watsonx.ai.
func (c *Client) ModelID() string { return c.modelID }

// WithIAMURL overrides the IBM IAM token endpoint used by this client.
// It is intended for tests that intercept IAM requests via httptest servers.
// Returns the receiver for chaining.
func (c *Client) WithIAMURL(u string) *Client {
	c.iamURL = u
	return c
}

// iamResponse is the JSON body returned by the IBM IAM token endpoint.
type iamResponse struct {
	AccessToken string `json:"access_token"` // never logged
	ExpiresIn   int    `json:"expires_in"`   // seconds
}

// accessToken returns a valid IAM bearer token, refreshing it when necessary.
// The token value is never written to any log.
func (c *Client) accessToken(ctx context.Context) (string, error) {
	c.token.mu.Lock()
	defer c.token.mu.Unlock()

	if c.token.accessToken != "" && time.Now().Before(c.token.expiry) {
		return c.token.accessToken, nil
	}

	form := url.Values{}
	form.Set("grant_type", "urn:ibm:params:oauth:grant-type:apikey")
	form.Set("apikey", c.apiKey) // sent in body, not logged

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.iamURL,
		strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("watsonx: iam request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("watsonx: iam http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limited, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		// Do NOT include the request body (contains apikey) in the error.
		return "", fmt.Errorf("watsonx: iam status %d: %s", resp.StatusCode, string(limited))
	}

	var tok iamResponse
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return "", fmt.Errorf("watsonx: iam decode: %w", err)
	}
	if tok.AccessToken == "" {
		return "", fmt.Errorf("watsonx: iam returned empty access_token")
	}

	ttl := time.Duration(tok.ExpiresIn) * time.Second
	if ttl <= 0 {
		ttl = 3600 * time.Second
	}
	c.token.accessToken = tok.AccessToken
	c.token.expiry = time.Now().Add(ttl - 60*time.Second) // refresh 60 s early
	return c.token.accessToken, nil
}

// generateRequest is the JSON body sent to /ml/v1/text/generation.
type generateRequest struct {
	ModelID    string                 `json:"model_id"`
	Input      string                 `json:"input"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	ProjectID  string                 `json:"project_id"`
}

// generateResponse is the JSON body returned by the watsonx endpoint.
type generateResponse struct {
	Results []struct {
		GeneratedText string `json:"generated_text"`
		StopReason    string `json:"stop_reason,omitempty"`
	} `json:"results"`
}

// Analyze sends findings text to Granite and returns the generated analysis.
// The returned string is the raw model output; callers may parse it further.
func (c *Client) Analyze(ctx context.Context, text string) (string, error) {
	bearer, err := c.accessToken(ctx)
	if err != nil {
		return "", err
	}

	body := generateRequest{
		ModelID:   c.modelID,
		Input:     buildPrompt(text),
		ProjectID: c.projectID,
		Parameters: map[string]interface{}{
			"max_new_tokens": defaultMaxNewTokens,
			"temperature":    0.2,
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("watsonx: marshal request: %w", err)
	}

	endpoint := c.baseURL + "/ml/v1/text/generation?version=2023-05-29"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("watsonx: create request: %w", err)
	}
	// Authorization header value (bearer) is an IAM token — never log it.
	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("watsonx: http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		limited, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return "", fmt.Errorf("watsonx: unexpected status %d: %s", resp.StatusCode, string(limited))
	}

	// Require valid JSON — reject plain-text or HTML error pages.
	var result generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("watsonx: decode response: %w", err)
	}
	if len(result.Results) == 0 || result.Results[0].GeneratedText == "" {
		return "", fmt.Errorf("watsonx: empty response from model")
	}
	return result.Results[0].GeneratedText, nil
}

// buildPrompt wraps user-supplied text in a structured Granite prompt.
func buildPrompt(input string) string {
	return fmt.Sprintf(`<|system|>
You are DeploySure AI, an expert in deployment risk analysis.
Analyze the following configuration or code for risks, anti-patterns, and security concerns.
Provide a clear, structured risk assessment with concrete recommendations.
<|user|>
%s
<|assistant|>`, input)
}
