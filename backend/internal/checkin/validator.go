package checkin

import (
	"errors"
	"fmt"
)

// Payload represents the request body for check-in submission.
type Payload struct {
	Energy   int    `json:"energy"`
	Mood     int    `json:"mood"`
	Focus    int    `json:"focus"`
	Physical int    `json:"physical"`
	Notes    string `json:"notes,omitempty"`
}

// ValidatePayload validates check-in submission data.
func ValidatePayload(payload *Payload) error {
	if payload == nil {
		return errors.New("payload cannot be nil")
	}

	if err := validateScale("energy", payload.Energy); err != nil {
		return err
	}
	if err := validateScale("mood", payload.Mood); err != nil {
		return err
	}
	if err := validateScale("focus", payload.Focus); err != nil {
		return err
	}
	if err := validateScale("physical", payload.Physical); err != nil {
		return err
	}

	if len(payload.Notes) > 1000 {
		return errors.New("notes cannot exceed 1000 characters")
	}

	return nil
}

func validateScale(fieldName string, value int) error {
	if value < 1 || value > 10 {
		return fmt.Errorf("%s must be between 1 and 10, got %d", fieldName, value)
	}
	return nil
}
