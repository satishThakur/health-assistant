package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/satishthakur/health-assistant/backend/internal/db"
	"github.com/satishthakur/health-assistant/backend/internal/models"
	"github.com/satishthakur/health-assistant/backend/internal/validation"
)

// GarminIngestionHandler handles Garmin data ingestion endpoints
type GarminIngestionHandler struct {
	eventRepo *db.EventRepository
}

// NewGarminIngestionHandler creates a new GarminIngestionHandler
func NewGarminIngestionHandler(eventRepo *db.EventRepository) *GarminIngestionHandler {
	return &GarminIngestionHandler{
		eventRepo: eventRepo,
	}
}

// HandleSleepIngestion handles POST /api/v1/garmin/ingest/sleep
func (h *GarminIngestionHandler) HandleSleepIngestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var payload validation.GarminSleepPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode sleep payload: %v", err)
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Validate payload
	if err := validation.ValidateSleepPayload(&payload); err != nil {
		log.Printf("Sleep payload validation failed: %v", err)
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	// Transform to internal model
	event, err := transformSleepToEvent(&payload)
	if err != nil {
		log.Printf("Failed to transform sleep data: %v", err)
		http.Error(w, fmt.Sprintf("Transformation error: %v", err), http.StatusInternalServerError)
		return
	}

	// Store in database
	result, err := h.eventRepo.InsertEvent(r.Context(), event)
	if err != nil {
		log.Printf("Failed to insert sleep event: %v", err)
		http.Error(w, "Failed to store event", http.StatusInternalServerError)
		return
	}

	action := "updated"
	if result.WasInserted {
		action = "inserted"
	}
	log.Printf("Successfully %s sleep data for user %s on %s", action, payload.UserID, payload.Date)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "success",
		"action":       action,
		"was_inserted": result.WasInserted,
	})
}

// HandleActivityIngestion handles POST /api/v1/garmin/ingest/activity
func (h *GarminIngestionHandler) HandleActivityIngestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var payload validation.GarminActivityPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode activity payload: %v", err)
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Validate payload
	if err := validation.ValidateActivityPayload(&payload); err != nil {
		log.Printf("Activity payload validation failed: %v", err)
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	// Transform to internal model
	event, err := transformActivityToEvent(&payload)
	if err != nil {
		log.Printf("Failed to transform activity data: %v", err)
		http.Error(w, fmt.Sprintf("Transformation error: %v", err), http.StatusInternalServerError)
		return
	}

	// Store in database
	result, err := h.eventRepo.InsertEvent(r.Context(), event)
	if err != nil {
		log.Printf("Failed to insert activity event: %v", err)
		http.Error(w, "Failed to store event", http.StatusInternalServerError)
		return
	}

	action := "updated"
	if result.WasInserted {
		action = "inserted"
	}
	log.Printf("Successfully %s activity data for user %s on %s", action, payload.UserID, payload.Date)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "success",
		"action":       action,
		"was_inserted": result.WasInserted,
	})
}

// HandleHRVIngestion handles POST /api/v1/garmin/ingest/hrv
func (h *GarminIngestionHandler) HandleHRVIngestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var payload validation.GarminHRVPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode HRV payload: %v", err)
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Validate payload
	if err := validation.ValidateHRVPayload(&payload); err != nil {
		log.Printf("HRV payload validation failed: %v", err)
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	// Transform to internal model
	event, err := transformHRVToEvent(&payload)
	if err != nil {
		log.Printf("Failed to transform HRV data: %v", err)
		http.Error(w, fmt.Sprintf("Transformation error: %v", err), http.StatusInternalServerError)
		return
	}

	// Store in database
	result, err := h.eventRepo.InsertEvent(r.Context(), event)
	if err != nil {
		log.Printf("Failed to insert HRV event: %v", err)
		http.Error(w, "Failed to store event", http.StatusInternalServerError)
		return
	}

	action := "updated"
	if result.WasInserted {
		action = "inserted"
	}
	log.Printf("Successfully %s HRV data for user %s on %s", action, payload.UserID, payload.Date)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "success",
		"action":       action,
		"was_inserted": result.WasInserted,
	})
}

// HandleStressIngestion handles POST /api/v1/garmin/ingest/stress
func (h *GarminIngestionHandler) HandleStressIngestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var payload validation.GarminStressPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode stress payload: %v", err)
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Validate payload
	if err := validation.ValidateStressPayload(&payload); err != nil {
		log.Printf("Stress payload validation failed: %v", err)
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	// Transform to internal model
	event, err := transformStressToEvent(&payload)
	if err != nil {
		log.Printf("Failed to transform stress data: %v", err)
		http.Error(w, fmt.Sprintf("Transformation error: %v", err), http.StatusInternalServerError)
		return
	}

	// Store in database
	result, err := h.eventRepo.InsertEvent(r.Context(), event)
	if err != nil {
		log.Printf("Failed to insert stress event: %v", err)
		http.Error(w, "Failed to store event", http.StatusInternalServerError)
		return
	}

	action := "updated"
	if result.WasInserted {
		action = "inserted"
	}
	log.Printf("Successfully %s stress data for user %s on %s", action, payload.UserID, payload.Date)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "success",
		"action":       action,
		"was_inserted": result.WasInserted,
	})
}

// Transform functions

