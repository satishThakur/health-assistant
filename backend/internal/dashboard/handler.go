package dashboard

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/satishthakur/health-assistant/backend/internal/checkin"
	"github.com/satishthakur/health-assistant/backend/internal/middleware"
)

// Handler handles dashboard and trends requests.
type Handler struct {
	checkinRepo *checkin.Repository
}

// NewHandler creates a new dashboard Handler.
func NewHandler(checkinRepo *checkin.Repository) *Handler {
	return &Handler{checkinRepo: checkinRepo}
}

// HandleGetToday handles GET /api/v1/dashboard/today
func (h *Handler) HandleGetToday(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

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
func (h *Handler) HandleGetWeekTrends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

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
func (h *Handler) HandleGetCorrelations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	daysParam := r.URL.Query().Get("days")
	days := 30
	if daysParam != "" {
		if _, err := fmt.Sscanf(daysParam, "%d", &days); err != nil {
			days = 30
		}
	}

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
