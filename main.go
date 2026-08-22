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
)

const apiURL = "https://api.exchangerate-api.com/v4/latest/"

// RateResponse represents the API response
type RateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

// Color styles using ANSI codes
var (
	reset    = "\033[0m"
	bold     = "\033[1m"

	// Ocean theme
	cyan         = "\033[0;36m"
	brightCyan   = "\033[1;36m"
	brightCyanBg = "\033[96m"
	blue         = "\033[0;34m"
	brightBlue   = "\033[1;34m"
	red          = "\033[0;31m"
	brightRed    = "\033[1;31m"

	// Sunset theme
	yellow        = "\033[0;33m"
	brightYellow  = "\033[1;33m"
	brightYellowBg = "\033[93m"
	orange        = "\033[38;5;208m"
	brightOrange  = "\033[1;38;5;208m"

	// Forest theme
	green         = "\033[0;32m"
	brightGreen   = "\033[1;32m"
	brightGreenBg = "\033[92m"
	lime          = "\033[38;5;154m"
	brightLime    = "\033[1;38;5;154m"
	brown         = "\033[38;5;130m"
	brightBrown   = "\033[1;38;5;130m"

	// Neon theme
	magenta        = "\033[0;35m"
	brightMagenta  = "\033[1;35m"
	brightMagentaBg = "\033[95m"
	brightGreenNeon = "\033[1;92m"
	brightBlueNeon  = "\033[1;94m"
)

// Style functions
var (
	borderStyle   func(string) string
	headerStyle   func(string) string
	valueStyle    func(string) string
	rateStyle     func(string) string
	errorStyle    func(string) string
	titleStyle    func(string) string
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

func parseAmount(s string) (float64, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if strings.HasSuffix(s, "k") {
		numStr := strings.TrimSuffix(s, "k")
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid amount format: %s", s)
		}
		return num * 1000, nil
	}
	return strconv.ParseFloat(s, 64)
}

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
		if count > 0 && count%3 == 0 && intPart[i] != '-' {
			result = append([]rune{','}, result...)
		}
		result = append([]rune{rune(intPart[i])}, result...)
		count++
	}

	return string(result) + decPart
}

func printTable(targets []string, rates map[string]float64, amount float64, base string) {
	title := fmt.Sprintf("Currency Exchange converting %.2f %s", amount, strings.ToUpper(base))
	fmt.Println()
	fmt.Println(titleStyle(title))
	fmt.Println()

	// Define table borders - two column layout
	topBorder := borderStyle("┌─────────────────────┬────────────────────┐")
	headerBorder := borderStyle("├─────────────────────┼────────────────────┤")
	bottomBorder := borderStyle("└─────────────────────┴────────────────────┘")
	leftPipe := borderStyle("│")

	fmt.Println(topBorder)

	// Print header
	header := fmt.Sprintf("%s %-21s %-21s%s",
		leftPipe,
		headerStyle("Converted Value"),
		headerStyle(fmt.Sprintf("Rate (1 %s = X)", strings.ToUpper(base))),
		borderStyle("│"))
	fmt.Println(header)
	fmt.Println(headerBorder)

	// Print data rows
	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]
		if !ok {
			row := fmt.Sprintf("%s %-21s %-21s%s",
				leftPipe,
				errorStyle("N/A"),
				errorStyle("Currency not found"),
				borderStyle("│"))
			fmt.Println(row)
			continue
		}

		converted := amount * rate
		valueWithCurrency := fmt.Sprintf("%s %s", formatNumber(converted), currUpper)
		row := fmt.Sprintf("%s %-21s %-21s%s",
			leftPipe,
			valueStyle(valueWithCurrency),
			rateStyle(fmt.Sprintf("1 %s = %.4f %s", strings.ToUpper(base), rate, currUpper)),
			borderStyle("│"))
		fmt.Println(row)
	}

	fmt.Println(bottomBorder)
	fmt.Println()
}

