package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const apiURL = "https://api.exchangerate-api.com/v4/latest/"

// RateResponse represents the API response
type RateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func getRates(base string) (map[string]float64, error) {
	resp, err := http.Get(apiURL + strings.ToUpper(base))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rateResp RateResponse
	if err := json.Unmarshal(body, &rateResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return rateResp.Rates, nil
}

func printTable(targets []string, rates map[string]float64, amount float64, base string) {
	// Print header
	fmt.Println("┌────────────┬─────────────┬────────────────────┐")
	fmt.Println("│ Currency   │   Converted │ Rate               │")
	fmt.Println("├────────────┼─────────────┼────────────────────┤")

	for _, curr := range targets {
		rate, ok := rates[strings.ToUpper(curr)]
		if !ok {
			fmt.Printf("│ %-10s │ %11s │ %-19s│\n", strings.ToUpper(curr), "N/A", "Currency not found")
			continue
		}
		converted := amount * rate
		fmt.Printf("│ %-10s │ %11.2f │ 1 %s = %.4f %s │\n", strings.ToUpper(curr), converted, strings.ToUpper(base), rate, strings.ToUpper(curr))
	}

	fmt.Println("└────────────┴─────────────┴────────────────────┘")
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: cex <amount> <from_currency> to <to_currency1> <to_currency2> ...")
		fmt.Println("Example: cex 100 czk to pln eur chf")
		os.Exit(1)
	}

	// Parse amount
	var amount float64
	if _, err := fmt.Sscanf(os.Args[1], "%f", &amount); err != nil {
		fmt.Printf("Invalid amount: %s\n", os.Args[1])
		os.Exit(1)
	}

	base := os.Args[2]

	// Skip "to" and get target currencies
	var targets []string
	for i := 4; i < len(os.Args); i++ {
		targets = append(targets, os.Args[i])
	}

	rates, err := getRates(base)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	printTable(targets, rates, amount, base)
}
