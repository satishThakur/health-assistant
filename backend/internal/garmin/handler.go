package garmin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/satishthakur/health-assistant/backend/internal/db"
	"github.com/satishthakur/health-assistant/backend/internal/models"
)

// Handler handles Garmin data ingestion endpoints.
type Handler struct {
	eventRepo *db.EventRepository
}

// NewHandler creates a new garmin Handler.
func NewHandler(eventRepo *db.EventRepository) *Handler {
	return &Handler{eventRepo: eventRepo}
}

// HandleSleepIngestion handles POST /api/v1/garmin/ingest/sleep
func (h *Handler) HandleSleepIngestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload SleepPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode sleep payload: %v", err)
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if err := ValidateSleepPayload(&payload); err != nil {
		log.Printf("Sleep payload validation failed: %v", err)
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	event, err := transformSleepToEvent(&payload)
	if err != nil {
		log.Printf("Failed to transform sleep data: %v", err)
		http.Error(w, fmt.Sprintf("Transformation error: %v", err), http.StatusInternalServerError)
		return
	}

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
func (h *Handler) HandleActivityIngestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload ActivityPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode activity payload: %v", err)
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if err := ValidateActivityPayload(&payload); err != nil {
		log.Printf("Activity payload validation failed: %v", err)
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	event, err := transformActivityToEvent(&payload)
	if err != nil {
		log.Printf("Failed to transform activity data: %v", err)
		http.Error(w, fmt.Sprintf("Transformation error: %v", err), http.StatusInternalServerError)
		return
	}

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
func (h *Handler) HandleHRVIngestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload HRVPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode HRV payload: %v", err)
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if err := ValidateHRVPayload(&payload); err != nil {
		log.Printf("HRV payload validation failed: %v", err)
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	event, err := transformHRVToEvent(&payload)
	if err != nil {
		log.Printf("Failed to transform HRV data: %v", err)
		http.Error(w, fmt.Sprintf("Transformation error: %v", err), http.StatusInternalServerError)
		return
	}

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
func (h *Handler) HandleStressIngestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload StressPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Failed to decode stress payload: %v", err)
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if err := ValidateStressPayload(&payload); err != nil {
		log.Printf("Stress payload validation failed: %v", err)
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	event, err := transformStressToEvent(&payload)
	if err != nil {
		log.Printf("Failed to transform stress data: %v", err)
		http.Error(w, fmt.Sprintf("Transformation error: %v", err), http.StatusInternalServerError)
		return
	}

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

// HandleDailyStatsIngestion handles POST /api/v1/garmin/ingest/daily-stats
func (h *Handler) HandleDailyStatsIngestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload DailyStatsPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Invalid JSON in daily stats request: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := ValidateDailyStatsPayload(&payload); err != nil {
		log.Printf("Validation failed for daily stats: %v", err)
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	event, err := transformDailyStatsToEvent(&payload)
	if err != nil {
		log.Printf("Failed to transform daily stats: %v", err)
		http.Error(w, "Failed to process daily stats", http.StatusInternalServerError)
		return
	}

	result, err := h.eventRepo.InsertEvent(r.Context(), event)
	if err != nil {
		log.Printf("Failed to store daily stats: %v", err)
		http.Error(w, "Failed to store event", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully inserted daily stats for user %s on %s", payload.UserID, payload.Date)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "success",
		"was_inserted": result.WasInserted,
	})
}

// HandleBodyBatteryIngestion handles POST /api/v1/garmin/ingest/body-battery
func (h *Handler) HandleBodyBatteryIngestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload BodyBatteryPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Invalid JSON in body battery request: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := ValidateBodyBatteryPayload(&payload); err != nil {
		log.Printf("Validation failed for body battery: %v", err)
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	event, err := transformBodyBatteryToEvent(&payload)
	if err != nil {
		log.Printf("Failed to transform body battery: %v", err)
		http.Error(w, "Failed to process body battery", http.StatusInternalServerError)
		return
	}

	result, err := h.eventRepo.InsertEvent(r.Context(), event)
	if err != nil {
		log.Printf("Failed to store body battery: %v", err)
		http.Error(w, "Failed to store event", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully inserted body battery for user %s on %s", payload.UserID, payload.Date)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "success",
		"was_inserted": result.WasInserted,
	})
}

// Transform functions

func transformSleepToEvent(payload *SleepPayload) (*models.Event, error) {
	var eventTime time.Time
	if endTimestamp, ok := payload.SleepData["sleep_end_timestamp_gmt"].(string); ok {
		t, err := time.Parse(time.RFC3339, endTimestamp)
		if err == nil {
			eventTime = t
		}
	}
	if eventTime.IsZero() {
		t, _ := time.Parse("2006-01-02", payload.Date)
		eventTime = t.Add(8 * time.Hour)
	}

	garminSleep := models.GarminSleep{
		DurationMinutes:   int(getFloat64Value(payload.SleepData, "sleep_time_seconds") / 60),
		DeepSleepMinutes:  int(getFloat64Value(payload.SleepData, "deep_sleep_seconds") / 60),
		LightSleepMinutes: int(getFloat64Value(payload.SleepData, "light_sleep_seconds") / 60),
		REMSleepMinutes:   int(getFloat64Value(payload.SleepData, "rem_sleep_seconds") / 60),
		AwakeMinutes:      int(getFloat64Value(payload.SleepData, "awake_seconds") / 60),
		HRVAvg:            getFloat64Value(payload.SleepData, "average_hrv"),
	}

	if sleepScores, ok := payload.SleepData["sleep_scores"].(map[string]interface{}); ok {
		garminSleep.SleepScore = int(getFloat64Value(sleepScores, "overall_score"))
	}

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

func transformActivityToEvent(payload *ActivityPayload) (*models.Event, error) {
	var eventTime time.Time
	if startTime, ok := payload.ActivityData["start_time_gmt"].(string); ok {
		t, err := time.Parse(time.RFC3339, startTime)
		if err == nil {
			eventTime = t
		}
	}
	if eventTime.IsZero() {
		t, _ := time.Parse("2006-01-02", payload.Date)
		eventTime = t.Add(12 * time.Hour)
	}

	garminActivity := models.GarminActivity{
		ActivityType:    getStringValue(payload.ActivityData, "activity_type"),
		DurationMinutes: int(getFloat64Value(payload.ActivityData, "duration_seconds") / 60),
		Calories:        int(getFloat64Value(payload.ActivityData, "calories")),
		AvgHR:           int(getFloat64Value(payload.ActivityData, "average_heart_rate")),
		MaxHR:           int(getFloat64Value(payload.ActivityData, "max_heart_rate")),
		Distance:        getFloat64Value(payload.ActivityData, "distance_meters"),
	}

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

func transformHRVToEvent(payload *HRVPayload) (*models.Event, error) {
	eventTime, _ := time.Parse("2006-01-02", payload.Date)

	hrvData := map[string]interface{}{
		"average_hrv": getFloat64Value(payload.HRVData, "average_hrv"),
	}
	if maxHRV, ok := payload.HRVData["max_hrv"]; ok {
		hrvData["max_hrv"] = maxHRV
	}
	if minHRV, ok := payload.HRVData["min_hrv"]; ok {
		hrvData["min_hrv"] = minHRV
	}

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

func transformStressToEvent(payload *StressPayload) (*models.Event, error) {
	eventTime, _ := time.Parse("2006-01-02", payload.Date)

	stressData := map[string]interface{}{
		"average_stress_level": getFloat64Value(payload.StressData, "average_stress_level"),
	}
	if maxStress, ok := payload.StressData["max_stress_level"]; ok {
		stressData["max_stress_level"] = maxStress
	}
	if restStress, ok := payload.StressData["rest_stress_duration"]; ok {
		stressData["rest_stress_duration"] = restStress
	}

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

func transformDailyStatsToEvent(payload *DailyStatsPayload) (*models.Event, error) {
	eventTime, _ := time.Parse("2006-01-02", payload.Date)

	dailyStats := models.GarminDailyStats{
		Steps:                    int(getFloat64Value(payload.DailyStatsData, "steps")),
		Calories:                 int(getFloat64Value(payload.DailyStatsData, "calories")),
		DistanceMeters:           int(getFloat64Value(payload.DailyStatsData, "distance_meters")),
		ActiveCalories:           int(getFloat64Value(payload.DailyStatsData, "active_calories")),
		BMRCalories:              int(getFloat64Value(payload.DailyStatsData, "bmr_calories")),
		MinHeartRate:             int(getFloat64Value(payload.DailyStatsData, "min_heart_rate")),
		MaxHeartRate:             int(getFloat64Value(payload.DailyStatsData, "max_heart_rate")),
		RestingHeartRate:         int(getFloat64Value(payload.DailyStatsData, "resting_heart_rate")),
		ModerateIntensityMinutes: int(getFloat64Value(payload.DailyStatsData, "moderate_intensity_minutes")),
		VigorousIntensityMinutes: int(getFloat64Value(payload.DailyStatsData, "vigorous_intensity_minutes")),
	}

	dataJSON, err := json.Marshal(dailyStats)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal daily stats: %w", err)
	}

	return &models.Event{
		Time:      eventTime,
		UserID:    payload.UserID,
		EventType: models.EventTypeGarminDailyStats,
		Source:    models.SourceGarmin,
		Data:      dataJSON,
	}, nil
}

func transformBodyBatteryToEvent(payload *BodyBatteryPayload) (*models.Event, error) {
	eventTime, _ := time.Parse("2006-01-02", payload.Date)

	bodyBattery := models.GarminBodyBattery{
		Charged:      int(getFloat64Value(payload.BodyBatteryData, "charged")),
		Drained:      int(getFloat64Value(payload.BodyBatteryData, "drained")),
		HighestValue: int(getFloat64Value(payload.BodyBatteryData, "highest_value")),
		LowestValue:  int(getFloat64Value(payload.BodyBatteryData, "lowest_value")),
	}

	dataJSON, err := json.Marshal(bodyBattery)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body battery: %w", err)
	}

	return &models.Event{
		Time:      eventTime,
		UserID:    payload.UserID,
		EventType: models.EventTypeGarminBodyBattery,
		Source:    models.SourceGarmin,
		Data:      dataJSON,
	}, nil
}

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
