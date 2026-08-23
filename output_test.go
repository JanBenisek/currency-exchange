package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var updateGolden = flag.Bool("update-golden", false, "Update golden test files")

// captureOutput captures stdout during function execution
func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

// normalizeOutput normalizes output for comparison
func normalizeOutput(output string) string {
	// Remove ANSI codes using the function from main.go
	output = stripANSI(output)
	// Trim trailing whitespace per line
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}
	return strings.Join(lines, "\n")
}

func TestTableFormatAlignment(t *testing.T) {
	// Mock rates for consistent testing
	rates := map[string]float64{
		"CZK": 25.7700,
		"EUR": 1.0700,
		"USD": 1.0,
	}
	amount := 761000.0
	base := "CHF"
	targets := []string{"CZK", "EUR"}

	// Initialize styles for consistent output
	themeName = "ocean"
	initStyles(themeName)

	output := captureOutput(func() {
		printTable(targets, rates, amount, base)
	})

	normalized := normalizeOutput(output)
	lines := strings.Split(normalized, "\n")

	// Check that the header row is present
	if !strings.Contains(normalized, "CURRENCY") {
		t.Error("Table missing 'CURRENCY' header")
	}
	if !strings.Contains(normalized, "CONVERTED") {
		t.Error("Table missing 'CONVERTED' header")
	}
	if !strings.Contains(normalized, "RATE") {
		t.Error("Table missing 'RATE' header")
	}

	// Verify data rows are present
	if !strings.Contains(normalized, "CZK") {
		t.Error("Table missing 'CZK' data")
	}
	if !strings.Contains(normalized, "EUR") {
		t.Error("Table missing 'EUR' data")
	}

	// Verify go-pretty rounded border characters are present
	if !strings.Contains(normalized, "╭") {
		t.Error("Table missing top-left border character")
	}
	if !strings.Contains(normalized, "╮") {
		t.Error("Table missing top-right border character")
	}
	if !strings.Contains(normalized, "╰") {
		t.Error("Table missing bottom-left border character")
	}
	if !strings.Contains(normalized, "╯") {
		t.Error("Table missing bottom-right border character")
	}

	// Verify column separators (│) are present
	hasColumnSeparators := false
	for _, line := range lines {
		if strings.Contains(line, "│") && strings.Contains(line, "CURRENCY") {
			hasColumnSeparators = true
			break
		}
	}
	if !hasColumnSeparators {
		t.Error("Table missing column separators")
	}
}

func TestMinimalFormatAlignment(t *testing.T) {
	rates := map[string]float64{
		"CZK": 25.7700,
		"EUR": 1.0700,
	}
	amount := 761000.0
	base := "CHF"
	targets := []string{"CZK", "EUR"}

	themeName = "ocean"
	initStyles(themeName)

	output := captureOutput(func() {
		printMinimal(targets, rates, amount, base)
	})

	normalized := normalizeOutput(output)
	lines := strings.Split(normalized, "\n")

	// Find lines with "@" symbol and verify they align
	var atSymbolPositions []int
	for i, line := range lines {
		if strings.Contains(line, "@") {
			pos := strings.Index(line, "@")
			atSymbolPositions = append(atSymbolPositions, pos)
			t.Logf("Line %d: '@' at position %d: %q", i, pos, line)
		}
	}

	// All "@" symbols should be at the same position (proper alignment)
	if len(atSymbolPositions) > 1 {
		firstPos := atSymbolPositions[0]
		for i, pos := range atSymbolPositions {
			if pos != firstPos {
				t.Errorf("'@' symbol misaligned: first at %d, line %d at %d",
					firstPos, i+1, pos)
			}
		}
	}
}

