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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// DB_DSN is read here for configuration; it is intentionally not used in
	// this synthetic application — its presence seeds defect SD-08: the env
	// var is consumed in code but is absent from docs/requirements.md.
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		dbDSN = "postgres://localhost:5432/orders?sslmode=disable"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	router := handlers.NewRouter()

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

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
