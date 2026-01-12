package models

import (
	"encoding/json"
	"time"
)

// Event represents a time-series event in the system
type Event struct {
	Time       time.Time       `json:"time" db:"time"`
	UserID     string          `json:"user_id" db:"user_id"`
	EventType  string          `json:"event_type" db:"event_type"`
	Source     string          `json:"source" db:"source"`
	Data       json.RawMessage `json:"data" db:"data"`
	Metadata   json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	Confidence *float64        `json:"confidence,omitempty" db:"confidence"`
}

// EventType constants
const (
	EventTypeGarminSleep      = "garmin_sleep"
	EventTypeGarminActivity   = "garmin_activity"
	EventTypeGarminHRV        = "garmin_hrv"
	EventTypeGarminStress     = "garmin_stress"
	EventTypeSubjectiveFeeling = "subjective_feeling"
	EventTypeMeal             = "meal"
	EventTypeSupplement       = "supplement"
	EventTypeBiomarker        = "biomarker"
)

// Source constants
const (
	SourceGarmin  = "garmin"
	SourceManual  = "manual"
	SourceParsed  = "parsed"
	SourceLLM     = "llm"
)

// SubjectiveFeeling represents daily subjective assessment
type SubjectiveFeeling struct {
	Energy   int    `json:"energy"`   // 1-10 scale
	Mood     int    `json:"mood"`     // 1-10 scale
	Focus    int    `json:"focus"`    // 1-10 scale
	Physical int    `json:"physical"` // 1-10 scale
	Notes    string `json:"notes,omitempty"`
}

// Meal represents a meal log
type Meal struct {
	MealType          string   `json:"meal_type"` // breakfast, lunch, dinner, snack
	PhotoURL          string   `json:"photo_url,omitempty"`
	Macros            Macros   `json:"macros"`
	Confidence        *float64 `json:"confidence,omitempty"`
	ManuallyVerified  bool     `json:"manually_verified"`
}

// Macros represents nutritional macros
type Macros struct {
	Calories  float64 `json:"calories"`
	ProteinG  float64 `json:"protein_g"`
	CarbsG    float64 `json:"carbs_g"`
	FatG      float64 `json:"fat_g"`
	FiberG    float64 `json:"fiber_g,omitempty"`
}

// Supplement represents a supplement log
type Supplement struct {
	Name          string    `json:"name"`
	Dosage        string    `json:"dosage"`
	Taken         bool      `json:"taken"`
	ScheduledTime string    `json:"scheduled_time"`
	ActualTime    time.Time `json:"actual_time,omitempty"`
}

// GarminSleep represents Garmin sleep data
type GarminSleep struct {
	DurationMinutes     int     `json:"duration_minutes"`
	DeepSleepMinutes    int     `json:"deep_sleep_minutes"`
	LightSleepMinutes   int     `json:"light_sleep_minutes"`
	REMSleepMinutes     int     `json:"rem_sleep_minutes"`
	AwakeMinutes        int     `json:"awake_minutes"`
	SleepScore          int     `json:"sleep_score"`
	HRVAvg              float64 `json:"hrv_avg,omitempty"`
}

// GarminActivity represents Garmin activity data
type GarminActivity struct {
	ActivityType string  `json:"activity_type"`
	DurationMinutes int  `json:"duration_minutes"`
	Calories     int     `json:"calories"`
	AvgHR        int     `json:"avg_hr,omitempty"`
	MaxHR        int     `json:"max_hr,omitempty"`
	Distance     float64 `json:"distance,omitempty"` // in meters
}

// Biomarker represents a lab test result
type Biomarker struct {
	TestName       string  `json:"test_name"`
	Value          float64 `json:"value"`
	Unit           string  `json:"unit"`
	ReferenceRange string  `json:"reference_range"`
	LabName        string  `json:"lab_name,omitempty"`
}