func TestCSVFormat(t *testing.T) {
	rates := map[string]float64{
		"CZK": 25.7700,
		"EUR": 1.0700,
	}
	amount := 100.0
	base := "USD"
	targets := []string{"CZK", "EUR"}

	output := captureOutput(func() {
		printCSV(targets, rates, amount, base)
	})

	lines := strings.Split(strings.TrimSpace(output), "\n")

	// First line should be header with new column order
	expectedHeader := "amount;currency_from;currency_to;rate;converted;date"
	if lines[0] != expectedHeader {
		t.Errorf("CSV header incorrect: got %q, want %q", lines[0], expectedHeader)
	}

	// Should have header + 2 data lines
	expectedLines := 3
	if len(lines) != expectedLines {
		t.Errorf("Expected %d lines, got %d", expectedLines, len(lines))
	}

	// Verify data lines have 6 semicolon-separated fields
	for i, line := range lines {
		if i == 0 {
			continue // Skip header
		}
		parts := strings.Split(line, ";")
		if len(parts) != 6 {
			t.Errorf("CSV line %d should have 6 semicolon-separated fields, got %d: %q",
				i, len(parts), line)
		}
		// Verify column order: amount;currency_from;currency_to;rate;converted;date
		if parts[1] != "USD" { // currency_from
			t.Errorf("Expected currency_from=USD at position 1, got %s", parts[1])
		}
		if parts[2] != "CZK" && parts[2] != "EUR" { // currency_to
			t.Errorf("Expected currency_to at position 2, got %s", parts[2])
		}
	}
}

func TestNumberFormat(t *testing.T) {
	rates := map[string]float64{
		"EUR": 0.85,
		"GBP": 0.73,
	}
	amount := 100.0
	base := "USD"
	targets := []string{"EUR", "GBP"}

	output := captureOutput(func() {
		printNumber(targets, rates, amount, base)
	})

	output = strings.TrimSpace(output)

	// Should output two space-separated numbers
	parts := strings.Split(output, " ")
	if len(parts) != 2 {
		t.Errorf("Expected 2 numbers, got %d: %q", len(parts), output)
	}

	// Numbers should not have thousand separators
	if strings.Contains(parts[0], ",") {
		t.Errorf("Numbers should not have thousand separators: got %q", parts[0])
	}

	// Each part should be a valid number
	for _, part := range parts {
		var f float64
		_, err := fmt.Sscanf(part, "%f", &f)
		if err != nil {
			t.Errorf("Invalid number format: %q (%v)", part, err)
		}
	}
}

func TestGoldenTableFormat(t *testing.T) {
	rates := map[string]float64{
		"CZK": 25.7700,
		"EUR": 1.0700,
	}
	amount := 1000.0
	base := "USD"
	targets := []string{"CZK", "EUR"}

	themeName = "ocean"
	initStyles(themeName)

	output := captureOutput(func() {
		printTable(targets, rates, amount, base)
	})

	normalized := normalizeOutput(output)
	goldenFile := "testdata/table.golden"

	// Update golden file if flag is set
	if *updateGolden {
		os.WriteFile(goldenFile, []byte(normalized), 0644)
		t.Logf("Updated golden file: %s", goldenFile)
		return
	}

	// Read golden file
	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatalf("Failed to read golden file: %v", err)
	}

	// Compare
	if normalized != string(expected) {
		t.Errorf("Table output does not match golden file")

		// Show diff
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(expected), normalized, false)
		t.Logf("Diff:\n%s", dmp.DiffPrettyText(diffs))
	}
}

func TestGoldenMinimalFormat(t *testing.T) {
	rates := map[string]float64{
		"CZK": 25.7700,
		"EUR": 1.0700,
	}
	amount := 1000.0
	base := "USD"
	targets := []string{"CZK", "EUR"}

	themeName = "ocean"
	initStyles(themeName)

	output := captureOutput(func() {
		printMinimal(targets, rates, amount, base)
	})

	normalized := normalizeOutput(output)
	goldenFile := "testdata/minimal.golden"

	// Update golden file if flag is set
	if *updateGolden {
		os.WriteFile(goldenFile, []byte(normalized), 0644)
		t.Logf("Updated golden file: %s", goldenFile)
		return
	}

	// Read golden file
	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatalf("Failed to read golden file: %v", err)
	}

	// Compare
	if normalized != string(expected) {
		t.Errorf("Minimal output does not match golden file")

		// Show diff
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(expected), normalized, false)
		t.Logf("Diff:\n%s", dmp.DiffPrettyText(diffs))
	}
}

