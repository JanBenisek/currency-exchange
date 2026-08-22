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

	"github.com/charmbracelet/lipgloss"
)

const apiURL = "https://api.exchangerate-api.com/v4/latest/"

// RateResponse represents the API response
type RateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

// Styles using Lip Gloss
var (
	borderStyle   lipgloss.Style
	headerStyle   lipgloss.Style
	valueStyle    lipgloss.Style
	rateStyle     lipgloss.Style
	errorStyle    lipgloss.Style
	titleStyle    lipgloss.Style
	resetStyle    lipgloss.Style
)

func initStyles(themeName string) {
	// Color palettes for different themes
	var (
		borderColor lipgloss.Color
		headerColor lipgloss.Color
		valueColor  lipgloss.Color
		rateColor   lipgloss.Color
		errorColor  lipgloss.Color
		titleColor  lipgloss.Color
	)

	switch themeName {
	case "ocean":
		borderColor = lipgloss.Color("36")  // Cyan
		headerColor = lipgloss.Color("96")  // Bright Cyan
		valueColor = lipgloss.Color("34")   // Blue
		rateColor = lipgloss.Color("34")    // Blue
		errorColor = lipgloss.Color("196") // Red
		titleColor = lipgloss.Color("96")   // Bright Cyan
	case "sunset":
		borderColor = lipgloss.Color("208") // Orange
		headerColor = lipgloss.Color("226") // Yellow
		valueColor = lipgloss.Color("208")  // Orange
		rateColor = lipgloss.Color("214")  // Yellow
		errorColor = lipgloss.Color("196") // Red
		titleColor = lipgloss.Color("226") // Yellow
	case "forest":
		borderColor = lipgloss.Color("154") // Lime
		headerColor = lipgloss.Color("46")  // Green
		valueColor = lipgloss.Color("154")  // Lime
		rateColor = lipgloss.Color("34")    // Green
		errorColor = lipgloss.Color("130") // Brown
		titleColor = lipgloss.Color("46")   // Green
	case "neon":
		borderColor = lipgloss.Color("201") // Magenta
		headerColor = lipgloss.Color("201") // Magenta
		valueColor = lipgloss.Color("46")   // Green
		rateColor = lipgloss.Color("27")    // Blue
		errorColor = lipgloss.Color("226") // Yellow
		titleColor = lipgloss.Color("201")  // Magenta
	default:
		borderColor = lipgloss.Color("36")
		headerColor = lipgloss.Color("96")
		valueColor = lipgloss.Color("34")
		rateColor = lipgloss.Color("34")
		errorColor = lipgloss.Color("196")
		titleColor = lipgloss.Color("96")
	}

	// Base styles
	borderStyle = lipgloss.NewStyle().
		Foreground(borderColor).
		Bold(true)

	headerStyle = lipgloss.NewStyle().
		Foreground(headerColor).
		Bold(true)

	valueStyle = lipgloss.NewStyle().
		Foreground(valueColor).
		Bold(true)

	rateStyle = lipgloss.NewStyle().
		Foreground(rateColor)

	errorStyle = lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true)

	titleStyle = lipgloss.NewStyle().
		Foreground(titleColor).
		Bold(true)

	resetStyle = lipgloss.NewStyle()
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
	str := fmt.Sprintf("%.2f", f)
	parts := strings.Split(str, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = "." + parts[1]
	}

	// Add thousand separators
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
	// Title
	title := fmt.Sprintf("Currency Exchange converting %.2f %s", amount, strings.ToUpper(base))
	fmt.Println()
	fmt.Println(titleStyle.Render(title))
	fmt.Println()

	// Table dimensions
	borderWidth := 71

	// Render border characters
	leftBorder := borderStyle.Render("┌")
	rightBorder := borderStyle.Render("┐")
	bottomLeft := borderStyle.Render("└")
	bottomRight := borderStyle.Render("┘")
	borderLine := borderStyle.Render("─")
	pipeChar := borderStyle.Render("│")
	middleLeft := borderStyle.Render("├")
	middleRight := borderStyle.Render("┤")

	// Top border
	fmt.Println(leftBorder + strings.Repeat(borderLine, borderWidth-2) + rightBorder)

	// Header
	headerCurrency := headerStyle.Render("Currency")
	headerConverted := headerStyle.Render("Converted")
	headerRate := headerStyle.Render(fmt.Sprintf("Rate (1 %s = X)", strings.ToUpper(base)))

	header := fmt.Sprintf("%s %-9s │ %-14s │ %-31s%s", pipeChar, headerCurrency, headerConverted, headerRate, pipeChar)
	fmt.Println(header)

	// Middle border
	fmt.Println(middleLeft + strings.Repeat(borderLine, borderWidth-2) + middleRight)

	// Data rows
	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]
		if !ok {
			row := fmt.Sprintf("%s %-9s │ %-14s │ %-31s%s", pipeChar, currUpper, errorStyle.Render("N/A"), "Currency not found", pipeChar)
			fmt.Println(row)
			continue
		}

		converted := amount * rate
		rateText := rateStyle.Render(fmt.Sprintf("1 %s = %.4f %s", strings.ToUpper(base), rate, currUpper))
		convertedText := valueStyle.Render(formatNumber(converted))

		row := fmt.Sprintf("%s %-9s │ %-14s │ %-31s%s", pipeChar, currUpper, convertedText, rateText, pipeChar)
		fmt.Println(row)
	}

	// Bottom border
	fmt.Println(bottomLeft + strings.Repeat(borderLine, borderWidth-2) + bottomRight)
	fmt.Println()
}

