package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/satishthakur/health-assistant/backend/internal/models"
)

// CheckinRepository handles database operations for check-ins
type CheckinRepository struct {
	db *Database
}

// NewCheckinRepository creates a new CheckinRepository
func NewCheckinRepository(db *Database) *CheckinRepository {
	return &CheckinRepository{db: db}
}

// DashboardData represents today's summary data
type DashboardData struct {
	Checkin *models.SubjectiveFeeling `json:"checkin,omitempty"`
	Garmin  *GarminSummary            `json:"garmin,omitempty"`
}

// GarminSummary represents aggregated Garmin data for today
type GarminSummary struct {
	Sleep       *models.GarminSleep       `json:"sleep,omitempty"`
	Activity    *models.GarminActivity    `json:"activity,omitempty"`
	HRV         *HRVData                  `json:"hrv,omitempty"`
	Stress      *StressData               `json:"stress,omitempty"`
	DailyStats  *models.GarminDailyStats  `json:"daily_stats,omitempty"`
	BodyBattery *models.GarminBodyBattery `json:"body_battery,omitempty"`
}

// HRVData represents HRV information
type HRVData struct {
	Average float64 `json:"average"`
}

// StressData represents stress information
type StressData struct {
	Average int    `json:"average"`
	Level   string `json:"level"` // low, moderate, high
}

// TrendData represents 7-day trend data
type TrendData struct {
	Date     string                    `json:"date"`
	Checkin  *models.SubjectiveFeeling `json:"checkin,omitempty"`
	Sleep    *models.GarminSleep       `json:"sleep,omitempty"`
	Activity *models.GarminActivity    `json:"activity,omitempty"`
}

// CorrelationInsight represents a correlation between metrics
type CorrelationInsight struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Confidence  float64                `json:"confidence"`
	SampleSize  int                    `json:"sample_size"`
	Details     map[string]interface{} `json:"details"`
}

// DailyData represents aggregated data for a single day
type DailyData struct {
	Feeling  *models.SubjectiveFeeling
	Sleep    *models.GarminSleep
	Activity *models.GarminActivity
}