func printCards(targets []string, rates map[string]float64, amount float64, base string) {
	title := fmt.Sprintf("Currency Exchange converting %.2f %s", amount, strings.ToUpper(base))
	fmt.Println()
	fmt.Println(titleStyle(title))
	fmt.Println()

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]
		if !ok {
			fmt.Println(borderStyle("┌────────────────────────────────┐"))
			fmt.Printf("%s %s%s\n", borderStyle("│"), errorStyle(fmt.Sprintf("Not found: %s", currUpper)), borderStyle(" │"))
			fmt.Println(borderStyle("└────────────────────────────────┘"))
			fmt.Println()
			continue
		}

		converted := amount * rate
		fmt.Println(borderStyle("┌────────────────────────────────┐"))
		fmt.Printf("%s %s%-10s%s %s\n", borderStyle("│"), headerStyle(""), currUpper, headerStyle(""), borderStyle(" │"))
		fmt.Printf("%s Converted: %s%-15s%s %s\n", borderStyle("│"), valueStyle(""), formatNumber(converted), valueStyle(""), borderStyle(" │"))
		fmt.Printf("%s Rate: %s1 %s = %.4f %s%s %s\n", borderStyle("│"), rateStyle(""), strings.ToUpper(base), rate, currUpper, strings.Repeat(" ", 11-len(currUpper)), borderStyle(" │"))
		fmt.Println(borderStyle("└────────────────────────────────┘"))
		fmt.Println()
	}
}

func printList(targets []string, rates map[string]float64, amount float64, base string) {
	title := fmt.Sprintf("Currency Exchange converting %.2f %s", amount, strings.ToUpper(base))
	fmt.Println()
	fmt.Println(titleStyle(title))
	fmt.Println()

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]
		if !ok {
			fmt.Printf("  %s%-6s%s → %s\n", headerStyle(""), currUpper, headerStyle(""), errorStyle("N/A"))
			continue
		}

		converted := amount * rate
		fmt.Printf("  %s%-6s%s → %s%-15s%s  %s\n",
			headerStyle(""), currUpper, headerStyle(""),
			valueStyle(""), formatNumber(converted), valueStyle(""),
			rateStyle(fmt.Sprintf("(1 %s = %.4f %s)", strings.ToUpper(base), rate, currUpper)))
	}

	fmt.Println()
}

func printMinimal(targets []string, rates map[string]float64, amount float64, base string) {
	titleLine := fmt.Sprintf("%.2f %s →", amount, strings.ToUpper(base))
	fmt.Println()
	fmt.Println(titleStyle(titleLine))

	var results []string
	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]
		if !ok {
			results = append(results, errorStyle(currUpper+"=N/A"))
			continue
		}

		converted := amount * rate
		results = append(results, valueStyle(currUpper)+"="+formatNumber(converted))
	}

	fmt.Printf("  %s\n\n", strings.Join(results, "  |  "))
}

func printResults(targets []string, rates map[string]float64, amount float64, base string, format string) {
	switch format {
	case "cards":
		printCards(targets, rates, amount, base)
	case "list":
		printList(targets, rates, amount, base)
	case "minimal":
		printMinimal(targets, rates, amount, base)
	default:
		printTable(targets, rates, amount, base)
	}
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "-help" || os.Args[1] == "--help") {
		printHelp()
		os.Exit(0)
	}

	flag.StringVar(&themeName, "theme", "ocean", "Color theme: ocean, sunset, forest, neon")
	flag.StringVar(&displayFormat, "format", "table", "Display format: table, cards, list, minimal")
	flag.Parse()

	initStyles(themeName)

	args := flag.Args()
	if len(args) < 4 {
		printHelp()
		os.Exit(1)
	}

	amount, err := parseAmount(args[0])
	if err != nil {
		fmt.Printf("Invalid amount: %s\n", args[0])
		os.Exit(1)
	}

	base := args[1]

	if args[2] != "to" {
		fmt.Println("Error: expected 'to' as third argument")
		fmt.Println("Usage: cex <amount> <from_currency> to <to_currency1> <to_currency2> ...")
		os.Exit(1)
	}

	var targets []string
	for i := 3; i < len(args); i++ {
		targets = append(targets, args[i])
	}

	rates, err := getRates(base)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	printResults(targets, rates, amount, base, displayFormat)
}

func printHelp() {
	fmt.Println("Currency Exchange CLI - Convert between currencies")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  cex [options] <amount> <from_currency> to <to_currency1> [to_currency2] ...")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  cex 100 czk to pln eur usd")
	fmt.Println("  cex 700k chf to czk")
	fmt.Println("  cex -theme sunset 50 usd to gbp jpy cad")
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
	fmt.Println()
	fmt.Println("AMOUNT FORMAT:")
	fmt.Println("  Supports 'k' suffix for thousands (e.g., 700k = 700,000)")
}

var (
	themeName     string
	displayFormat string
)
