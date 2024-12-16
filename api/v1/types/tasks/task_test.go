package tasks

import (
	"testing"
	"time"
)

func TestConvertUTCString(t *testing.T) {

	testCases := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "Valid UTC string",
			input:   "2024-12-15T18:30:00.000+0000",
			want:    time.Date(2024, 12, 15, 18, 30, 0, 0, time.UTC).Local(), // Expected local time
			wantErr: false,
		},
		{
			name:    "Invalid time format",
			input:   "2024-12-15 18:30:00",
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   "",
			want:    time.Time{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				got := convertUTCString(tc.input)

				if (got == time.Time{}) != tc.wantErr {
					t.Errorf("convertUTCString() error = %v, wantErr %v", got, tc.wantErr)
					return
				}

				if !tc.wantErr && !got.Equal(tc.want) {
					t.Errorf("convertUTCString() got = %v, want %v", got, tc.want)
				}
			},
		)
	}
}

func TestConvertLocalTime(t *testing.T) {
	testCases := []struct {
		name  string
		input time.Time
		want  string
	}{
		{
			name:  "Local time to UTC string",
			input: time.Date(2024, 12, 15, 10, 30, 0, 0, time.Local), // Example local time
			want:  "2024-12-15T15:30:00.000+0000",                    // Expected UTC string
		},
		{
			name:  "Zero time",
			input: time.Time{},
			want:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				got := convertLocalTime(tc.input)
				if got != tc.want {
					t.Errorf("convertLocalTime() = %q, want %q", got, tc.want)
				}
			},
		)
	}
}
