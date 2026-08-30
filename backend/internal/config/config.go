// Package config loads server configuration from environment variables.
// No .env file is read; all values come from the process environment or
// the documented defaults below.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration for the DeploySure backend.
type Config struct {
	// Port is the TCP port the HTTP server listens on (default: 8080).
	Port string
	// ReportsDir is the directory from which report files are read (default: reports).
	ReportsDir string
	// WatsonxAPIKey is the IBM Cloud API key used to obtain IAM tokens. Never logged.
	// Read from IBM_CLOUD_API_KEY (preferred) or WATSONX_API_KEY (legacy fallback).
	WatsonxAPIKey string
	// WatsonxProjectID is the IBM watsonx.ai project ID.
	WatsonxProjectID string
	// WatsonxURL is the base URL for the watsonx.ai Inference API.
	WatsonxURL string
	// WatsonxModelID is the Granite model identifier sent to watsonx.ai.
	// Defaults to ibm/granite-13b-instruct-v2 when unset.
	WatsonxModelID string
	// ReadTimeout is the maximum duration for reading the entire request.
	ReadTimeout time.Duration
	// WriteTimeout is the maximum duration before timing out writes of the response.
	WriteTimeout time.Duration
	// IdleTimeout is the maximum amount of time to wait for the next request.
	IdleTimeout time.Duration
	// ShutdownTimeout is the grace period allowed to drain in-flight requests on shutdown.
	ShutdownTimeout time.Duration
}

// Load reads configuration from environment variables and returns a Config.
// Missing optional variables fall back to the documented defaults.
// Returns an error if any required variable is missing or any value is invalid.
func Load() (*Config, error) {
	port := envOr("PORT", "8080")
	if err := validatePort(port); err != nil {
		return nil, fmt.Errorf("config: invalid PORT %q: %w", port, err)
	}

	readTimeout, err := envDuration("READ_TIMEOUT_SEC", 30)
	if err != nil {
		return nil, fmt.Errorf("config: READ_TIMEOUT_SEC: %w", err)
	}
	writeTimeout, err := envDuration("WRITE_TIMEOUT_SEC", 30)
	if err != nil {
		return nil, fmt.Errorf("config: WRITE_TIMEOUT_SEC: %w", err)
	}
	idleTimeout, err := envDuration("IDLE_TIMEOUT_SEC", 60)
	if err != nil {
		return nil, fmt.Errorf("config: IDLE_TIMEOUT_SEC: %w", err)
	}
	shutdownTimeout, err := envDuration("SHUTDOWN_TIMEOUT_SEC", 30)
	if err != nil {
		return nil, fmt.Errorf("config: SHUTDOWN_TIMEOUT_SEC: %w", err)
	}

	return &Config{
		Port:             port,
		ReportsDir:       envOr("REPORTS_DIR", "reports"),
		WatsonxAPIKey:    ibmAPIKey(),                      // never log this value
		WatsonxProjectID: os.Getenv("WATSONX_PROJECT_ID"), // never log this value
		WatsonxURL:       envOr("WATSONX_URL", "https://us-south.ml.cloud.ibm.com"),
		WatsonxModelID:   os.Getenv("WATSONX_MODEL_ID"),
		ReadTimeout:      readTimeout,
		WriteTimeout:     writeTimeout,
		IdleTimeout:      idleTimeout,
		ShutdownTimeout:  shutdownTimeout,
	}, nil
}


// ibmAPIKey returns the IBM Cloud API key from IBM_CLOUD_API_KEY (preferred)
// or WATSONX_API_KEY (legacy). The value is never logged.
func ibmAPIKey() string {
	if v := os.Getenv("IBM_CLOUD_API_KEY"); v != "" {
		return v
	}
	return os.Getenv("WATSONX_API_KEY")
}


func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envDuration(key string, defaultSec int) (time.Duration, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return time.Duration(defaultSec) * time.Second, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("must be a positive integer, got %q", raw)
	}
	return time.Duration(n) * time.Second, nil
}

func validatePort(port string) error {
	n, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("not a number")
	}
	if n < 1 || n > 65535 {
		return fmt.Errorf("out of range 1-65535")
	}
	return nil
}