// Test edge cases for alignment
func TestTableFormatLongNumbers(t *testing.T) {
	// Test with realistic numbers that fit in column widths
	rates := map[string]float64{
		"JPY": 145.30,
		"GBP": 0.79,
	}
	amount := 999999.0 // Large but manageable number
	base := "USD"
	targets := []string{"JPY", "GBP"}

	themeName = "ocean"
	initStyles(themeName)

	output := captureOutput(func() {
		printTable(targets, rates, amount, base)
	})

	normalized := normalizeOutput(output)

	// Verify the table contains the expected data
	if !strings.Contains(normalized, "JPY") {
		t.Error("Table should contain JPY")
	}
	if !strings.Contains(normalized, "GBP") {
		t.Error("Table should contain GBP")
	}
	if !strings.Contains(normalized, "145") {
		t.Error("Table should contain JPY rate")
	}
	if !strings.Contains(normalized, "0.79") {
		t.Error("Table should contain GBP rate")
	}
}

// TestTableColumnSeparatorPositions tests that column separators (│) are in consistent positions
func TestTableColumnSeparatorPositions(t *testing.T) {
	rates := map[string]float64{
		"EUR": 0.85,
		"GBP": 0.75,
	}
	amount := 100.0
	base := "USD"
	targets := []string{"EUR", "GBP"}

	themeName = "ocean"
	initStyles(themeName)

	output := captureOutput(func() {
		printTable(targets, rates, amount, base)
	})

	normalized := normalizeOutput(output)
	lines := strings.Split(normalized, "\n")

	// Find positions of all │ separators in each line
	type separatorPositions struct {
		first  int // Position of first │
		second int // Position of second │
		third  int // Position of third │
		fourth int // Position of fourth │
	}

	var positions []separatorPositions

	for _, line := range lines {
		if strings.Contains(line, "│") && line != "" {
			// Find all │ positions
			pipes := []int{}
			for i, c := range line {
				if c == '│' {
					pipes = append(pipes, i)
				}
			}

			if len(pipes) == 4 {
				positions = append(positions, separatorPositions{
					first:  pipes[0],
					second: pipes[1],
					third:  pipes[2],
					fourth: pipes[3],
				})
			}
		}
	}

	if len(positions) == 0 {
		t.Fatal("No properly formatted table rows found")
	}

	// All rows should have separator │ at the exact same positions
	firstPos := positions[0]
	for i, pos := range positions {
		if pos.first != firstPos.first {
			t.Errorf("Row %d: first │ at position %d, expected %d", i+1, pos.first, firstPos.first)
		}
		if pos.second != firstPos.second {
			t.Errorf("Row %d: second │ at position %d, expected %d", i+1, pos.second, firstPos.second)
		}
		if pos.third != firstPos.third {
			t.Errorf("Row %d: third │ at position %d, expected %d", i+1, pos.third, firstPos.third)
		}
		if pos.fourth != firstPos.fourth {
			t.Errorf("Row %d: fourth │ at position %d, expected %d", i+1, pos.fourth, firstPos.fourth)
		}
	}

	t.Logf("All column separators aligned at positions: %d, %d, %d, %d",
		firstPos.first, firstPos.second, firstPos.third, firstPos.fourth)
}

func TestMinimalFormatSingleCurrency(t *testing.T) {
	rates := map[string]float64{
		"EUR": 0.85,
	}
	amount := 100.0
	base := "USD"
	targets := []string{"EUR"}

	themeName = "ocean"
	initStyles(themeName)

	output := captureOutput(func() {
		printMinimal(targets, rates, amount, base)
	})

	normalized := normalizeOutput(output)

	// Should contain the currency and converted amount
	if !strings.Contains(normalized, "EUR") {
		t.Error("Output should contain 'EUR'")
	}
	if !strings.Contains(normalized, "@") {
		t.Error("Output should contain '@' separator")
	}
}
