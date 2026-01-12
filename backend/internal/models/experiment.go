package models

import (
	"encoding/json"
	"time"
)

// Experiment represents a health experiment
type Experiment struct {
	ID               string          `json:"id" db:"id"`
	UserID           string          `json:"user_id" db:"user_id"`
	Name             string          `json:"name" db:"name"`
	Hypothesis       string          `json:"hypothesis" db:"hypothesis"`
	Status           string          `json:"status" db:"status"`
	Intervention     json.RawMessage `json:"intervention" db:"intervention"`
	ControlCondition json.RawMessage `json:"control_condition,omitempty" db:"control_condition"`
	DurationDays     int             `json:"duration_days" db:"duration_days"`
	StartDate        *time.Time      `json:"start_date,omitempty" db:"start_date"`
	EndDate          *time.Time      `json:"end_date,omitempty" db:"end_date"`
	ComplianceRate   *float64        `json:"compliance_rate,omitempty" db:"compliance_rate"`
	Results          json.RawMessage `json:"results,omitempty" db:"results"`
	PosteriorBeliefs json.RawMessage `json:"posterior_beliefs,omitempty" db:"posterior_beliefs"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

// ExperimentStatus constants
const (
	ExperimentStatusProposed  = "proposed"
	ExperimentStatusAccepted  = "accepted"
	ExperimentStatusActive    = "active"
	ExperimentStatusCompleted = "completed"
	ExperimentStatusAbandoned = "abandoned"
)

// ExperimentResults represents statistical outcomes
type ExperimentResults struct {
	Effects map[string]Effect `json:"effects"`
	Summary string            `json:"summary"`
}

// Effect represents the effect of an intervention on a metric
type Effect struct {
	Mean                 float64   `json:"mean"`
	StdDev               float64   `json:"std"`
	CredibleInterval95   [2]float64 `json:"credible_interval_95"`
	ProbabilityPositive  float64   `json:"probability_positive"`
}
