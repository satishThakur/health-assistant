package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/satishthakur/health-assistant/backend/internal/db"
)

// DashboardHandler handles dashboard and trends requests
type DashboardHandler struct {
	checkinRepo *db.CheckinRepository
}

// NewDashboardHandler creates a new DashboardHandler
func NewDashboardHandler(checkinRepo *db.CheckinRepository) *DashboardHandler {
	return &DashboardHandler{
		checkinRepo: checkinRepo,
	}
}

// HandleGetTodayDashboard handles GET /api/v1/dashboard/today
func (h *DashboardHandler) HandleGetTodayDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Extract user_id from JWT token
	userID := "00000000-0000-0000-0000-000000000001"

	// Get today's dashboard data
	dashboard, err := h.checkinRepo.GetTodayDashboard(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to fetch today's dashboard: %v", err)
		http.Error(w, "Failed to fetch dashboard", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   dashboard,
	})
}

// HandleGetWeekTrends handles GET /api/v1/trends/week
func (h *DashboardHandler) HandleGetWeekTrends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Extract user_id from JWT token
	userID := "00000000-0000-0000-0000-000000000001"

	// Get 7-day trends
	trends, err := h.checkinRepo.GetWeekTrends(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to fetch week trends: %v", err)
		http.Error(w, "Failed to fetch trends", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"count":  len(trends),
		"trends": trends,
	})
}

// HandleGetCorrelations handles GET /api/v1/insights/correlations
func (h *DashboardHandler) HandleGetCorrelations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Extract user_id from JWT token
	userID := "00000000-0000-0000-0000-000000000001"

	// Parse days parameter (default 30)
	daysParam := r.URL.Query().Get("days")
	days := 30
	if daysParam != "" {
		if _, err := fmt.Sscanf(daysParam, "%d", &days); err != nil {
			days = 30
		}
	}

	// Get correlations
	insights, err := h.checkinRepo.GetCorrelations(r.Context(), userID, days)
	if err != nil {
		log.Printf("Failed to calculate correlations: %v", err)
		http.Error(w, "Failed to calculate insights", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "success",
		"count":        len(insights),
		"correlations": insights,
	})
}
