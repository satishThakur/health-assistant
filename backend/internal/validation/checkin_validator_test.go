package validation

import (
	"strings"
	"testing"
)

func TestValidateCheckinPayload(t *testing.T) {
	tests := []struct {
		name    string
		payload *CheckinPayload
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid payload with all fields",
			payload: &CheckinPayload{
				Energy:   8,
				Mood:     7,
				Focus:    9,
				Physical: 7,
				Notes:    "Felt great today",
			},
			wantErr: false,
		},
		{
			name: "valid payload without notes",
			payload: &CheckinPayload{
				Energy:   5,
				Mood:     6,
				Focus:    7,
				Physical: 8,
			},
			wantErr: false,
		},
		{
			name: "valid payload with min values",
			payload: &CheckinPayload{
				Energy:   1,
				Mood:     1,
				Focus:    1,
				Physical: 1,
			},
			wantErr: false,
		},
		{
			name: "valid payload with max values",
			payload: &CheckinPayload{
				Energy:   10,
				Mood:     10,
				Focus:    10,
				Physical: 10,
			},
			wantErr: false,
		},
		{
			name:    "nil payload",
			payload: nil,
			wantErr: true,
			errMsg:  "payload cannot be nil",
		},
		{
			name: "energy too low",
			payload: &CheckinPayload{
				Energy:   0,
				Mood:     7,
				Focus:    8,
				Physical: 9,
			},
			wantErr: true,
			errMsg:  "energy must be between 1 and 10",
		},
		{
			name: "energy too high",
			payload: &CheckinPayload{
				Energy:   11,
				Mood:     7,
				Focus:    8,
				Physical: 9,
			},
			wantErr: true,
			errMsg:  "energy must be between 1 and 10",
		},
		{
			name: "mood too low",
			payload: &CheckinPayload{
				Energy:   7,
				Mood:     -1,
				Focus:    8,
				Physical: 9,
			},
			wantErr: true,
			errMsg:  "mood must be between 1 and 10",
		},
		{
			name: "mood too high",
			payload: &CheckinPayload{
				Energy:   7,
				Mood:     12,
				Focus:    8,
				Physical: 9,
			},
			wantErr: true,
			errMsg:  "mood must be between 1 and 10",
		},
		{
			name: "focus too low",
			payload: &CheckinPayload{
				Energy:   7,
				Mood:     8,
				Focus:    0,
				Physical: 9,
			},
			wantErr: true,
			errMsg:  "focus must be between 1 and 10",
		},
		{
			name: "focus too high",
			payload: &CheckinPayload{
				Energy:   7,
				Mood:     8,
				Focus:    15,
				Physical: 9,
			},
			wantErr: true,
			errMsg:  "focus must be between 1 and 10",
		},
		{
			name: "physical too low",
			payload: &CheckinPayload{
				Energy:   7,
				Mood:     8,
				Focus:    9,
				Physical: 0,
			},
			wantErr: true,
			errMsg:  "physical must be between 1 and 10",
		},
		{
			name: "physical too high",
			payload: &CheckinPayload{
				Energy:   7,
				Mood:     8,
				Focus:    9,
				Physical: 11,
			},
			wantErr: true,
			errMsg:  "physical must be between 1 and 10",
		},
		{
			name: "notes too long",
			payload: &CheckinPayload{
				Energy:   7,
				Mood:     8,
				Focus:    9,
				Physical: 7,
				Notes:    strings.Repeat("a", 1001),
			},
			wantErr: true,
			errMsg:  "notes cannot exceed 1000 characters",
		},
		{
			name: "notes exactly 1000 characters",
			payload: &CheckinPayload{
				Energy:   7,
				Mood:     8,
				Focus:    9,
				Physical: 7,
				Notes:    strings.Repeat("a", 1000),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCheckinPayload(tt.payload)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateCheckinPayload() expected error but got nil")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateCheckinPayload() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateCheckinPayload() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestValidateScale(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		value     int
		wantErr   bool
	}{
		{"valid min", "energy", 1, false},
		{"valid mid", "energy", 5, false},
		{"valid max", "energy", 10, false},
		{"too low", "energy", 0, true},
		{"negative", "energy", -5, true},
		{"too high", "energy", 11, true},
		{"way too high", "energy", 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateScale(tt.fieldName, tt.value)

			if tt.wantErr && err == nil {
				t.Errorf("validateScale() expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("validateScale() unexpected error = %v", err)
			}
		})
	}
}
