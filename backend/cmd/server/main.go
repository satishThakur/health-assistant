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

	"github.com/satishthakur/health-assistant/backend/internal/audit"
	"github.com/satishthakur/health-assistant/backend/internal/auth"
	"github.com/satishthakur/health-assistant/backend/internal/checkin"
	"github.com/satishthakur/health-assistant/backend/internal/config"
	"github.com/satishthakur/health-assistant/backend/internal/dashboard"
	"github.com/satishthakur/health-assistant/backend/internal/db"
	"github.com/satishthakur/health-assistant/backend/internal/garmin"
	"github.com/satishthakur/health-assistant/backend/internal/middleware"
)

func main() {
	log.Println("Starting Health Assistant Server...")

	// Load configuration
	cfg := config.Load()

	// Initialize JWT token service
	tokenService, err := auth.NewTokenService(cfg.Auth.JWTSecret, cfg.Auth.TokenDuration)
	if err != nil {
		log.Fatalf("Invalid JWT configuration: %v", err)
	}

	// Initialize Google verifier (empty clientID = dev mode, audience check skipped)
	googleVerifier := auth.NewGoogleVerifier(os.Getenv("GOOGLE_CLIENT_ID"))

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
	userRepo := auth.NewUserRepository(database)
	checkinRepo := checkin.NewRepository(database)
	auditRepo := audit.NewRepository(database)

	// Create handlers
	authHandler := auth.NewHandler(googleVerifier, userRepo, tokenService)
	garminHandler := garmin.NewHandler(eventRepo)
	auditHandler := audit.NewHandler(auditRepo)
	checkinHandler := checkin.NewHandler(eventRepo, checkinRepo)
	dashboardHandler := dashboard.NewHandler(checkinRepo)

	// Build middleware
	requireAuth := middleware.WithAuth(tokenService)
	requireIngest := middleware.WithIngestSecret(os.Getenv("GARMIN_INGEST_SECRET"))

	// Setup routes
	mux := http.NewServeMux()

	// Health check endpoint (public)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
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
			"service":         "health-assistant",
			"database_status": dbStatus,
		})
	})

	// Auth endpoint (public — no auth middleware)
	mux.HandleFunc("/api/v1/auth/google", authHandler.HandleGoogleAuth)

	// Garmin ingestion endpoints (ingest-secret protected, server-to-server)
	mux.Handle("/api/v1/garmin/ingest/sleep", requireIngest(http.HandlerFunc(garminHandler.HandleSleepIngestion)))
	mux.Handle("/api/v1/garmin/ingest/activity", requireIngest(http.HandlerFunc(garminHandler.HandleActivityIngestion)))
	mux.Handle("/api/v1/garmin/ingest/hrv", requireIngest(http.HandlerFunc(garminHandler.HandleHRVIngestion)))
	mux.Handle("/api/v1/garmin/ingest/stress", requireIngest(http.HandlerFunc(garminHandler.HandleStressIngestion)))
	mux.Handle("/api/v1/garmin/ingest/daily-stats", requireIngest(http.HandlerFunc(garminHandler.HandleDailyStatsIngestion)))
	mux.Handle("/api/v1/garmin/ingest/body-battery", requireIngest(http.HandlerFunc(garminHandler.HandleBodyBatteryIngestion)))

	// Audit endpoints (JWT protected)
	mux.Handle("/api/v1/audit/sync", requireAuth(http.HandlerFunc(auditHandler.HandlePostSyncAudit)))
	mux.Handle("/api/v1/audit/sync/recent", requireAuth(http.HandlerFunc(auditHandler.HandleGetRecentSyncAudits)))
	mux.Handle("/api/v1/audit/sync/by-type", requireAuth(http.HandlerFunc(auditHandler.HandleGetSyncAuditsByType)))
	mux.Handle("/api/v1/audit/sync/stats", requireAuth(http.HandlerFunc(auditHandler.HandleGetSyncAuditStats)))

	// Check-in endpoints (JWT protected)
	mux.Handle("/api/v1/checkin", requireAuth(http.HandlerFunc(checkinHandler.HandleSubmission)))
	mux.Handle("/api/v1/checkin/latest", requireAuth(http.HandlerFunc(checkinHandler.HandleGetLatest)))
	mux.Handle("/api/v1/checkin/history", requireAuth(http.HandlerFunc(checkinHandler.HandleGetHistory)))

	// Dashboard and trends endpoints (JWT protected)
	mux.Handle("/api/v1/dashboard/today", requireAuth(http.HandlerFunc(dashboardHandler.HandleGetToday)))
	mux.Handle("/api/v1/trends/week", requireAuth(http.HandlerFunc(dashboardHandler.HandleGetWeekTrends)))
	mux.Handle("/api/v1/insights/correlations", requireAuth(http.HandlerFunc(dashboardHandler.HandleGetCorrelations)))

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
		log.Printf("Health Assistant Server listening on port %s", port)
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
