// Command server starts the DeploySure AI HTTP server.
//
// Run from the repository root:
//
//	go run ./backend/cmd/server
//
// Configuration is read exclusively from environment variables — no .env file.
//
// Environment variables:
//
//	PORT                 TCP port (default: 8080, range 1-65535)
//	REPORTS_DIR          Directory containing JSON/markdown report files (default: reports)
//	IBM_CLOUD_API_KEY    IBM Cloud API key used for IAM token exchange (never logged)
//	WATSONX_API_KEY      Legacy alias for IBM_CLOUD_API_KEY (never logged)
//	WATSONX_PROJECT_ID   IBM watsonx.ai project ID (never logged)
//	WATSONX_URL          watsonx.ai base URL (default: https://us-south.ml.cloud.ibm.com)
//	WATSONX_MODEL_ID     Granite model ID (default: ibm/granite-13b-instruct-v2)
//	READ_TIMEOUT_SEC     HTTP read timeout seconds (default: 30)
//	WRITE_TIMEOUT_SEC    HTTP write timeout seconds (default: 30)
//	IDLE_TIMEOUT_SEC     HTTP idle timeout seconds (default: 60)
//	SHUTDOWN_TIMEOUT_SEC Graceful-shutdown grace period seconds (default: 30)
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/api"
	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/config"
	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/reports"
	"github.com/YOUR-USERNAME/deploysure-ai/backend/internal/watsonx"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("level=fatal msg=\"invalid configuration\" err=%q", err.Error())
	}

	svc := reports.NewService(cfg.ReportsDir)

	var wx *watsonx.Client
	if cfg.WatsonxAPIKey != "" && cfg.WatsonxProjectID != "" {
		wx = watsonx.NewClient(cfg.WatsonxAPIKey, cfg.WatsonxProjectID, cfg.WatsonxURL, cfg.WatsonxModelID)
		log.Printf("level=info msg=\"watsonx integration enabled\" url=%s", cfg.WatsonxURL)
	} else {
		log.Printf("level=warn msg=\"watsonx integration disabled\" reason=\"WATSONX_API_KEY or WATSONX_PROJECT_ID not set\"")
	}

	h := api.NewHandler(svc, wx)
	router := api.NewRouter(h, "frontend")

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// errCh receives a non-nil error if ListenAndServe fails for any reason
	// other than a clean shutdown.  Using a channel avoids calling log.Fatalf
	// inside a goroutine (which bypasses deferred cleanup).
	errCh := make(chan error, 1)
	go func() {
		log.Printf("level=info msg=\"server starting\" addr=:%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	// Block until an OS signal or a listen error arrives.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Printf("level=info msg=\"shutdown signal received\" signal=%s", sig)
	case err := <-errCh:
		log.Printf("level=fatal msg=\"server error\" err=%q", err.Error())
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	log.Printf("level=info msg=\"graceful shutdown initiated\" timeout=%s", cfg.ShutdownTimeout)
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("level=error msg=\"shutdown error\" err=%q", err.Error())
		os.Exit(1)
	}
	log.Printf("level=info msg=\"server stopped\"")
}
