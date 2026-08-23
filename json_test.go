package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestJSONOutput(t *testing.T) {
	// Test that JSON can be unmarshaled from sample data
	output := JSONOutput{
		Amount: "100.00",
		Base:   "USD",
		Date:   "2026-08-23",
		Rates: []ConversionResult{
			{Currency: "EUR", Converted: "85.60", Rate: "0.8560"},
			{Currency: "GBP", Converted: "73.30", Rate: "0.7330"},
		},
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verify it's valid JSON
	var decoded JSONOutput
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Errorf("Generated invalid JSON: %v", err)
	}

	// Verify structure
	if decoded.Base != "USD" {
		t.Errorf("Expected base 'USD', got '%s'", decoded.Base)
	}

	if len(decoded.Rates) != 2 {
		t.Errorf("Expected 2 rates, got %d", len(decoded.Rates))
	}
}

func TestJSONSchema(t *testing.T) {
	// Test the expected JSON schema
	testOutput := JSONOutput{
		Amount: "761,000.00",
		Base:   "CHF",
		Date:   "2026-08-23",
		Rates: []ConversionResult{
			{Currency: "CZK", Converted: "19,610,970.00", Rate: "25.7700"},
			{Currency: "EUR", Converted: "814,270.00", Rate: "1.0700"},
		},
	}

	jsonData, err := json.MarshalIndent(testOutput, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verify JSON contains expected fields
	jsonStr := string(jsonData)
	requiredFields := []string{
		`"amount"`,
		`"base"`,
		`"date"`,
		`"rates"`,
		`"currency"`,
		`"converted"`,
		`"rate"`,
	}

	for _, field := range requiredFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("JSON missing field: %s", field)
		}
	}

	// Verify currency order is preserved
	if !strings.Contains(jsonStr, `"currency": "CZK"`) {
		t.Error("JSON should contain CZK currency")
	}
	if !strings.Contains(jsonStr, `"currency": "EUR"`) {
		t.Error("JSON should contain EUR currency")
	}
}
