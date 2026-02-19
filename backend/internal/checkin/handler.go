package checkin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/satishthakur/health-assistant/backend/internal/db"
	"github.com/satishthakur/health-assistant/backend/internal/middleware"
	"github.com/satishthakur/health-assistant/backend/internal/models"
)

// Handler handles check-in related requests.
type Handler struct {
	eventRepo   *db.EventRepository
	checkinRepo *Repository
}

// NewHandler creates a new check-in Handler.
func NewHandler(eventRepo *db.EventRepository, checkinRepo *Repository) *Handler {
	return &Handler{
		eventRepo:   eventRepo,
		checkinRepo: checkinRepo,
	}
}

// HandleSubmission handles POST /api/v1/checkin
func (h *Handler) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to parse checkin payload: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := ValidatePayload(&payload); err != nil {
		log.Printf("Checkin validation failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Validation failed",
			"message": err.Error(),
		})
		return
	}

	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	feeling := models.SubjectiveFeeling{
		Energy:   payload.Energy,
		Mood:     payload.Mood,
		Focus:    payload.Focus,
		Physical: payload.Physical,
		Notes:    payload.Notes,
	}

	feelingJSON, err := json.Marshal(feeling)
	if err != nil {
		log.Printf("Failed to marshal feeling data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	event := &models.Event{
		Time:      startOfDay,
		UserID:    userID,
		EventType: models.EventTypeSubjectiveFeeling,
		Source:    models.SourceManual,
		Data:      feelingJSON,
	}

	result, err := h.eventRepo.InsertEvent(r.Context(), event)
	if err != nil {
		log.Printf("Failed to store checkin: %v", err)
		http.Error(w, "Failed to store checkin", http.StatusInternalServerError)
		return
	}

	action := "updated"
	if result.WasInserted {
		action = "inserted"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"action":    action,
		"timestamp": now,
		"data":      feeling,
	})

	log.Printf("Check-in %s for user %s: energy=%d, mood=%d, focus=%d, physical=%d",
		action, userID, payload.Energy, payload.Mood, payload.Focus, payload.Physical)
}

// HandleGetLatest handles GET /api/v1/checkin/latest
func (h *Handler) HandleGetLatest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	events, err := h.eventRepo.GetEventsByUserAndType(
		r.Context(), userID, models.EventTypeSubjectiveFeeling, startOfDay, endOfDay,
	)
	if err != nil {
		log.Printf("Failed to fetch latest checkin: %v", err)
		http.Error(w, "Failed to fetch checkin", http.StatusInternalServerError)
		return
	}

	if len(events) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"checkin": nil,
			"message": "No check-in for today",
		})
		return
	}

	var feeling models.SubjectiveFeeling
	if err := json.Unmarshal(events[0].Data, &feeling); err != nil {
		log.Printf("Failed to parse feeling data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"timestamp": events[0].Time,
		"checkin":   feeling,
	})
}

// HandleGetHistory handles GET /api/v1/checkin/history?days=30
func (h *Handler) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
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

	now := time.Now()
	startDate := now.AddDate(0, 0, -days)
	startTime := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	events, err := h.eventRepo.GetEventsByUserAndType(
		r.Context(), userID, models.EventTypeSubjectiveFeeling, startTime, now,
	)
	if err != nil {
		log.Printf("Failed to fetch checkin history: %v", err)
		http.Error(w, "Failed to fetch history", http.StatusInternalServerError)
		return
	}

	type HistoryItem struct {
		Date    string                   `json:"date"`
		Checkin models.SubjectiveFeeling `json:"checkin"`
	}

	history := make([]HistoryItem, 0, len(events))
	for _, event := range events {
		var feeling models.SubjectiveFeeling
		if err := json.Unmarshal(event.Data, &feeling); err != nil {
			log.Printf("Failed to parse feeling data: %v", err)
			continue
		}
		history = append(history, HistoryItem{
			Date:    event.Time.Format("2006-01-02"),
			Checkin: feeling,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"count":   len(history),
		"history": history,
	})
}
