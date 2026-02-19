package checkin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/satishthakur/health-assistant/backend/internal/db"
	"github.com/satishthakur/health-assistant/backend/internal/models"
)

// Repository handles database operations for check-in and dashboard queries.
type Repository struct {
	db *db.Database
}

// NewRepository creates a new Repository.
func NewRepository(database *db.Database) *Repository {
	return &Repository{db: database}
}

// DashboardData represents today's summary data.
type DashboardData struct {
	Checkin *models.SubjectiveFeeling `json:"checkin,omitempty"`
	Garmin  *GarminSummary            `json:"garmin,omitempty"`
}

// GarminSummary represents aggregated Garmin data for today.
type GarminSummary struct {
	Sleep       *models.GarminSleep       `json:"sleep,omitempty"`
	Activity    *models.GarminActivity    `json:"activity,omitempty"`
	HRV         *HRVData                  `json:"hrv,omitempty"`
	Stress      *StressData               `json:"stress,omitempty"`
	DailyStats  *models.GarminDailyStats  `json:"daily_stats,omitempty"`
	BodyBattery *models.GarminBodyBattery `json:"body_battery,omitempty"`
}

// HRVData represents HRV information.
type HRVData struct {
	Average float64 `json:"average"`
}

// StressData represents stress information.
type StressData struct {
	Average int    `json:"average"`
	Level   string `json:"level"` // low, moderate, high
}

// TrendData represents 7-day trend data.
type TrendData struct {
	Date     string                    `json:"date"`
	Checkin  *models.SubjectiveFeeling `json:"checkin,omitempty"`
	Sleep    *models.GarminSleep       `json:"sleep,omitempty"`
	Activity *models.GarminActivity    `json:"activity,omitempty"`
}

// CorrelationInsight represents a correlation between metrics.
type CorrelationInsight struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Confidence  float64                `json:"confidence"`
	SampleSize  int                    `json:"sample_size"`
	Details     map[string]interface{} `json:"details"`
}

// dailyData represents aggregated data for a single day (used internally).
type dailyData struct {
	Feeling  *models.SubjectiveFeeling
	Sleep    *models.GarminSleep
	Activity *models.GarminActivity
}

// GetTodayDashboard retrieves today's check-in and Garmin data.
func (r *Repository) GetTodayDashboard(ctx context.Context, userID string) (*DashboardData, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

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

	dashboard := &DashboardData{Garmin: &GarminSummary{}}

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
			var hrvData map[string]interface{}
			if err := json.Unmarshal(data, &hrvData); err == nil {
				if avg, ok := hrvData["average"].(float64); ok {
					dashboard.Garmin.HRV = &HRVData{Average: avg}
				}
			}
		case models.EventTypeGarminStress:
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
					dashboard.Garmin.Stress = &StressData{Average: avgInt, Level: level}
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

// GetWeekTrends retrieves 7-day trend data.
func (r *Repository) GetWeekTrends(ctx context.Context, userID string) ([]TrendData, error) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -6)
	startOfWeek := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

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

	trends := make([]TrendData, 0, len(trendsByDate))
	for _, trend := range trendsByDate {
		trends = append(trends, *trend)
	}

	return trends, nil
}

// GetCorrelations calculates simple correlations between Garmin data and feelings.
func (r *Repository) GetCorrelations(ctx context.Context, userID string, days int) ([]CorrelationInsight, error) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -days)
	startTime := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

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

	byDate := make(map[string]*dailyData)

	for rows.Next() {
		var eventTime time.Time
		var eventType string
		var data []byte

		if err := rows.Scan(&eventTime, &eventType, &data); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		dateKey := eventTime.Format("2006-01-02")
		if _, exists := byDate[dateKey]; !exists {
			byDate[dateKey] = &dailyData{}
		}
		daily := byDate[dateKey]

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

	return calculateCorrelations(byDate), nil
}

