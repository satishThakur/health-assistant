package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/satishthakur/health-assistant/backend/internal/config"
	"github.com/satishthakur/health-assistant/backend/internal/db"
	"github.com/satishthakur/health-assistant/backend/internal/handlers"
)

func main() {
	log.Println("Starting Ingestion Service...")

	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	ctx := context.Background()
	database, err := db.NewDatabase(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("Database connection established")

	// Create repositories
	eventRepo := db.NewEventRepository(database)

	// Create handlers
	garminHandler := handlers.NewGarminIngestionHandler(eventRepo)

	// Setup routes
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database health
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		dbStatus := "connected"
		if err := database.Health(ctx); err != nil {
			dbStatus = "disconnected"
			log.Printf("Database health check failed: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":          "healthy",
			"service":         "ingestion-service",
			"database_status": dbStatus,
		})
	})

	// Garmin ingestion endpoints
	mux.HandleFunc("/api/v1/garmin/ingest/sleep", garminHandler.HandleSleepIngestion)
	mux.HandleFunc("/api/v1/garmin/ingest/activity", garminHandler.HandleActivityIngestion)
	mux.HandleFunc("/api/v1/garmin/ingest/hrv", garminHandler.HandleHRVIngestion)
	mux.HandleFunc("/api/v1/garmin/ingest/stress", garminHandler.HandleStressIngestion)

	// Create HTTP server
	port := ":8083"
	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Ingestion Service listening on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
