package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const apiURL = "https://api.exchangerate-api.com/v4/latest/"

// RateResponse represents the API response
type RateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

// Color themes
type ColorTheme struct {
	Name   string
	Reset  string
	Header string
	Border string
	Value  string
	Rate   string
	Error  string
}

// ANSI color codes
const (
	reset   = "\033[0m"

	// Ocean theme (default)
	cyanBold      = "\033[1;36m"
	cyan          = "\033[0;36m"
	blueBold      = "\033[1;34m"
	blue          = "\033[0;34m"
	whiteBold     = "\033[1;37m"
	brightCyan    = "\033[96m"

	// Sunset theme
	yellowBold    = "\033[1;33m"
	yellow        = "\033[0;33m"
	orangeBold    = "\033[1;38;5;208m"
	orange        = "\033[0;38;5;208m"
	redBold       = "\033[1;31m"
	red           = "\033[0;31m"

	// Forest theme
	greenBold     = "\033[1;32m"
	green         = "\033[0;32m"
	limeBold      = "\033[1;38;5;154m"
	lime          = "\033[0;38;5;154m"
	brownBold     = "\033[1;38;5;130m"
	brown         = "\033[0;38;5;130m"

	// Neon theme
	magentaBold   = "\033[1;35m"
	magenta       = "\033[0;35m"
	brightMagenta = "\033[1;95m"
	brightGreen   = "\033[1;92m"
	brightYellow  = "\033[1;93m"
	brightBlue    = "\033[1;94m"
)

var themes = map[string]ColorTheme{
	"ocean": {
		Name:   "Ocean",
		Reset:  reset,
		Header: cyanBold,
		Border: cyan,
		Value:  blueBold,
		Rate:   blue,
		Error:  redBold,
	},
	"sunset": {
		Name:   "Sunset",
		Reset:  reset,
		Header: yellowBold,
		Border: orange,
		Value:  orangeBold,
		Rate:   yellow,
		Error:  redBold,
	},
	"forest": {
		Name:   "Forest",
		Reset:  reset,
		Header: greenBold,
		Border: lime,
		Value:  limeBold,
		Rate:   green,
		Error:  brownBold,
	},
	"neon": {
		Name:   "Neon",
		Reset:  reset,
		Header: brightMagenta,
		Border: magenta,
		Value:  brightGreen,
		Rate:   brightBlue,
		Error:  brightYellow,
	},
}

var (
	themeName   string
	displayFormat string
)

func init() {
	flag.StringVar(&themeName, "theme", "ocean", "Color theme: ocean, sunset, forest, neon")
	flag.StringVar(&displayFormat, "format", "table", "Display format: table, cards, list, minimal")
	flag.Parse()
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

// parseAmount parses an amount string with optional 'k' suffix (e.g., "700k" -> 700000)
func parseAmount(s string) (float64, error) {
	s = strings.ToLower(strings.TrimSpace(s))

	// Check for 'k' suffix
	if strings.HasSuffix(s, "k") {
		numStr := strings.TrimSuffix(s, "k")
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid amount format: %s", s)
		}
		return num * 1000, nil
	}

	// Parse as regular float
	return strconv.ParseFloat(s, 64)
}

// formatNumber formats a number with thousand separators
func formatNumber(f float64) string {
	// Split into integer and decimal parts
	str := fmt.Sprintf("%.2f", f)
	parts := strings.Split(str, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = "." + parts[1]
	}

	// Add thousand separators to integer part
	var result []rune
	count := 0
	for i := len(intPart) - 1; i >= 0; i-- {
		if count > 0 && count%3 == 0 && intPart[i] != '-' {
			result = append([]rune{','}, result...)
		}
		result = append([]rune{rune(intPart[i])}, result...)
		count++
		if !unicode.IsDigit(rune(intPart[i])) && intPart[i] != '-' {
			count = 0
		}
	}

	return string(result) + decPart
}