func calculateCorrelations(byDate map[string]*dailyData) []CorrelationInsight {
	var insights []CorrelationInsight

	if insight := sleepEnergyCorrelation(byDate); insight != nil {
		insights = append(insights, *insight)
	}
	if insight := activityMoodCorrelation(byDate); insight != nil {
		insights = append(insights, *insight)
	}
	if insight := sleepFocusCorrelation(byDate); insight != nil {
		insights = append(insights, *insight)
	}

	return insights
}

func sleepEnergyCorrelation(byDate map[string]*dailyData) *CorrelationInsight {
	var withGood, withPoor []int
	for _, d := range byDate {
		if d.Feeling == nil || d.Sleep == nil {
			continue
		}
		if float64(d.Sleep.DurationMinutes)/60.0 >= 7.0 {
			withGood = append(withGood, d.Feeling.Energy)
		} else {
			withPoor = append(withPoor, d.Feeling.Energy)
		}
	}
	if len(withGood) < 5 || len(withPoor) < 5 {
		return nil
	}
	avgGood, avgPoor := average(withGood), average(withPoor)
	improvement := ((avgGood - avgPoor) / avgPoor) * 100
	if improvement < 5 {
		return nil
	}
	return &CorrelationInsight{
		Type:        "sleep_energy",
		Description: fmt.Sprintf("Your energy is %.0f%% higher when you sleep 7+ hours", improvement),
		Confidence:  0.85,
		SampleSize:  len(withGood) + len(withPoor),
		Details: map[string]interface{}{
			"condition": "sleep >= 7 hours", "avg_energy_with": avgGood,
			"avg_energy_without": avgPoor, "improvement_percent": improvement,
		},
	}
}

func activityMoodCorrelation(byDate map[string]*dailyData) *CorrelationInsight {
	var withActivity, withoutActivity []int
	for _, d := range byDate {
		if d.Feeling == nil || d.Activity == nil {
			continue
		}
		if d.Activity.DurationMinutes >= 30 {
			withActivity = append(withActivity, d.Feeling.Mood)
		} else {
			withoutActivity = append(withoutActivity, d.Feeling.Mood)
		}
	}
	if len(withActivity) < 5 || len(withoutActivity) < 5 {
		return nil
	}
	avgWith, avgWithout := average(withActivity), average(withoutActivity)
	improvement := ((avgWith - avgWithout) / avgWithout) * 100
	if improvement < 5 {
		return nil
	}
	return &CorrelationInsight{
		Type:        "activity_mood",
		Description: fmt.Sprintf("Your mood improves by %.0f%% on active days (30+ min)", improvement),
		Confidence:  0.78,
		SampleSize:  len(withActivity) + len(withoutActivity),
		Details: map[string]interface{}{
			"condition": "activity >= 30 minutes", "avg_mood_with": avgWith,
			"avg_mood_without": avgWithout, "improvement_percent": improvement,
		},
	}
}

func sleepFocusCorrelation(byDate map[string]*dailyData) *CorrelationInsight {
	var withGood, withPoor []int
	for _, d := range byDate {
		if d.Feeling == nil || d.Sleep == nil {
			continue
		}
		if d.Sleep.SleepScore >= 80 {
			withGood = append(withGood, d.Feeling.Focus)
		} else {
			withPoor = append(withPoor, d.Feeling.Focus)
		}
	}
	if len(withGood) < 5 || len(withPoor) < 5 {
		return nil
	}
	avgGood, avgPoor := average(withGood), average(withPoor)
	improvement := ((avgGood - avgPoor) / avgPoor) * 100
	if improvement < 5 {
		return nil
	}
	return &CorrelationInsight{
		Type:        "sleep_focus",
		Description: fmt.Sprintf("Your focus is %.0f%% better after quality sleep (score 80+)", improvement),
		Confidence:  0.82,
		SampleSize:  len(withGood) + len(withPoor),
		Details: map[string]interface{}{
			"condition": "sleep_score >= 80", "avg_focus_with": avgGood,
			"avg_focus_without": avgPoor, "improvement_percent": improvement,
		},
	}
}

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
