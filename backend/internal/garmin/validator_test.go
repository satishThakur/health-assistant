package garmin

import (
	"testing"
)

func TestValidateSleepPayload(t *testing.T) {
	tests := []struct {
		name    string
		payload *SleepPayload
		wantErr bool
	}{
		{
			name: "valid payload",
			payload: &SleepPayload{
				UserID: "00000000-0000-0000-0000-000000000001",
				Date:   "2026-01-28",
				SleepData: map[string]interface{}{
					"sleep_time_seconds": float64(28800),
					"deep_sleep_seconds": float64(7200),
				},
			},
			wantErr: false,
		},
		{
			name: "missing user_id",
			payload: &SleepPayload{
				UserID: "",
				Date:   "2026-01-28",
				SleepData: map[string]interface{}{
					"sleep_time_seconds": float64(28800),
				},
			},
			wantErr: true,
		},
		{
			name: "missing date",
			payload: &SleepPayload{
				UserID: "00000000-0000-0000-0000-000000000001",
				Date:   "",
				SleepData: map[string]interface{}{
					"sleep_time_seconds": float64(28800),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid date format",
			payload: &SleepPayload{
				UserID: "00000000-0000-0000-0000-000000000001",
				Date:   "01/28/2026",
				SleepData: map[string]interface{}{
					"sleep_time_seconds": float64(28800),
				},
			},
			wantErr: true,
		},
		{
			name: "missing sleep_data",
			payload: &SleepPayload{
				UserID:    "00000000-0000-0000-0000-000000000001",
				Date:      "2026-01-28",
				SleepData: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid sleep_time_seconds",
			payload: &SleepPayload{
				UserID: "00000000-0000-0000-0000-000000000001",
				Date:   "2026-01-28",
				SleepData: map[string]interface{}{
					"sleep_time_seconds": float64(-100),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSleepPayload(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSleepPayload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateActivityPayload(t *testing.T) {
	tests := []struct {
		name    string
		payload *ActivityPayload
		wantErr bool
	}{
		{
			name: "valid payload",
			payload: &ActivityPayload{
				UserID: "00000000-0000-0000-0000-000000000001",
				Date:   "2026-01-28",
				ActivityData: map[string]interface{}{
					"activity_type":    "running",
					"duration_seconds": float64(2700),
				},
			},
			wantErr: false,
		},
		{
			name: "missing activity_type",
			payload: &ActivityPayload{
				UserID: "00000000-0000-0000-0000-000000000001",
				Date:   "2026-01-28",
				ActivityData: map[string]interface{}{
					"duration_seconds": float64(2700),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid duration",
			payload: &ActivityPayload{
				UserID: "00000000-0000-0000-0000-000000000001",
				Date:   "2026-01-28",
				ActivityData: map[string]interface{}{
					"activity_type":    "running",
					"duration_seconds": float64(0),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateActivityPayload(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateActivityPayload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateHRVPayload(t *testing.T) {
	tests := []struct {
		name    string
		payload *HRVPayload
		wantErr bool
	}{
		{
			name: "valid payload",
			payload: &HRVPayload{
				UserID: "00000000-0000-0000-0000-000000000001",
				Date:   "2026-01-28",
				HRVData: map[string]interface{}{
					"average_hrv": float64(65.5),
				},
			},
			wantErr: false,
		},
		{
			name: "negative hrv",
			payload: &HRVPayload{
				UserID: "00000000-0000-0000-0000-000000000001",
				Date:   "2026-01-28",
				HRVData: map[string]interface{}{
					"average_hrv": float64(-10),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHRVPayload(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHRVPayload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateStressPayload(t *testing.T) {
	tests := []struct {
		name    string
		payload *StressPayload
		wantErr bool
	}{
		{
			name: "valid payload",
			payload: &StressPayload{
				UserID: "00000000-0000-0000-0000-000000000001",
				Date:   "2026-01-28",
				StressData: map[string]interface{}{
					"average_stress_level": float64(45),
				},
			},
			wantErr: false,
		},
		{
			name: "stress level too high",
			payload: &StressPayload{
				UserID: "00000000-0000-0000-0000-000000000001",
				Date:   "2026-01-28",
				StressData: map[string]interface{}{
					"average_stress_level": float64(150),
				},
			},
			wantErr: true,
		},
		{
			name: "negative stress level",
			payload: &StressPayload{
				UserID: "00000000-0000-0000-0000-000000000001",
				Date:   "2026-01-28",
				StressData: map[string]interface{}{
					"average_stress_level": float64(-10),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStressPayload(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStressPayload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFloat64(t *testing.T) {
	tests := []struct {
		name   string
		data   map[string]interface{}
		key    string
		want   float64
		wantOk bool
	}{
		{
			name:   "float64 value",
			data:   map[string]interface{}{"key": float64(42.5)},
			key:    "key",
			want:   42.5,
			wantOk: true,
		},
		{
			name:   "int value",
			data:   map[string]interface{}{"key": int(42)},
			key:    "key",
			want:   42.0,
			wantOk: true,
		},
		{
			name:   "missing key",
			data:   map[string]interface{}{"other": float64(42)},
			key:    "key",
			want:   0,
			wantOk: false,
		},
		{
			name:   "invalid type",
			data:   map[string]interface{}{"key": "not a number"},
			key:    "key",
			want:   0,
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOk := getFloat64(tt.data, tt.key)
			if got != tt.want {
				t.Errorf("getFloat64() got = %v, want %v", got, tt.want)
			}
			if gotOk != tt.wantOk {
				t.Errorf("getFloat64() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