// GetTodayDashboard retrieves today's check-in and Garmin data
func (r *CheckinRepository) GetTodayDashboard(ctx context.Context, userID string) (*DashboardData, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Query all today's events
	query := `
		SELECT time, event_type, data
		FROM events
		WHERE user_id = $1
			AND time >= $2
			AND time < $3
		ORDER BY time DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to query today's events: %w", err)
	}
	defer rows.Close()

	dashboard := &DashboardData{
		Garmin: &GarminSummary{},
	}

	for rows.Next() {
		var eventTime time.Time
		var eventType string
		var data []byte

		if err := rows.Scan(&eventTime, &eventType, &data); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		switch eventType {
		case models.EventTypeSubjectiveFeeling:
			var feeling models.SubjectiveFeeling
			if err := json.Unmarshal(data, &feeling); err == nil {
				dashboard.Checkin = &feeling
			}

		case models.EventTypeGarminSleep:
			var sleep models.GarminSleep
			if err := json.Unmarshal(data, &sleep); err == nil {
				dashboard.Garmin.Sleep = &sleep
			}

		case models.EventTypeGarminActivity:
			var activity models.GarminActivity
			if err := json.Unmarshal(data, &activity); err == nil {
				dashboard.Garmin.Activity = &activity
			}

		case models.EventTypeGarminHRV:
			// HRV data structure from Garmin (assuming it has avg field)
			var hrvData map[string]interface{}
			if err := json.Unmarshal(data, &hrvData); err == nil {
				if avg, ok := hrvData["average"].(float64); ok {
					dashboard.Garmin.HRV = &HRVData{Average: avg}
				}
			}

		case models.EventTypeGarminStress:
			// Stress data structure
			var stressData map[string]interface{}
			if err := json.Unmarshal(data, &stressData); err == nil {
				if avg, ok := stressData["average_stress_level"].(float64); ok {
					avgInt := int(avg)
					level := "low"
					if avgInt >= 26 && avgInt <= 50 {
						level = "moderate"
					} else if avgInt > 50 {
						level = "high"
					}
					dashboard.Garmin.Stress = &StressData{
						Average: avgInt,
						Level:   level,
					}
				}
			}

		case models.EventTypeGarminDailyStats:
			var dailyStats models.GarminDailyStats
			if err := json.Unmarshal(data, &dailyStats); err == nil {
				dashboard.Garmin.DailyStats = &dailyStats
			}

		case models.EventTypeGarminBodyBattery:
			var bodyBattery models.GarminBodyBattery
			if err := json.Unmarshal(data, &bodyBattery); err == nil {
				dashboard.Garmin.BodyBattery = &bodyBattery
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	return dashboard, nil
}

// GetWeekTrends retrieves 7-day trend data
func (r *CheckinRepository) GetWeekTrends(ctx context.Context, userID string) ([]TrendData, error) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -6) // Last 7 days including today
	startOfWeek := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	// Query all events for the past 7 days
	query := `
		SELECT time, event_type, data
		FROM events
		WHERE user_id = $1
			AND time >= $2
		ORDER BY time ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, startOfWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to query week trends: %w", err)
	}
	defer rows.Close()

	// Group by date
	trendsByDate := make(map[string]*TrendData)

	for rows.Next() {
		var eventTime time.Time
		var eventType string
		var data []byte

		if err := rows.Scan(&eventTime, &eventType, &data); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		dateKey := eventTime.Format("2006-01-02")
		if _, exists := trendsByDate[dateKey]; !exists {
			trendsByDate[dateKey] = &TrendData{Date: dateKey}
		}

		trend := trendsByDate[dateKey]

		switch eventType {
		case models.EventTypeSubjectiveFeeling:
			var feeling models.SubjectiveFeeling
			if err := json.Unmarshal(data, &feeling); err == nil {
				trend.Checkin = &feeling
			}

		case models.EventTypeGarminSleep:
			var sleep models.GarminSleep
			if err := json.Unmarshal(data, &sleep); err == nil {
				trend.Sleep = &sleep
			}

		case models.EventTypeGarminActivity:
			var activity models.GarminActivity
			if err := json.Unmarshal(data, &activity); err == nil {
				trend.Activity = &activity
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	// Convert map to sorted array
	trends := make([]TrendData, 0, len(trendsByDate))
	for _, trend := range trendsByDate {
		trends = append(trends, *trend)
	}

	return trends, nil
}

// GetCorrelations calculates simple correlations between Garmin data and feelings
func (r *CheckinRepository) GetCorrelations(ctx context.Context, userID string, days int) ([]CorrelationInsight, error) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -days)
	startTime := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	// Query all events for the specified period
	query := `
		SELECT time, event_type, data
		FROM events
		WHERE user_id = $1
			AND time >= $2
		ORDER BY time ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query correlation data: %w", err)
	}
	defer rows.Close()

	// Group data by date
	dailyData := make(map[string]*DailyData)

	for rows.Next() {
		var eventTime time.Time
		var eventType string
		var data []byte

		if err := rows.Scan(&eventTime, &eventType, &data); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		dateKey := eventTime.Format("2006-01-02")
		if _, exists := dailyData[dateKey]; !exists {
			dailyData[dateKey] = &DailyData{}
		}

		daily := dailyData[dateKey]

		switch eventType {
		case models.EventTypeSubjectiveFeeling:
			var feeling models.SubjectiveFeeling
			if err := json.Unmarshal(data, &feeling); err == nil {
				daily.Feeling = &feeling
			}

		case models.EventTypeGarminSleep:
			var sleep models.GarminSleep
			if err := json.Unmarshal(data, &sleep); err == nil {
				daily.Sleep = &sleep
			}

		case models.EventTypeGarminActivity:
			var activity models.GarminActivity
			if err := json.Unmarshal(data, &activity); err == nil {
				daily.Activity = &activity
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	// Calculate correlations
	insights := r.calculateCorrelations(dailyData)

	return insights, nil
}

// calculateCorrelations performs simple correlation calculations
func (r *CheckinRepository) calculateCorrelations(dailyData map[string]*DailyData) []CorrelationInsight {
	insights := []CorrelationInsight{}

	// Sleep vs Energy correlation
	if insight := r.calculateSleepEnergyCorrelation(dailyData); insight != nil {
		insights = append(insights, *insight)
	}

	// Activity vs Mood correlation
	if insight := r.calculateActivityMoodCorrelation(dailyData); insight != nil {
		insights = append(insights, *insight)
	}

	// Sleep vs Focus correlation
	if insight := r.calculateSleepFocusCorrelation(dailyData); insight != nil {
		insights = append(insights, *insight)
	}

	return insights
}

// calculateSleepEnergyCorrelation calculates correlation between sleep duration and energy
func (r *CheckinRepository) calculateSleepEnergyCorrelation(dailyData map[string]*DailyData) *CorrelationInsight {
	var energyWithGoodSleep []int
	var energyWithPoorSleep []int

	for _, data := range dailyData {
		if data.Feeling == nil || data.Sleep == nil {
			continue
		}

		sleepHours := float64(data.Sleep.DurationMinutes) / 60.0
		if sleepHours >= 7.0 {
			energyWithGoodSleep = append(energyWithGoodSleep, data.Feeling.Energy)
		} else {
			energyWithPoorSleep = append(energyWithPoorSleep, data.Feeling.Energy)
		}
	}

	// Need at least 5 samples in each group
	if len(energyWithGoodSleep) < 5 || len(energyWithPoorSleep) < 5 {
		return nil
	}

	avgWithGoodSleep := average(energyWithGoodSleep)
	avgWithPoorSleep := average(energyWithPoorSleep)
	improvement := ((avgWithGoodSleep - avgWithPoorSleep) / avgWithPoorSleep) * 100

	if improvement < 5 { // Only show if improvement is at least 5%
		return nil
	}

	return &CorrelationInsight{
		Type:        "sleep_energy",
		Description: fmt.Sprintf("Your energy is %.0f%% higher when you sleep 7+ hours", improvement),
		Confidence:  0.85,
		SampleSize:  len(energyWithGoodSleep) + len(energyWithPoorSleep),
		Details: map[string]interface{}{
			"condition":           "sleep >= 7 hours",
			"avg_energy_with":     avgWithGoodSleep,
			"avg_energy_without":  avgWithPoorSleep,
			"improvement_percent": improvement,
		},
	}
}

// calculateActivityMoodCorrelation calculates correlation between activity and mood
func (r *CheckinRepository) calculateActivityMoodCorrelation(dailyData map[string]*DailyData) *CorrelationInsight {
	var moodWithActivity []int
	var moodWithoutActivity []int

	for _, data := range dailyData {
		if data.Feeling == nil || data.Activity == nil {
			continue
		}

		if data.Activity.DurationMinutes >= 30 {
			moodWithActivity = append(moodWithActivity, data.Feeling.Mood)
		} else {
			moodWithoutActivity = append(moodWithoutActivity, data.Feeling.Mood)
		}
	}

	if len(moodWithActivity) < 5 || len(moodWithoutActivity) < 5 {
		return nil
	}

	avgWithActivity := average(moodWithActivity)
	avgWithoutActivity := average(moodWithoutActivity)
	improvement := ((avgWithActivity - avgWithoutActivity) / avgWithoutActivity) * 100

	if improvement < 5 {
		return nil
	}

	return &CorrelationInsight{
		Type:        "activity_mood",
		Description: fmt.Sprintf("Your mood improves by %.0f%% on active days (30+ min)", improvement),
		Confidence:  0.78,
		SampleSize:  len(moodWithActivity) + len(moodWithoutActivity),
		Details: map[string]interface{}{
			"condition":           "activity >= 30 minutes",
			"avg_mood_with":       avgWithActivity,
			"avg_mood_without":    avgWithoutActivity,
			"improvement_percent": improvement,
		},
	}
}

// calculateSleepFocusCorrelation calculates correlation between sleep quality and focus
func (r *CheckinRepository) calculateSleepFocusCorrelation(dailyData map[string]*DailyData) *CorrelationInsight {
	var focusWithGoodSleep []int
	var focusWithPoorSleep []int

	for _, data := range dailyData {
		if data.Feeling == nil || data.Sleep == nil {
			continue
		}

		if data.Sleep.SleepScore >= 80 {
			focusWithGoodSleep = append(focusWithGoodSleep, data.Feeling.Focus)
		} else {
			focusWithPoorSleep = append(focusWithPoorSleep, data.Feeling.Focus)
		}
	}

	if len(focusWithGoodSleep) < 5 || len(focusWithPoorSleep) < 5 {
		return nil
	}

	avgWithGoodSleep := average(focusWithGoodSleep)
	avgWithPoorSleep := average(focusWithPoorSleep)
	improvement := ((avgWithGoodSleep - avgWithPoorSleep) / avgWithPoorSleep) * 100

	if improvement < 5 {
		return nil
	}

	return &CorrelationInsight{
		Type:        "sleep_focus",
		Description: fmt.Sprintf("Your focus is %.0f%% better after quality sleep (score 80+)", improvement),
		Confidence:  0.82,
		SampleSize:  len(focusWithGoodSleep) + len(focusWithPoorSleep),
		Details: map[string]interface{}{
			"condition":           "sleep_score >= 80",
			"avg_focus_with":      avgWithGoodSleep,
			"avg_focus_without":   avgWithPoorSleep,
			"improvement_percent": improvement,
		},
	}
}

// average calculates the average of a slice of integers
func average(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0
	for _, v := range values {
		sum += v
	}
	return float64(sum) / float64(len(values))
}