func printCards(targets []string, rates map[string]float64, amount float64, base string) {
	// Title
	title := fmt.Sprintf("Currency Exchange converting %.2f %s", amount, strings.ToUpper(base))
	fmt.Println()
	fmt.Println(titleStyle.Render(title))
	fmt.Println()

	cardWidth := 34
	topBorder := borderStyle.Render("┌" + strings.Repeat("─", cardWidth-2) + "┐")
	bottomBorder := borderStyle.Render("└" + strings.Repeat("─", cardWidth-2) + "┘")

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]
		if !ok {
			fmt.Println(topBorder)
			errorText := errorStyle.Render(fmt.Sprintf("  Not found: %s", currUpper))
			padding := cardWidth - 4 - len(currUpper) - 10
			fmt.Printf("%s %s%s%s\n", borderStyle.Render("│"), errorText, strings.Repeat(" ", padding), borderStyle.Render("│"))
			fmt.Println(bottomBorder)
			fmt.Println()
			continue
		}

		converted := amount * rate

		fmt.Println(topBorder)

		// Header line
		headerText := headerStyle.Render("  " + currUpper)
		padding := cardWidth - 4 - len(currUpper)
		fmt.Printf("%s %s%s%s\n", borderStyle.Render("│"), headerText, strings.Repeat(" ", padding), borderStyle.Render("│"))

		// Converted line
		convertedText := valueStyle.Render(formatNumber(converted))
		padding = cardWidth - 14 - len(formatNumber(converted))
		fmt.Printf("%s Converted:  %s%s%s\n", borderStyle.Render("│"), convertedText, strings.Repeat(" ", padding), borderStyle.Render("│"))

		// Rate line
		rateText := rateStyle.Render(fmt.Sprintf("1 %s = %.4f %s", strings.ToUpper(base), rate, currUpper))
		padding = cardWidth - 10 - len(fmt.Sprintf("1 %s = %.4f %s", strings.ToUpper(base), rate, currUpper))
		fmt.Printf("%s Rate:      %s%s%s\n", borderStyle.Render("│"), rateText, strings.Repeat(" ", padding), borderStyle.Render("│"))

		fmt.Println(bottomBorder)
		fmt.Println()
	}
}

func printList(targets []string, rates map[string]float64, amount float64, base string) {
	// Title
	title := fmt.Sprintf("Currency Exchange converting %.2f %s", amount, strings.ToUpper(base))
	fmt.Println()
	fmt.Println(titleStyle.Render(title))
	fmt.Println()

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]
		if !ok {
			line := fmt.Sprintf("%s %s%-6s%s → %s",
				borderStyle.Render("│"),
				headerStyle,
				currUpper,
				headerStyle,
				errorStyle.Render("N/A"))
			fmt.Println(line)
			continue
		}

		converted := amount * rate
		rateText := rateStyle.Render(fmt.Sprintf("(1 %s = %.4f %s)", strings.ToUpper(base), rate, currUpper))

		line := fmt.Sprintf("%s %s%-6s%s → %s%-15s%s  %s",
			borderStyle.Render("│"),
			headerStyle,
			currUpper,
			headerStyle,
			valueStyle,
			formatNumber(converted),
			valueStyle,
			rateText)
		fmt.Println(line)
	}

	fmt.Println()
}

func printMinimal(targets []string, rates map[string]float64, amount float64, base string) {
	// Title line
	titleLine := fmt.Sprintf("%.2f %s →", amount, strings.ToUpper(base))
	fmt.Println()
	fmt.Println(titleStyle.Render(titleLine))

	var results []string
	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)
		rate, ok := rates[currUpper]
		if !ok {
			results = append(results, errorStyle.Render(currUpper+"=N/A"))
			continue
		}

		converted := amount * rate
		result := fmt.Sprintf("%s=%s", valueStyle.Render(currUpper), formatNumber(converted))
		results = append(results, result)
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
	// Check for help
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "-help" || os.Args[1] == "--help") {
		printHelp()
		os.Exit(0)
	}

	// Parse flags
	flag.StringVar(&themeName, "theme", "ocean", "Color theme: ocean, sunset, forest, neon")
	flag.StringVar(&displayFormat, "format", "table", "Display format: table, cards, list, minimal")
	flag.Parse()

	// Initialize styles with selected theme
	initStyles(themeName)

	// Get remaining args after flags
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

func init() {
	// Default values are set in flag.StringVar
}
