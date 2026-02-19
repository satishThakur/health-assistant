package garmin

import (
	"errors"
	"time"
)

// SleepPayload represents the incoming sleep data from Python scheduler.
type SleepPayload struct {
	UserID    string                 `json:"user_id"`
	Date      string                 `json:"date"`
	SleepData map[string]interface{} `json:"sleep_data"`
}

// ActivityPayload represents the incoming activity data from Python scheduler.
type ActivityPayload struct {
	UserID       string                 `json:"user_id"`
	Date         string                 `json:"date"`
	ActivityData map[string]interface{} `json:"activity_data"`
}

// HRVPayload represents the incoming HRV data from Python scheduler.
type HRVPayload struct {
	UserID  string                 `json:"user_id"`
	Date    string                 `json:"date"`
	HRVData map[string]interface{} `json:"hrv_data"`
}

// StressPayload represents the incoming stress data from Python scheduler.
type StressPayload struct {
	UserID     string                 `json:"user_id"`
	Date       string                 `json:"date"`
	StressData map[string]interface{} `json:"stress_data"`
}

// DailyStatsPayload represents the incoming daily stats from Python scheduler.
type DailyStatsPayload struct {
	UserID         string                 `json:"user_id"`
	Date           string                 `json:"date"`
	DailyStatsData map[string]interface{} `json:"daily_stats_data"`
}

// BodyBatteryPayload represents the incoming body battery data from Python scheduler.
type BodyBatteryPayload struct {
	UserID          string                 `json:"user_id"`
	Date            string                 `json:"date"`
	BodyBatteryData map[string]interface{} `json:"body_battery_data"`
}

// ValidateSleepPayload validates the sleep data payload.
func ValidateSleepPayload(payload *SleepPayload) error {
	if payload.UserID == "" {
		return errors.New("user_id is required")
	}

	if payload.Date == "" {
		return errors.New("date is required")
	}

	if _, err := time.Parse("2006-01-02", payload.Date); err != nil {
		return errors.New("date must be in YYYY-MM-DD format")
	}

	if payload.SleepData == nil {
		return errors.New("sleep_data is required")
	}

	sleepTime, ok := getFloat64(payload.SleepData, "sleep_time_seconds")
	if !ok || sleepTime <= 0 {
		return errors.New("sleep_time_seconds must be a positive number")
	}

	if endTimestamp, exists := payload.SleepData["sleep_end_timestamp_gmt"]; exists {
		if str, ok := endTimestamp.(string); ok {
			if _, err := time.Parse(time.RFC3339, str); err != nil {
				return errors.New("sleep_end_timestamp_gmt must be in RFC3339 format")
			}
		}
	}

	return nil
}

// ValidateActivityPayload validates the activity data payload.
func ValidateActivityPayload(payload *ActivityPayload) error {
	if payload.UserID == "" {
		return errors.New("user_id is required")
	}

	if payload.Date == "" {
		return errors.New("date is required")
	}

	if _, err := time.Parse("2006-01-02", payload.Date); err != nil {
		return errors.New("date must be in YYYY-MM-DD format")
	}

	if payload.ActivityData == nil {
		return errors.New("activity_data is required")
	}

	activityType, ok := payload.ActivityData["activity_type"].(string)
	if !ok || activityType == "" {
		return errors.New("activity_type must be a non-empty string")
	}

	duration, ok := getFloat64(payload.ActivityData, "duration_seconds")
	if !ok || duration <= 0 {
		return errors.New("duration_seconds must be a positive number")
	}

	if startTime, exists := payload.ActivityData["start_time_gmt"]; exists {
		if str, ok := startTime.(string); ok {
			if _, err := time.Parse(time.RFC3339, str); err != nil {
				return errors.New("start_time_gmt must be in RFC3339 format")
			}
		}
	}

	return nil
}

// ValidateHRVPayload validates the HRV data payload.
func ValidateHRVPayload(payload *HRVPayload) error {
	if payload.UserID == "" {
		return errors.New("user_id is required")
	}

	if payload.Date == "" {
		return errors.New("date is required")
	}

	if _, err := time.Parse("2006-01-02", payload.Date); err != nil {
		return errors.New("date must be in YYYY-MM-DD format")
	}

	if payload.HRVData == nil {
		return errors.New("hrv_data is required")
	}

	hrvValue, ok := getFloat64(payload.HRVData, "average_hrv")
	if !ok || hrvValue < 0 {
		return errors.New("average_hrv must be a non-negative number")
	}

	return nil
}

// ValidateStressPayload validates the stress data payload.
func ValidateStressPayload(payload *StressPayload) error {
	if payload.UserID == "" {
		return errors.New("user_id is required")
	}

	if payload.Date == "" {
		return errors.New("date is required")
	}

	if _, err := time.Parse("2006-01-02", payload.Date); err != nil {
		return errors.New("date must be in YYYY-MM-DD format")
	}

	if payload.StressData == nil {
		return errors.New("stress_data is required")
	}

	if stressLevel, exists := payload.StressData["average_stress_level"]; exists {
		if level, ok := getFloat64OrInt(stressLevel); ok {
			if level < 0 || level > 100 {
				return errors.New("average_stress_level must be between 0 and 100")
			}
		}
	}

	return nil
}

// ValidateDailyStatsPayload validates the daily stats payload.
func ValidateDailyStatsPayload(payload *DailyStatsPayload) error {
	if payload.UserID == "" {
		return errors.New("user_id is required")
	}

	if payload.Date == "" {
		return errors.New("date is required")
	}

	if _, err := time.Parse("2006-01-02", payload.Date); err != nil {
		return errors.New("date must be in YYYY-MM-DD format")
	}

	if payload.DailyStatsData == nil {
		return errors.New("daily_stats_data is required")
	}

	steps, ok := getFloat64(payload.DailyStatsData, "steps")
	if !ok || steps < 0 {
		return errors.New("steps must be a non-negative number")
	}

	return nil
}

// ValidateBodyBatteryPayload validates the body battery payload.
func ValidateBodyBatteryPayload(payload *BodyBatteryPayload) error {
	if payload.UserID == "" {
		return errors.New("user_id is required")
	}

	if payload.Date == "" {
		return errors.New("date is required")
	}

	if _, err := time.Parse("2006-01-02", payload.Date); err != nil {
		return errors.New("date must be in YYYY-MM-DD format")
	}

	if payload.BodyBatteryData == nil {
		return errors.New("body_battery_data is required")
	}

	charged, okCharged := getFloat64(payload.BodyBatteryData, "charged")
	drained, okDrained := getFloat64(payload.BodyBatteryData, "drained")

	if (!okCharged && !okDrained) || (charged < 0 && drained < 0) {
		return errors.New("charged or drained must be valid non-negative numbers")
	}

	return nil
}

func getFloat64(data map[string]interface{}, key string) (float64, bool) {
	val, exists := data[key]
	if !exists {
		return 0, false
	}

	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

func getFloat64OrInt(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}
