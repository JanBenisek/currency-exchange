package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

const (
	apiURL      = "https://api.exchangerate-api.com/v4/latest/"
	httpTimeout = 15 * time.Second
)

// Version is set via ldflags during build
// If not set, defaults to "dev"
var Version = "dev"

// RateResponse represents the current ExchangeRate-API response.
type RateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

// ConversionResult represents one currency conversion for JSON output.
type ConversionResult struct {
	Currency  string `json:"currency"`
	Converted string `json:"converted"` // Use string to preserve formatting
	Rate      string `json:"rate"`      // Use string for consistent decimal places
}

// JSONOutput represents the complete JSON output structure.
type JSONOutput struct {
	Amount string             `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  []ConversionResult `json:"rates"`
}

// Config holds parsed flag values.
type Config struct {
	showVersion bool
	themeName   string
	format      string
	apiProvider string
	positional  []string
}

// Color styles using ANSI codes.
var (
	reset = "\033[0m"
	bold  = "\033[1m"

	// Ocean theme
	cyan         = "\033[0;36m"
	brightCyan   = "\033[1;36m"
	brightCyanBg = "\033[96m"
	blue         = "\033[0;34m"
	brightBlue   = "\033[1;34m"
	red          = "\033[0;31m"
	brightRed    = "\033[1;31m"

	// Sunset theme
	yellow         = "\033[0;33m"
	brightYellow   = "\033[1;33m"
	brightYellowBg = "\033[93m"
	orange         = "\033[38;5;208m"
	brightOrange   = "\033[1;38;5;208m"

	// Forest theme
	green         = "\033[0;32m"
	brightGreen   = "\033[1;32m"
	brightGreenBg = "\033[92m"
	lime          = "\033[38;5;154m"
	brightLime    = "\033[1;38;5;154m"
	brown         = "\033[38;5;130m"
	brightBrown   = "\033[1;38;5;130m"

	// Neon theme
	magenta         = "\033[0;35m"
	brightMagenta   = "\033[1;35m"
	brightMagentaBg = "\033[95m"
	brightGreenNeon = "\033[1;92m"
	brightBlueNeon  = "\033[1;94m"
)

// Style functions.
var (
	borderStyle func(string) string
	headerStyle func(string) string
	valueStyle  func(string) string
	rateStyle   func(string) string
	errorStyle  func(string) string
	titleStyle  func(string) string
)

// Command options.
var (
	themeName     = "ocean"
	displayFormat = "table"
)

func initStyles(themeName string) {
	switch themeName {
	case "ocean":
		borderStyle = func(s string) string { return cyan + bold + s + reset }
		headerStyle = func(s string) string { return brightCyanBg + bold + s + reset }
		valueStyle = func(s string) string { return brightBlue + bold + s + reset }
		rateStyle = func(s string) string { return blue + s + reset }
		errorStyle = func(s string) string { return brightRed + bold + s + reset }
		titleStyle = func(s string) string { return brightCyanBg + bold + s + reset }

	case "sunset":
		borderStyle = func(s string) string { return orange + bold + s + reset }
		headerStyle = func(s string) string { return brightYellowBg + bold + s + reset }
		valueStyle = func(s string) string { return brightOrange + bold + s + reset }
		rateStyle = func(s string) string { return brightYellow + s + reset }
		errorStyle = func(s string) string { return brightRed + bold + s + reset }
		titleStyle = func(s string) string { return brightYellowBg + bold + s + reset }

	case "forest":
		borderStyle = func(s string) string { return lime + bold + s + reset }
		headerStyle = func(s string) string { return brightGreenBg + bold + s + reset }
		valueStyle = func(s string) string { return brightLime + bold + s + reset }
		rateStyle = func(s string) string { return green + s + reset }
		errorStyle = func(s string) string { return brightBrown + bold + s + reset }
		titleStyle = func(s string) string { return brightGreenBg + bold + s + reset }

	case "neon":
		borderStyle = func(s string) string { return brightMagenta + bold + s + reset }
		headerStyle = func(s string) string { return brightMagentaBg + bold + s + reset }
		valueStyle = func(s string) string { return brightGreenNeon + bold + s + reset }
		rateStyle = func(s string) string { return brightBlueNeon + s + reset }
		errorStyle = func(s string) string { return brightYellow + bold + s + reset }
		titleStyle = func(s string) string { return brightMagentaBg + bold + s + reset }

	default:
		borderStyle = func(s string) string { return cyan + bold + s + reset }
		headerStyle = func(s string) string { return brightCyanBg + bold + s + reset }
		valueStyle = func(s string) string { return brightBlue + bold + s + reset }
		rateStyle = func(s string) string { return blue + s + reset }
		errorStyle = func(s string) string { return brightRed + bold + s + reset }
		titleStyle = func(s string) string { return brightCyanBg + bold + s + reset }
	}
}

// parseFlags parses flags from anywhere in the argument list.
// It extracts flags and returns a Config with the parsed values.
func parseFlags(args []string) Config {
	config := Config{
		showVersion: false,
		themeName:   "ocean",
		format:      "table",
		apiProvider: "exchangerates",
		positional:  make([]string, 0, len(args)),
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "-v" || arg == "--version":
			config.showVersion = true
		case arg == "-t" || arg == "-theme" || arg == "--theme":
			if i+1 < len(args) {
				config.themeName = args[i+1]
				i++ // skip next arg
			}
		case arg == "-f" || arg == "-format" || arg == "--format":
			if i+1 < len(args) {
				config.format = args[i+1]
				i++ // skip next arg
			}
		case arg == "-a" || arg == "-api" || arg == "--api":
			if i+1 < len(args) {
				config.apiProvider = args[i+1]
				i++ // skip next arg
			}
		case arg == "-help" || arg == "--help":
			// Help is handled before parseFlags is called
			config.positional = append(config.positional, arg)
		default:
			config.positional = append(config.positional, arg)
		}
	}

	return config
}

// parseAmount parses an amount and supports k.
func parseAmount(s string) (float64, error) {
	s = strings.ToLower(strings.TrimSpace(s))

	if strings.HasSuffix(s, "k") {
		numStr := strings.TrimSuffix(s, "k")

		num, err := strconv.ParseFloat(numStr, 64)

		if err != nil {
			return 0, fmt.Errorf(
				"invalid amount format: %s",
				s,
			)
		}

		return num * 1000, nil
	}

	return strconv.ParseFloat(s, 64)
}

// formatNumber formats a number with commas and two decimals.
func formatNumber(f float64) string {
	str := fmt.Sprintf("%.2f", f)

	parts := strings.Split(str, ".")

	intPart := parts[0]

	decPart := ""

	if len(parts) > 1 {
		decPart = "." + parts[1]
	}

	var result []rune

	count := 0

	for i := len(intPart) - 1; i >= 0; i-- {
		if count > 0 &&
			count%3 == 0 &&
			intPart[i] != '-' {
			result = append(
				[]rune{','},
				result...,
			)
		}

		result = append(
			[]rune{rune(intPart[i])},
			result...,
		)

		count++
	}

	return string(result) + decPart
}

// stripANSI removes ANSI color sequences.
func stripANSI(s string) string {
	for {
		start := strings.Index(s, "\033[")

		if start == -1 {
			return s
		}

		end := strings.Index(
			s[start:],
			"m",
		)

		if end == -1 {
			return s
		}

		end += start

		s = s[:start] + s[end+1:]
	}
}

// visibleWidth returns terminal-visible width.
func visibleWidth(s string) int {
	return len([]rune(stripANSI(s)))
}

// padANSI pads colored text to a fixed visible width.
func padANSI(s string, width int) string {
	current := visibleWidth(s)

	if current >= width {
		return s
	}

	return s + strings.Repeat(
		" ",
		width-current,
	)
}

// printTable displays the conversion in a formatted table using go-pretty.
func printTable(
	targets []string,
	rates map[string]float64,
	amount float64,
	base string,
) {
	// Print title row
	printTitle(amount, base)

	// Create table with go-pretty
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)

	// Set alignment - converted and rate columns right-aligned
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft, AlignFooter: text.AlignLeft, AlignHeader: text.AlignLeft},    // Currency
		{Number: 2, Align: text.AlignRight, AlignFooter: text.AlignRight, AlignHeader: text.AlignRight}, // Converted
		{Number: 3, Align: text.AlignRight, AlignFooter: text.AlignRight, AlignHeader: text.AlignRight}, // Rate
		// Enable max width to handle large numbers
		{Number: 2, WidthMax: 20},
		// Use transformers to apply colors - format first, then wrap with colors
		{Number: 1, Transformer: text.Transformer(func(val interface{}) string {
			if str, ok := val.(string); ok {
				// Check if this looks like a currency code (3 letters, all caps)
				if len(str) == 3 && str == strings.ToUpper(str) {
					return valueStyle(str)
				}
			}
			return fmt.Sprintf("%v", val)
		})},
		{Number: 2, Transformer: text.Transformer(func(val interface{}) string {
			return valueStyle(fmt.Sprintf("%v", val))
		})},
		{Number: 3, Transformer: text.Transformer(func(val interface{}) string {
			return rateStyle(fmt.Sprintf("%v", val))
		})},
	})

	// Set headers with colors
	t.AppendHeader(table.Row{
		headerStyle("Currency"),
		headerStyle("Converted"),
		headerStyle("Rate"),
	})
	t.AppendSeparator()

	// Build data rows (colors will be applied by transformers)
	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]

		if !ok {
			t.AppendRow(table.Row{currUpper, "N/A", "Not found"})
			continue
		}

		converted := amount * rate
		t.AppendRow(table.Row{
			currUpper,
			formatNumber(converted),
			fmt.Sprintf("%.4f", rate),
		})
	}

	// Render the table
	t.Render()
	fmt.Println()
}

// printTitle prints the title row with amount and base currency
func printTitle(amount float64, base string) {
	baseUpper := strings.ToUpper(base)
	amountStr := formatNumber(amount)
	fmt.Printf("%s %s\n", valueStyle(amountStr), headerStyle(baseUpper))
}

// printCSV displays the conversion in CSV format.
func printCSV(
	targets []string,
	rates map[string]float64,
	amount float64,
	base string,
) {
	dateStr := time.Now().Format("2006-01-02")
	baseUpper := strings.ToUpper(base)

	// Print header with new column order
	fmt.Println("amount;currency_from;currency_to;rate;converted;date")

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]

		if !ok {
			fmt.Printf("%s;%s;%s;%s;%s;%s\n", formatNumber(amount), baseUpper, currUpper, "N/A", "N/A", dateStr)
			continue
		}

		converted := amount * rate
		fmt.Printf("%s;%s;%s;%.4f;%s;%s\n", formatNumber(amount), baseUpper, currUpper, rate, formatNumber(converted), dateStr)
	}
}

// printJSON displays the conversion in JSON format.
func printJSON(
	targets []string,
	rates map[string]float64,
	amount float64,
	base string,
) {
	baseUpper := strings.ToUpper(base)
	dateStr := time.Now().Format("2006-01-02")

	output := JSONOutput{
		Amount: formatNumber(amount),
		Base:   baseUpper,
		Date:   dateStr,
		Rates:  []ConversionResult{},
	}

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]

		if !ok {
			// Skip currencies not found
			continue
		}

		converted := amount * rate
		output.Rates = append(output.Rates, ConversionResult{
			Currency:  currUpper,
			Converted: formatNumber(converted),
			Rate:      fmt.Sprintf("%.4f", rate),
		})
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Printf("Error: failed to generate JSON: %v\n", err)
		return
	}

	fmt.Println(string(jsonData))
}

// printNumber displays only the converted values as plain numbers.
func printNumber(
	targets []string,
	rates map[string]float64,
	amount float64,
	base string,
) {
	var results []string

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]

		if !ok {
			results = append(results, "N/A")
			continue
		}

		converted := amount * rate
		// Format without thousand separators
		results = append(results, fmt.Sprintf("%.2f", converted))
	}

	fmt.Println(strings.Join(results, " "))
}

// printMinimal displays the conversion in a minimal format.
func printMinimal(
	targets []string,
	rates map[string]float64,
	amount float64,
	base string,
) {
	baseUpper := strings.ToUpper(base)

	fmt.Println()
	fmt.Println(valueStyle(formatNumber(amount)) + " " + headerStyle(baseUpper))
	fmt.Println()

	// Find max width for converted values to align the "@" column
	maxConvertedWidth := 0
	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		if rate, ok := rates[currUpper]; ok {
			converted := amount * rate
			convertedStr := formatNumber(converted)
			if len(convertedStr) > maxConvertedWidth {
				maxConvertedWidth = len(convertedStr)
			}
		}
	}
	// Ensure minimum width for "N/A" case
	if maxConvertedWidth < 3 {
		maxConvertedWidth = 3
	}

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]

		if !ok {
			fmt.Printf("  %-5s  %s\n", errorStyle(currUpper), errorStyle("N/A"))
			continue
		}

		converted := amount * rate
		convertedStr := formatNumber(converted)
		// Right-align converted values with proper spacing before "@"
		fmt.Printf("  %-5s  %s  @  %s\n", valueStyle(currUpper), valueStyle(fmt.Sprintf("%"+strconv.Itoa(maxConvertedWidth)+"s", convertedStr)), rateStyle(fmt.Sprintf("%.4f", rate)))
	}

	fmt.Println()
}

// printResults chooses the output format.
func printResults(
	targets []string,
	rates map[string]float64,
	amount float64,
	base string,
	format string,
) {
	switch format {
	case "csv":
		printCSV(
			targets,
			rates,
			amount,
			base,
		)

	case "json":
		printJSON(
			targets,
			rates,
			amount,
			base,
		)

	case "number":
		printNumber(
			targets,
			rates,
			amount,
			base,
		)

	case "minimal":
		printMinimal(
			targets,
			rates,
			amount,
			base,
		)

	default:
		printTable(
			targets,
			rates,
			amount,
			base,
		)
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func printHelp() {
	fmt.Println("Currency Exchange CLI - Convert between currencies")
	fmt.Println()
	fmt.Println("VERSION:")
	fmt.Printf("  cex %s\n", Version)
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  cex [options] <amount> <from_currency> to <to_currency1> <to_currency2> ...")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  cex 100 czk to pln eur usd")
	fmt.Println("  cex 700k chf to czk")
	fmt.Println("  cex -t sunset 50 usd to gbp jpy cad")
	fmt.Println("  cex -a fixer 100 usd to eur")
	fmt.Println("  cex 100 -v usd to eur")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -v, --version")
	fmt.Println("      Show version information.")
	fmt.Println()
	fmt.Println("  -t, -theme, --theme string")
	fmt.Println("      Color theme (default \"ocean\")")
	fmt.Println("      Available: ocean, sunset, forest, neon")
	fmt.Println()
	fmt.Println("  -f, -format, --format string")
	fmt.Println("      Output format (default \"table\")")
	fmt.Println("      Available: table, minimal, csv, json, number")
	fmt.Println()
	fmt.Println("  -a, -api, --api string")
	fmt.Println("      API provider (default \"exchangerates\")")
	fmt.Println("      Available: exchangerates, fixer, currencylayer, openexchangerates")
	fmt.Println()
	fmt.Println("  -help, --help")
	fmt.Println("      Show this help.")
	fmt.Println()
	fmt.Println("THEME PREVIEW:")
	fmt.Println("  Ocean:   Blues and cyans - professional look")
	fmt.Println("  Sunset:  Warm tones - yellows and oranges")
	fmt.Println("  Forest:  Greens and earthy tones")
	fmt.Println("  Neon:    Bright high-contrast colors")
	fmt.Println()
	fmt.Println("FORMAT PREVIEW:")
	fmt.Println("  table:   Formatted table with headers")
	fmt.Println("  minimal: Compact format with currency codes")
	fmt.Println("  csv:     Semicolon-delimited values with header")
	fmt.Println("  json:    Structured JSON output")
	fmt.Println("  number:  Only converted values (space-separated numbers)")
	fmt.Println()
	fmt.Println("API PROVIDERS:")
	fmt.Println("  exchangerates:      Exchange Rates API (free, no API key required)")
	fmt.Println("  fixer:              Fixer.io (requires FIXER_API_KEY env var)")
	fmt.Println("  currencylayer:      CurrencyLayer (requires CURRENCYLAYER_API_KEY env var)")
	fmt.Println("  openexchangerates:  Open Exchange Rates (requires OPENEXCHANGERATES_API_KEY env var)")
	fmt.Println()
	fmt.Println("AMOUNT FORMAT:")
	fmt.Println("  Supports 'k' suffix for thousands")
	fmt.Println("  (e.g., 700k = 700,000)")
}

func main() {
	// Check for help flags first
	for _, arg := range os.Args[1:] {
		if arg == "-help" || arg == "--help" {
			printHelp()
			return
		}
	}

	// Parse flags (they can appear anywhere in the command line)
	config := parseFlags(os.Args[1:])

	// Handle version flag
	if config.showVersion {
		fmt.Printf("cex %s\n", Version)
		return
	}

	// Initialize styles with the parsed theme
	themeName = config.themeName
	displayFormat = config.format

	// Validate theme value
	validThemes := map[string]bool{
		"ocean":  true,
		"sunset": true,
		"forest": true,
		"neon":   true,
	}
	if !validThemes[themeName] {
		fmt.Printf("Error: invalid theme '%s'\n", themeName)
		fmt.Println("Valid themes: ocean, sunset, forest, neon")
		os.Exit(1)
	}

	// Validate format value
	validFormats := map[string]bool{
		"table":   true,
		"minimal": true,
		"csv":     true,
		"json":    true,
		"number":  true,
	}
	if !validFormats[displayFormat] {
		fmt.Printf("Error: invalid format '%s'\n", displayFormat)
		fmt.Println("Valid formats: table, minimal, csv, json, short")
		os.Exit(1)
	}

	// Validate API provider value
	validProviders := map[string]bool{
		"exchangerates":     true,
		"fixer":             true,
		"currencylayer":     true,
		"openexchangerates": true,
	}
	if !validProviders[config.apiProvider] {
		fmt.Printf("Error: invalid API provider '%s'\n", config.apiProvider)
		fmt.Println("Valid providers: exchangerates, fixer, currencylayer, openexchangerates")
		os.Exit(1)
	}

	initStyles(themeName)

	args := config.positional

	// Validate we have at least 3 arguments (amount, base, to)
	if len(args) < 3 {
		printHelp()
		os.Exit(1)
	}

	// Parse amount
	amount, err := parseAmount(args[0])
	if err != nil {
		fmt.Printf("Error: invalid amount '%s'\n", args[0])
		fmt.Println("Amount should be a number (e.g., 100, 99.99)")
		fmt.Println("You can use 'k' suffix for thousands (e.g., 100k = 100,000)")
		os.Exit(1)
	}

	base := args[1]

	// Validate "to" keyword (if we have a third argument)
	if len(args) >= 3 && strings.ToLower(args[2]) != "to" {
		fmt.Println("Error: expected 'to' as third argument")
		fmt.Println()
		fmt.Println("Usage: cex [options] <amount> <from_currency> to <to_currency1> <to_currency2> ...")
		fmt.Println()
		fmt.Println("Example: cex 100 usd to eur gbp")
		os.Exit(1)
	}

	// Collect target currencies
	var targets []string
	for i := 3; i < len(args); i++ {
		targets = append(targets, args[i])
	}

	// Validate that we have at least one target currency
	if len(targets) == 0 {
		fmt.Println("Error: no target currencies specified")
		fmt.Println()
		fmt.Println("Usage: cex [options] <amount> <from_currency> to <to_currency1> <to_currency2> ...")
		fmt.Println()
		fmt.Println("Example: cex 100 usd to eur gbp")
		os.Exit(1)
	}

	// Create API provider
	provider, err := NewProvider(APIProvider(config.apiProvider))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println()
		if strings.Contains(err.Error(), "API_KEY") {
			fmt.Println("To use this provider, set the required API key as an environment variable:")
			switch config.apiProvider {
			case "fixer":
				fmt.Println("  export FIXER_API_KEY=your_key_here")
			case "currencylayer":
				fmt.Println("  export CURRENCYLAYER_API_KEY=your_key_here")
			case "openexchangerates":
				fmt.Println("  export OPENEXCHANGERATES_API_KEY=your_key_here")
			}
		}
		os.Exit(1)
	}

	// Fetch exchange rates
	rates, err := getRates(base, provider)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("Please use a valid 3-letter currency code (e.g., USD, EUR, GBP)")
		} else {
			fmt.Printf("Error: %v\n", err)
		}
		os.Exit(1)
	}

	// Check if base currency exists in the response
	baseUpper := strings.ToUpper(base)
	if _, exists := rates[baseUpper]; !exists && baseUpper != "USD" {
		// USD is the implicit base in most APIs, so it might not be in the rates map
		fmt.Printf("Error: currency '%s' not found\n", baseUpper)
		fmt.Println("Please use a valid 3-letter currency code (e.g., USD, EUR, GBP)")
		os.Exit(1)
	}

	// Count how many target currencies were actually found
	validTargets := 0
	for _, curr := range targets {
		if _, exists := rates[strings.ToUpper(curr)]; exists {
			validTargets++
		}
	}

	if validTargets == 0 {
		fmt.Println("Error: none of the specified currencies were found")
		fmt.Println("Please use valid 3-letter currency codes (e.g., USD, EUR, GBP)")
		os.Exit(1)
	}

	// Print results
	printResults(
		targets,
		rates,
		amount,
		base,
		displayFormat,
	)
}
