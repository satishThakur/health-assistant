package audit

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Handler handles sync audit endpoints.
type Handler struct {
	repo *Repository
}

// NewHandler creates a new audit Handler.
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// HandlePostSyncAudit handles POST /api/v1/audit/sync
func (h *Handler) HandlePostSyncAudit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload SyncAudit
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode audit payload: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.repo.InsertSyncAudit(r.Context(), &payload); err != nil {
		log.Printf("Failed to insert sync audit: %v", err)
		http.Error(w, "Failed to store audit", http.StatusInternalServerError)
		return
	}

	log.Printf("Audit recorded: user=%s, type=%s, date=%s, status=%s",
		payload.UserID, payload.DataType, payload.TargetDate, payload.Status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"id":     payload.ID,
	})
}

// HandleGetRecentSyncAudits handles GET /api/v1/audit/sync/recent?user_id=X&limit=Y
func (h *Handler) HandleGetRecentSyncAudits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
			limit = l
		}
	}

	audits, err := h.repo.GetRecentSyncAudits(r.Context(), userID, limit)
	if err != nil {
		log.Printf("Failed to get recent sync audits: %v", err)
		http.Error(w, "Failed to retrieve audits", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(audits)
}

// HandleGetSyncAuditsByType handles GET /api/v1/audit/sync/by-type?data_type=X&limit=Y
func (h *Handler) HandleGetSyncAuditsByType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dataType := r.URL.Query().Get("data_type")
	if dataType == "" {
		http.Error(w, "data_type parameter is required", http.StatusBadRequest)
		return
	}

	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
			limit = l
		}
	}

	audits, err := h.repo.GetSyncAuditsByDataType(r.Context(), dataType, limit)
	if err != nil {
		log.Printf("Failed to get sync audits by type: %v", err)
		http.Error(w, "Failed to retrieve audits", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(audits)
}

// HandleGetSyncAuditStats handles GET /api/v1/audit/sync/stats?user_id=X&start=Y&end=Z
func (h *Handler) HandleGetSyncAuditStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	if startStr := r.URL.Query().Get("start"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = t
		}
	}

	if endStr := r.URL.Query().Get("end"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = t
		}
	}

	stats, err := h.repo.GetSyncAuditStats(r.Context(), userID, startDate, endDate)
	if err != nil {
		log.Printf("Failed to get sync audit stats: %v", err)
		http.Error(w, "Failed to retrieve stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}
