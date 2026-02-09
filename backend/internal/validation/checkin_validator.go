package validation

import (
	"errors"
	"fmt"
)

// CheckinPayload represents the request payload for check-in submission
type CheckinPayload struct {
	Energy   int    `json:"energy"`
	Mood     int    `json:"mood"`
	Focus    int    `json:"focus"`
	Physical int    `json:"physical"`
	Notes    string `json:"notes,omitempty"`
}

// ValidateCheckinPayload validates check-in submission data
func ValidateCheckinPayload(payload *CheckinPayload) error {
	if payload == nil {
		return errors.New("payload cannot be nil")
	}

	// Validate energy (1-10 scale)
	if err := validateScale("energy", payload.Energy); err != nil {
		return err
	}

	// Validate mood (1-10 scale)
	if err := validateScale("mood", payload.Mood); err != nil {
		return err
	}

	// Validate focus (1-10 scale)
	if err := validateScale("focus", payload.Focus); err != nil {
		return err
	}

	// Validate physical (1-10 scale)
	if err := validateScale("physical", payload.Physical); err != nil {
		return err
	}

	// Notes are optional, but if provided should not exceed 1000 characters
	if len(payload.Notes) > 1000 {
		return errors.New("notes cannot exceed 1000 characters")
	}

	return nil
}

// validateScale validates a value is between 1 and 10
func validateScale(fieldName string, value int) error {
	if value < 1 || value > 10 {
		return fmt.Errorf("%s must be between 1 and 10, got %d", fieldName, value)
	}
	return nil
}