func printTable(targets []string, rates map[string]float64, amount float64, base string) {
	theme := themes[themeName]

	// Calculate column widths
	currencyColWidth := 10
	convertedColWidth := 19
	rateColWidth := 30

	// Print title with amount
	fmt.Printf("\n%s%s%s %s %.2f %s%s\n", theme.Header, "Currency Exchange", theme.Reset, "converting", amount, strings.ToUpper(base), theme.Reset)
	fmt.Println()

	// Print header
	borderWidth := currencyColWidth + convertedColWidth + rateColWidth + 10
	fmt.Printf("%s%s%s\n", theme.Border, strings.Repeat("─", borderWidth), theme.Reset)
	fmt.Printf("%s│%s %-10s %s│ %s%17s %s│ %s%-28s %s│\n",
		theme.Border, theme.Reset,
		theme.Header+"Currency"+theme.Reset,
		theme.Border, theme.Reset,
		theme.Header+"Converted"+theme.Reset,
		theme.Border, theme.Reset,
		theme.Header+"Rate (1 "+strings.ToUpper(base)+") = X"+theme.Reset,
		theme.Border)
	fmt.Printf("%s%s%s\n", theme.Border, strings.Repeat("─", borderWidth), theme.Reset)

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]
		if !ok {
			fmt.Printf("%s│%s %-10s %s│ %s%17s %s│ %s%-28s %s│\n",
				theme.Border, theme.Reset,
				currUpper,
				theme.Border, theme.Reset,
				theme.Error+"N/A"+theme.Reset,
				theme.Border, theme.Reset,
				"Currency not found",
				theme.Border)
			continue
		}
		converted := amount * rate
		fmt.Printf("%s│%s %-10s %s│ %s%17s %s│ %s1 %s = %.4f %-28s│\n",
			theme.Border, theme.Reset,
			currUpper,
			theme.Border, theme.Reset,
			theme.Value+formatNumber(converted)+theme.Reset,
			theme.Border, theme.Reset,
			strings.ToUpper(base), rate, theme.Rate+currUpper+theme.Reset)
	}

	fmt.Printf("%s%s%s\n", theme.Border, strings.Repeat("─", borderWidth), theme.Reset)
	fmt.Println()
}

func printCards(targets []string, rates map[string]float64, amount float64, base string) {
	theme := themes[themeName]

	fmt.Printf("\n%s%s%s %s %.2f %s%s\n\n", theme.Header, "Currency Exchange", theme.Reset, "converting", amount, strings.ToUpper(base), theme.Reset)

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]
		if !ok {
			fmt.Printf("%s┌────────────────────────────────┐%s\n", theme.Border, theme.Reset)
			fmt.Printf("%s│%s %s%-10s%s %s│%s\n", theme.Border, theme.Reset, theme.Header, currUpper, theme.Reset, strings.Repeat(" ", 17), theme.Border, theme.Reset)
			fmt.Printf("%s│%s %s%-30s%s %s│%s\n", theme.Border, theme.Reset, theme.Error, "Currency not found", theme.Reset, strings.Repeat(" ", 2), theme.Border, theme.Reset)
			fmt.Printf("%s└────────────────────────────────┘%s\n\n", theme.Border, theme.Reset)
			continue
		}

		converted := amount * rate
		fmt.Printf("%s┌────────────────────────────────┐%s\n", theme.Border, theme.Reset)
		fmt.Printf("%s│%s %-31s%s│\n", theme.Border, theme.Reset, theme.Header+currUpper+theme.Reset+strings.Repeat(" ", 31-len(currUpper)), theme.Border, theme.Reset)
		fmt.Printf("%s│%s Converted:  %s%-15s%s   %s│\n", theme.Border, theme.Reset, theme.Value, formatNumber(converted), theme.Reset, theme.Border, theme.Reset)
		fmt.Printf("%s│%s Rate:       %s1 %s = %.4f %s%s%s   %s│\n", theme.Border, theme.Reset, theme.Rate, strings.ToUpper(base), rate, currUpper, strings.Repeat(" ", 10-len(currUpper)), theme.Border, theme.Reset)
		fmt.Printf("%s└────────────────────────────────┘%s\n\n", theme.Border, theme.Reset)
	}
}

