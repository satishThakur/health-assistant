package models

import (
	"encoding/json"
	"time"
)

// User represents a user in the system
type User struct {
	ID               string          `json:"id" db:"id"`
	Email            string          `json:"email" db:"email"`
	PasswordHash     string          `json:"-" db:"password_hash"` // Never expose in JSON
	GarminOAuthToken json.RawMessage `json:"-" db:"garmin_oauth_token"`
	Preferences      json.RawMessage `json:"preferences,omitempty" db:"preferences"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
}

// UserPreferences represents user preferences
type UserPreferences struct {
	TimeZone             string   `json:"timezone"`
	SubjectiveReminders  []string `json:"subjective_reminders"` // e.g., ["08:00", "22:00"]
	SupplementReminders  bool     `json:"supplement_reminders"`
	ExperimentNotifications bool  `json:"experiment_notifications"`
}
