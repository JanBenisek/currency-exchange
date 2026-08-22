package main

import (
	"testing"
)

func TestParseAmount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{"regular number", "100", 100.0, false},
		{"with k suffix", "700k", 700000.0, false},
		{"decimal", "99.99", 99.99, false},
		{"k with decimal", "1.5k", 1500.0, false},
		{"invalid", "abc", 0, true},
		{"empty", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAmount(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("parseAmount() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{"small number", 100.0, "100.00"},
		{"thousands", 18046000.0, "18,046,000.00"},
		{"millions", 1234567.89, "1,234,567.89"},
		{"negative", -1000.0, "-1,000.00"},
		{"zero", 0.0, "0.00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatNumber(tt.input)
			if result != tt.expected {
				t.Errorf("formatNumber() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInitStyles(t *testing.T) {
	themes := []string{"ocean", "sunset", "forest", "neon", "invalid"}

	for _, theme := range themes {
		t.Run(theme, func(t *testing.T) {
			// This should not panic for any theme
			initStyles(theme)

			// Check that styles are initialized (not nil)
			if borderStyle == nil || headerStyle == nil || valueStyle == nil {
				t.Error("Styles were not initialized properly")
			}
		})
	}
}
