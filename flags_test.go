package main

import (
	"reflect"
	"testing"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected Config
	}{
		{
			name:  "version flag",
			input: []string{"-v"},
			expected: Config{
				showVersion: true,
				themeName:   "ocean",
				format:      "table",
				positional:  []string{},
			},
		},
		{
			name:  "version long form",
			input: []string{"--version"},
			expected: Config{
				showVersion: true,
				themeName:   "ocean",
				format:      "table",
				positional:  []string{},
			},
		},
		{
			name:  "theme shorthand",
			input: []string{"-t", "sunset"},
			expected: Config{
				showVersion: false,
				themeName:   "sunset",
				format:      "table",
				positional:  []string{},
			},
		},
		{
			name:  "theme long form",
			input: []string{"--theme", "forest"},
			expected: Config{
				showVersion: false,
				themeName:   "forest",
				format:      "table",
				positional:  []string{},
			},
		},
		{
			name:  "format flag",
			input: []string{"--format", "minimal"},
			expected: Config{
				showVersion: false,
				themeName:   "ocean",
				format:      "minimal",
				positional:  []string{},
			},
		},
		{
			name:  "format shorthand",
			input: []string{"-f", "csv"},
			expected: Config{
				showVersion: false,
				themeName:   "ocean",
				format:      "csv",
				positional:  []string{},
			},
		},
		{
			name:  "format single dash",
			input: []string{"-format", "minimal"},
			expected: Config{
				showVersion: false,
				themeName:   "ocean",
				format:      "minimal",
				positional:  []string{},
			},
		},
		{
			name:  "theme single dash",
			input: []string{"-theme", "sunset"},
			expected: Config{
				showVersion: false,
				themeName:   "sunset",
				format:      "table",
				positional:  []string{},
			},
		},
		{
			name:  "positional args only",
			input: []string{"100", "usd", "to", "eur"},
			expected: Config{
				showVersion: false,
				themeName:   "ocean",
				format:      "table",
				positional:  []string{"100", "usd", "to", "eur"},
			},
		},
		{
			name:  "flags at beginning",
			input: []string{"-v", "100", "usd", "to", "eur"},
			expected: Config{
				showVersion: true,
				themeName:   "ocean",
				format:      "table",
				positional:  []string{"100", "usd", "to", "eur"},
			},
		},
		{
			name:  "flags in middle",
			input: []string{"100", "-t", "neon", "usd", "to", "eur"},
			expected: Config{
				showVersion: false,
				themeName:   "neon",
				format:      "table",
				positional:  []string{"100", "usd", "to", "eur"},
			},
		},
		{
			name:  "flags at end",
			input: []string{"100", "usd", "to", "eur", "--format", "cards"},
			expected: Config{
				showVersion: false,
				themeName:   "ocean",
				format:      "cards",
				positional:  []string{"100", "usd", "to", "eur"},
			},
		},
		{
			name:  "multiple flags",
			input: []string{"-t", "sunset", "--format", "csv", "100", "usd", "to", "eur"},
			expected: Config{
				showVersion: false,
				themeName:   "sunset",
				format:      "csv",
				positional:  []string{"100", "usd", "to", "eur"},
			},
		},
		{
			name:  "both shorthands",
			input: []string{"-t", "neon", "-f", "minimal", "100", "usd", "to", "eur"},
			expected: Config{
				showVersion: false,
				themeName:   "neon",
				format:      "minimal",
				positional:  []string{"100", "usd", "to", "eur"},
			},
		},
		{
			name:  "version with other flags",
			input: []string{"-t", "neon", "-v"},
			expected: Config{
				showVersion: true,
				themeName:   "neon",
				format:      "table",
				positional:  []string{},
			},
		},
		{
			name:  "empty args",
			input: []string{},
			expected: Config{
				showVersion: false,
				themeName:   "ocean",
				format:      "table",
				positional:  []string{},
			},
		},
		{
			name:  "unknown flag treated as positional",
			input: []string{"--unknown", "100", "usd", "to", "eur"},
			expected: Config{
				showVersion: false,
				themeName:   "ocean",
				format:      "table",
				positional:  []string{"--unknown", "100", "usd", "to", "eur"},
			},
		},
		{
			name:  "help flag treated as positional",
			input: []string{"-help", "100", "usd", "to", "eur"},
			expected: Config{
				showVersion: false,
				themeName:   "ocean",
				format:      "table",
				positional:  []string{"-help", "100", "usd", "to", "eur"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseFlags(tt.input)
			if result.showVersion != tt.expected.showVersion {
				t.Errorf("showVersion = %v, want %v", result.showVersion, tt.expected.showVersion)
			}
			if result.themeName != tt.expected.themeName {
				t.Errorf("themeName = %v, want %v", result.themeName, tt.expected.themeName)
			}
			if result.format != tt.expected.format {
				t.Errorf("format = %v, want %v", result.format, tt.expected.format)
			}
			if !reflect.DeepEqual(result.positional, tt.expected.positional) {
				t.Errorf("positional = %v, want %v", result.positional, tt.expected.positional)
			}
		})
	}
}

func TestParseFlagsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected Config
	}{
		{
			name:  "flag without value takes next arg",
			input: []string{"-t", "100", "usd", "to", "eur"},
			expected: Config{
				showVersion: false,
				themeName:   "100", // Will take "100" as theme name
				format:      "table",
				positional:  []string{"usd", "to", "eur"},
			},
		},
		{
			name:  "format at end of args",
			input: []string{"100", "usd", "to", "eur", "--format"},
			expected: Config{
				showVersion: false,
				themeName:   "ocean",
				format:      "table",                             // No value after --format, stays default
				positional:  []string{"100", "usd", "to", "eur"}, // --format is consumed but has no value
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseFlags(tt.input)
			if result.showVersion != tt.expected.showVersion {
				t.Errorf("showVersion = %v, want %v", result.showVersion, tt.expected.showVersion)
			}
			if result.themeName != tt.expected.themeName {
				t.Errorf("themeName = %v, want %v", result.themeName, tt.expected.themeName)
			}
			if result.format != tt.expected.format {
				t.Errorf("format = %v, want %v", result.format, tt.expected.format)
			}
			if !reflect.DeepEqual(result.positional, tt.expected.positional) {
				t.Errorf("positional = %v, want %v", result.positional, tt.expected.positional)
			}
		})
	}
}