func transformSleepToEvent(payload *validation.GarminSleepPayload) (*models.Event, error) {
	// Determine event time (use sleep_end_timestamp_gmt if available, otherwise use date at midnight)
	var eventTime time.Time
	if endTimestamp, ok := payload.SleepData["sleep_end_timestamp_gmt"].(string); ok {
		t, err := time.Parse(time.RFC3339, endTimestamp)
		if err == nil {
			eventTime = t
		}
	}
	if eventTime.IsZero() {
		t, _ := time.Parse("2006-01-02", payload.Date)
		eventTime = t.Add(8 * time.Hour) // Assume sleep ends at 8 AM if no timestamp
	}

	// Build GarminSleep model
	garminSleep := models.GarminSleep{
		DurationMinutes:   int(getFloat64Value(payload.SleepData, "sleep_time_seconds") / 60),
		DeepSleepMinutes:  int(getFloat64Value(payload.SleepData, "deep_sleep_seconds") / 60),
		LightSleepMinutes: int(getFloat64Value(payload.SleepData, "light_sleep_seconds") / 60),
		REMSleepMinutes:   int(getFloat64Value(payload.SleepData, "rem_sleep_seconds") / 60),
		AwakeMinutes:      int(getFloat64Value(payload.SleepData, "awake_seconds") / 60),
		HRVAvg:            getFloat64Value(payload.SleepData, "average_hrv"),
	}

	// Extract sleep score if available
	if sleepScores, ok := payload.SleepData["sleep_scores"].(map[string]interface{}); ok {
		garminSleep.SleepScore = int(getFloat64Value(sleepScores, "overall_score"))
	}

	// Marshal to JSON
	dataJSON, err := json.Marshal(garminSleep)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sleep data: %w", err)
	}

	return &models.Event{
		Time:      eventTime,
		UserID:    payload.UserID,
		EventType: models.EventTypeGarminSleep,
		Source:    models.SourceGarmin,
		Data:      dataJSON,
	}, nil
}

func transformActivityToEvent(payload *validation.GarminActivityPayload) (*models.Event, error) {
	// Determine event time (use start_time_gmt if available, otherwise use date at noon)
	var eventTime time.Time
	if startTime, ok := payload.ActivityData["start_time_gmt"].(string); ok {
		t, err := time.Parse(time.RFC3339, startTime)
		if err == nil {
			eventTime = t
		}
	}
	if eventTime.IsZero() {
		t, _ := time.Parse("2006-01-02", payload.Date)
		eventTime = t.Add(12 * time.Hour) // Assume activity at noon if no timestamp
	}

	// Build GarminActivity model
	garminActivity := models.GarminActivity{
		ActivityType:    getStringValue(payload.ActivityData, "activity_type"),
		DurationMinutes: int(getFloat64Value(payload.ActivityData, "duration_seconds") / 60),
		Calories:        int(getFloat64Value(payload.ActivityData, "calories")),
		AvgHR:           int(getFloat64Value(payload.ActivityData, "average_heart_rate")),
		MaxHR:           int(getFloat64Value(payload.ActivityData, "max_heart_rate")),
		Distance:        getFloat64Value(payload.ActivityData, "distance_meters"),
	}

	// Marshal to JSON
	dataJSON, err := json.Marshal(garminActivity)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal activity data: %w", err)
	}

	return &models.Event{
		Time:      eventTime,
		UserID:    payload.UserID,
		EventType: models.EventTypeGarminActivity,
		Source:    models.SourceGarmin,
		Data:      dataJSON,
	}, nil
}

func transformHRVToEvent(payload *validation.GarminHRVPayload) (*models.Event, error) {
	// Use date at midnight for HRV data
	eventTime, _ := time.Parse("2006-01-02", payload.Date)

	// Build HRV data structure
	hrvData := map[string]interface{}{
		"average_hrv": getFloat64Value(payload.HRVData, "average_hrv"),
	}

	// Include additional HRV fields if present
	if maxHRV, ok := payload.HRVData["max_hrv"]; ok {
		hrvData["max_hrv"] = maxHRV
	}
	if minHRV, ok := payload.HRVData["min_hrv"]; ok {
		hrvData["min_hrv"] = minHRV
	}

	// Marshal to JSON
	dataJSON, err := json.Marshal(hrvData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal HRV data: %w", err)
	}

	return &models.Event{
		Time:      eventTime,
		UserID:    payload.UserID,
		EventType: models.EventTypeGarminHRV,
		Source:    models.SourceGarmin,
		Data:      dataJSON,
	}, nil
}

func transformStressToEvent(payload *validation.GarminStressPayload) (*models.Event, error) {
	// Use date at midnight for stress data
	eventTime, _ := time.Parse("2006-01-02", payload.Date)

	// Build stress data structure
	stressData := map[string]interface{}{
		"average_stress_level": getFloat64Value(payload.StressData, "average_stress_level"),
	}

	// Include additional stress fields if present
	if maxStress, ok := payload.StressData["max_stress_level"]; ok {
		stressData["max_stress_level"] = maxStress
	}
	if restStress, ok := payload.StressData["rest_stress_duration"]; ok {
		stressData["rest_stress_duration"] = restStress
	}

	// Marshal to JSON
	dataJSON, err := json.Marshal(stressData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stress data: %w", err)
	}

	return &models.Event{
		Time:      eventTime,
		UserID:    payload.UserID,
		EventType: models.EventTypeGarminStress,
		Source:    models.SourceGarmin,
		Data:      dataJSON,
	}, nil
}

// Helper functions

func getFloat64Value(data map[string]interface{}, key string) float64 {
	val, exists := data[key]
	if !exists {
		return 0
	}

	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

func getStringValue(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}
