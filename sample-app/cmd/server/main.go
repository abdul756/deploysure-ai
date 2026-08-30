// Package main is the entry point for the orders-api server.
// It reads configuration from environment variables, wires the HTTP router,
// and handles graceful shutdown on SIGTERM / SIGINT.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abdul756/deploysure-ai/sample-app/internal/handlers"
)

// buildServer constructs an *http.Server with the application router and
// standard timeouts. It is exported for testing via main_test.go.
func buildServer(port string) *http.Server {
	return &http.Server{
		Addr:         ":" + port,
		Handler:      handlers.NewRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// DB_DSN is read here for configuration; it is intentionally not used in
	// this synthetic application.
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		dbDSN = "postgres://localhost:5432/orders?sslmode=disable"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	srv := buildServer(port)

	// Start server in a goroutine so the main goroutine can block on signals.
	go func() {
		log.Printf("level=%s msg=\"server starting\" addr=:%s", logLevel, port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Block until SIGTERM or SIGINT.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("msg=\"shutdown signal received, draining requests\"")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("msg=\"server exited cleanly\"")
}