func printList(targets []string, rates map[string]float64, amount float64, base string) {
	theme := themes[themeName]

	fmt.Printf("\n%s%s%s %s %.2f %s%s\n\n", theme.Header, "Currency Exchange", theme.Reset, "converting", amount, strings.ToUpper(base), theme.Reset)

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]
		if !ok {
			fmt.Printf("%s  %s%-6s%s → %sN/A%s\n", theme.Border, theme.Header, currUpper, theme.Reset, theme.Error, theme.Reset)
			continue
		}

		converted := amount * rate
		fmt.Printf("%s  %s%-6s%s → %s%-15s%s  %s(1 %s = %.4f %s)%s\n",
			theme.Border, theme.Header, currUpper, theme.Reset,
			theme.Value, formatNumber(converted), theme.Reset,
			theme.Rate, strings.ToUpper(base), rate, currUpper, theme.Reset)
	}
	fmt.Println()
}

func printMinimal(targets []string, rates map[string]float64, amount float64, base string) {
	theme := themes[themeName]

	fmt.Printf("\n%s%.2f %s%s →\n", theme.Header, amount, strings.ToUpper(base), theme.Reset)

	var results []string
	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]
		if !ok {
			results = append(results, fmt.Sprintf("%s%s%s=N/A", theme.Error, currUpper, theme.Reset))
			continue
		}
		converted := amount * rate
		results = append(results, fmt.Sprintf("%s%s%s=%s", theme.Value, currUpper, theme.Reset, formatNumber(converted)))
	}

	fmt.Printf("  %s\n\n", strings.Join(results, "  |  "))
}

func printResults(targets []string, rates map[string]float64, amount float64, base string) {
	switch displayFormat {
	case "cards":
		printCards(targets, rates, amount, base)
	case "list":
		printList(targets, rates, amount, base)
	case "minimal":
		printMinimal(targets, rates, amount, base)
	default: // table
		printTable(targets, rates, amount, base)
	}
}

func main() {
	// Check if -h or -help is passed
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "-help" || os.Args[1] == "--help") {
		printHelp()
		os.Exit(0)
	}

	// After flag parsing, we need to get the non-flag args
	args := flag.Args()
	if len(args) < 4 {
		printHelp()
		os.Exit(1)
	}

	// Parse amount (supports "k" suffix like "700k" -> 700000)
	amount, err := parseAmount(args[0])
	if err != nil {
		fmt.Printf("Invalid amount: %s\n", args[0])
		os.Exit(1)
	}

	base := args[1]

	// Check for "to" separator and get target currencies
	var targets []string
	if args[2] != "to" {
		fmt.Println("Error: expected 'to' as third argument")
		fmt.Println("Usage: cex <amount> <from_currency> to <to_currency1> <to_currency2> ...")
		os.Exit(1)
	}
	for i := 3; i < len(args); i++ {
		targets = append(targets, args[i])
	}

	rates, err := getRates(base)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	printResults(targets, rates, amount, base)
}

func printHelp() {
	fmt.Println("Currency Exchange CLI - Convert between currencies")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  cex <amount> <from_currency> to <to_currency1> [to_currency2] ...")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  cex 100 czk to pln eur usd")
	fmt.Println("  cex 50 usd to gbp jpy cad")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -theme string")
	fmt.Println("      Color theme (default \"ocean\")")
	fmt.Println("      Available: ocean, sunset, forest, neon")
	fmt.Println()
	fmt.Println("  -format string")
	fmt.Println("      Display format (default \"table\")")
	fmt.Println("      Available: table, cards, list, minimal")
	fmt.Println()
	fmt.Println("THEME PREVIEW:")
	fmt.Println("  Ocean:   Blues and cyans - professional look")
	fmt.Println("  Sunset:  Warm tones - yellows and oranges")
	fmt.Println("  Forest:  Greens and earthy tones")
	fmt.Println("  Neon:    Bright high-contrast colors")
}
