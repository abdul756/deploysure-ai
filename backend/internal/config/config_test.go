package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/config"
)

// setenv sets env variables for the duration of the test and restores them
// via t.Cleanup.
func setenv(t *testing.T, pairs ...string) {
	t.Helper()
	for i := 0; i+1 < len(pairs); i += 2 {
		key, val := pairs[i], pairs[i+1]
		prev, hadPrev := os.LookupEnv(key)
		if val == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, val)
		}
		t.Cleanup(func() {
			if hadPrev {
				os.Setenv(key, prev)
			} else {
				os.Unsetenv(key)
			}
		})
	}
}

func TestLoad_Defaults(t *testing.T) {
	// Ensure none of the env vars are set so we get pure defaults.
	setenv(t,
		"PORT", "",
		"REPORTS_DIR", "",
		"IBM_CLOUD_API_KEY", "",
		"WATSONX_API_KEY", "",
		"WATSONX_PROJECT_ID", "",
		"WATSONX_URL", "",
		"WATSONX_MODEL_ID", "",
		"READ_TIMEOUT_SEC", "",
		"WRITE_TIMEOUT_SEC", "",
		"IDLE_TIMEOUT_SEC", "",
		"SHUTDOWN_TIMEOUT_SEC", "",
	)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want %q", cfg.Port, "8080")
	}
	if cfg.ReportsDir != "reports" {
		t.Errorf("ReportsDir = %q, want %q", cfg.ReportsDir, "reports")
	}
	if cfg.ReadTimeout != 30*time.Second {
		t.Errorf("ReadTimeout = %v, want 30s", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout != 30*time.Second {
		t.Errorf("WriteTimeout = %v, want 30s", cfg.WriteTimeout)
	}
	if cfg.IdleTimeout != 60*time.Second {
		t.Errorf("IdleTimeout = %v, want 60s", cfg.IdleTimeout)
	}
	if cfg.ShutdownTimeout != 30*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 30s", cfg.ShutdownTimeout)
	}
}

func TestLoad_CustomPort(t *testing.T) {
	setenv(t, "PORT", "9090")
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9090")
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	cases := []string{"abc", "0", "99999", "-1"}
	for _, tc := range cases {
		t.Run(tc, func(t *testing.T) {
			setenv(t, "PORT", tc)
			_, err := config.Load()
			if err == nil {
				t.Fatalf("expected error for PORT=%q, got nil", tc)
			}
		})
	}
}

func TestLoad_InvalidTimeoutSec(t *testing.T) {
	setenv(t, "READ_TIMEOUT_SEC", "not-a-number")
	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for invalid READ_TIMEOUT_SEC")
	}
}

func TestLoad_ZeroTimeoutSec(t *testing.T) {
	setenv(t, "WRITE_TIMEOUT_SEC", "0")
	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for WRITE_TIMEOUT_SEC=0")
	}
}

func TestLoad_WatsonxCredentials(t *testing.T) {
	// Clear IBM_CLOUD_API_KEY so the legacy WATSONX_API_KEY fallback is used.
	setenv(t,
		"IBM_CLOUD_API_KEY", "",
		"WATSONX_API_KEY", "testkey",
		"WATSONX_PROJECT_ID", "proj-123",
	)
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.WatsonxAPIKey != "testkey" {
		t.Errorf("WatsonxAPIKey = %q, want %q", cfg.WatsonxAPIKey, "testkey")
	}
	if cfg.WatsonxProjectID != "proj-123" {
		t.Errorf("WatsonxProjectID mismatch")
	}
}

func TestLoad_IBMCloudAPIKeyTakesPriority(t *testing.T) {
	// IBM_CLOUD_API_KEY must win over WATSONX_API_KEY when both are set.
	setenv(t,
		"IBM_CLOUD_API_KEY", "ibm-key",
		"WATSONX_API_KEY", "legacy-key",
		"WATSONX_PROJECT_ID", "proj-456",
	)
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.WatsonxAPIKey != "ibm-key" {
		t.Errorf("WatsonxAPIKey = %q, want %q (IBM_CLOUD_API_KEY should take priority)", cfg.WatsonxAPIKey, "ibm-key")
	}
}
