package youtube

import (
	"testing"
	"time"
)

func TestParseISO8601Duration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"P0D", 0},
		{"PT5M", 5 * time.Minute},
		{"PT1H30M", 90 * time.Minute},
		{"PT3M42S", 3*time.Minute + 42*time.Second},
		{"P1DT5H3M42S", 29*time.Hour + 3*time.Minute + 42*time.Second},
		{"P2DT0S", 48 * time.Hour},
		{"P1D", 24 * time.Hour},
		{"PT0S", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseISO8601Duration(tt.input)
			if err != nil {
				t.Fatalf("parseISO8601Duration(%q) returned error: %v", tt.input, err)
			}
			if got != tt.expected {
				t.Errorf("parseISO8601Duration(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
