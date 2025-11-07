package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bgc-app/internal/app"
	"bgc-app/internal/config"
	"bgc-app/internal/observability/tracing"
	"bgc-app/internal/repository/postgres"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize OpenTelemetry tracer
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	tracerShutdown, err := tracing.InitTracer("bgc-api", environment)
	if err != nil {
		log.Printf("Warning: Failed to initialize tracer: %v (tracing disabled)", err)
	} else {
		defer func() {
			if err := tracerShutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer: %v", err)
			}
		}()
	}

	db := postgres.MustConnect(cfg)
	defer db.Close()

	server := app.NewServer(cfg, db)

	// Graceful shutdown
	go func() {
		if err := server.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down BGC API...")
}
