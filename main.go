package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	apiURL = "https://api.exchangerate-api.com/v4/latest/"

	// Free historical API — no key required.
	frankfurterURL = "https://api.frankfurter.dev/v2/rates"

	historyDays = 90

	// Width of each compact history chart.
	chartWidth  = 20
	chartHeight = 5
	chartGap    = 2 

	httpTimeout = 15 * time.Second
)

// RateResponse represents the current ExchangeRate-API response.
type RateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

// FrankfurterRate represents one historical rate row.
type FrankfurterRate struct {
	Date  string  `json:"date"`
	Base  string  `json:"base"`
	Quote string  `json:"quote"`
	Rate  float64 `json:"rate"`
}

// HistoryPoint is one point used by the terminal chart.
type HistoryPoint struct {
	Date time.Time
	Rate float64
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
	themeName     string
	displayFormat string
	showHistory   bool
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

// getRates gets the current rates.
// This preserves your existing current-conversion behavior.
func getRates(base string) (map[string]float64, error) {
	client := &http.Client{
		Timeout: httpTimeout,
	}

	resp, err := client.Get(apiURL + strings.ToUpper(base))
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

// getHistoricalRates gets the entire 3-month series in ONE request.
//
// Frankfurter supports:
//   ?from=YYYY-MM-DD&to=YYYY-MM-DD&base=CZK&quotes=EUR,GBP,USD
//
// No API key is required.
func getHistoricalRates(
	base string,
	targets []string,
) (map[string][]HistoryPoint, error) {
	if len(targets) == 0 {
		return nil, fmt.Errorf("no target currencies")
	}

	now := time.Now()

	endDate := now
	startDate := now.AddDate(0, 0, -historyDays)

	quotes := make([]string, 0, len(targets))

	for _, target := range targets {
		quotes = append(
			quotes,
			strings.ToUpper(target),
		)
	}

	params := url.Values{}
	params.Set("from", startDate.Format("2006-01-02"))
	params.Set("to", endDate.Format("2006-01-02"))
	params.Set("base", strings.ToUpper(base))
	params.Set("quotes", strings.Join(quotes, ","))

	requestURL := frankfurterURL + "?" + params.Encode()

	client := &http.Client{
		Timeout: httpTimeout,
	}

	resp, err := client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to fetch historical rates: %w",
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		if len(body) > 0 {
			return nil, fmt.Errorf(
				"Frankfurter returned HTTP %d: %s",
				resp.StatusCode,
				strings.TrimSpace(string(body)),
			)
		}

		return nil, fmt.Errorf(
			"Frankfurter returned HTTP %d",
			resp.StatusCode,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read historical response: %w",
			err,
		)
	}

	var rows []FrankfurterRate

	if err := json.Unmarshal(body, &rows); err != nil {
		return nil, fmt.Errorf(
			"failed to parse historical response: %w",
			err,
		)
	}

	history := make(map[string][]HistoryPoint)

	for _, target := range targets {
		history[strings.ToUpper(target)] = []HistoryPoint{}
	}

	for _, row := range rows {
		date, err := time.Parse(
			"2006-01-02",
			row.Date,
		)

		if err != nil {
			continue
		}

		quote := strings.ToUpper(row.Quote)

		if _, exists := history[quote]; !exists {
			continue
		}

		history[quote] = append(
			history[quote],
			HistoryPoint{
				Date: date,
				Rate: row.Rate,
			},
		)
	}

	return history, nil
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

// buildChart generates a compact Unicode line chart.

// buildChart creates a compact dot-based chart.
//
// Example:
//
// ┌──────────────────┐
// │ EUR  3M          │
// │       ..         │
// │    ...  ...      │
// │  ..        ..    │
// │..            ... │
// │ 06/2026  07/2026 │
// └──────────────────┘
//
// The dots are deliberately not connected. This keeps the tiny
// terminal charts much easier to read.
func buildChart(history []HistoryPoint) []string {
	lines := make([]string, chartHeight)

	for i := range lines {
		lines[i] = strings.Repeat(" ", chartWidth)
	}

	if len(history) < 2 {
		lines[chartHeight/2] = "no data"
		return lines
	}

	// Find min/max rate.
	minRate := history[0].Rate
	maxRate := history[0].Rate

	for _, point := range history {
		if point.Rate < minRate {
			minRate = point.Rate
		}

		if point.Rate > maxRate {
			maxRate = point.Rate
		}
	}

	// Avoid division by zero when all values are identical.
	if maxRate == minRate {
		row := []rune(
			strings.Repeat(".", chartWidth),
		)

		lines[chartHeight/2] = string(row)

		return lines
	}

	// Create a blank plotting canvas.
	canvas := make([][]rune, chartHeight)

	for y := 0; y < chartHeight; y++ {
		canvas[y] = make([]rune, chartWidth)

		for x := 0; x < chartWidth; x++ {
			canvas[y][x] = ' '
		}
	}

	// Plot observations.
	//
	// We use one x-coordinate for each horizontal terminal column.
	// Multiple daily observations that fall into the same column are
	// represented by the average value for that column.
	buckets := make([][]float64, chartWidth)

	for _, point := range history {
		first := history[0].Date
		last := history[len(history)-1].Date

		totalSeconds := last.Sub(first).Seconds()

		if totalSeconds <= 0 {
			continue
		}

		position :=
			point.Date.Sub(first).Seconds() /
				totalSeconds

		x := int(
			position *
				float64(chartWidth-1),
		)

		if x < 0 {
			x = 0
		}

		if x >= chartWidth {
			x = chartWidth - 1
		}

		buckets[x] = append(
			buckets[x],
			point.Rate,
		)
	}

	// Plot each horizontal bucket.
	for x, bucket := range buckets {
		if len(bucket) == 0 {
			continue
		}

		var sum float64

		for _, value := range bucket {
			sum += value
		}

		value := sum / float64(len(bucket))

		normalized :=
			(value - minRate) /
				(maxRate - minRate)

		y := chartHeight -
			1 -
			int(
				normalized*
					float64(chartHeight-1),
			)

		if y < 0 {
			y = 0
		}

		if y >= chartHeight {
			y = chartHeight - 1
		}

		canvas[y][x] = '.'
	}

	// Guarantee that the first and last observations are visible.
	plotPoint := func(point HistoryPoint) {
		first := history[0].Date
		last := history[len(history)-1].Date

		totalSeconds := last.Sub(first).Seconds()

		if totalSeconds <= 0 {
			return
		}

		position :=
			point.Date.Sub(first).Seconds() /
				totalSeconds

		x := int(
			position *
				float64(chartWidth-1),
		)

		if x < 0 {
			x = 0
		}

		if x >= chartWidth {
			x = chartWidth - 1
		}

		normalized :=
			(point.Rate - minRate) /
				(maxRate - minRate)

		y := chartHeight -
			1 -
			int(
				normalized*
					float64(chartHeight-1),
			)

		if y < 0 {
			y = 0
		}

		if y >= chartHeight {
			y = chartHeight - 1
		}

		canvas[y][x] = '.'
	}

	plotPoint(history[0])
	plotPoint(history[len(history)-1])

	for y := 0; y < chartHeight; y++ {
		lines[y] = string(canvas[y])
	}

	return lines
}

func buildMonthAxis(history []HistoryPoint) string {
	if len(history) < 2 {
		return fmt.Sprintf(
			"%-*s",
			chartWidth,
			"",
		)
	}

	first := history[0].Date
	last := history[len(history)-1].Date

	type monthTick struct {
		position int
		label    string
	}

	var ticks []monthTick

	// Start at the first day of the first month.
	current := time.Date(
		first.Year(),
		first.Month(),
		1,
		0,
		0,
		0,
		0,
		first.Location(),
	)

	for !current.After(last) {
		totalSeconds := last.Sub(first).Seconds()

		position := 0

		if totalSeconds > 0 {
			position = int(
				current.Sub(first).Seconds() /
					totalSeconds *
					float64(chartWidth-1),
			)
		}

		if position < 0 {
			position = 0
		}

		if position >= chartWidth {
			position = chartWidth - 1
		}

		ticks = append(
			ticks,
			monthTick{
				position: position,
				label: current.Format("01/2006"),
			},
		)

		current = current.AddDate(0, 1, 0)
	}

	// Build the tick row.
	row := make([]rune, chartWidth)

	for i := range row {
		row[i] = ' '
	}

	// Put a vertical tick at every month boundary.
	for _, tick := range ticks {
		if tick.position >= 0 &&
			tick.position < chartWidth {
			row[tick.position] = '|'
		}
	}

	// We have limited room. Put month labels on a separate row
	// and make them fit without overlapping where possible.
	labelRow := make([]rune, chartWidth)

	for i := range labelRow {
		labelRow[i] = ' '
	}

	for _, tick := range ticks {
		label := tick.label

		// Keep MM/YYYY intact. If there isn't enough room for the
		// whole label, skip that tick rather than producing garbage.
		labelWidth := len([]rune(label))

		start := tick.position - labelWidth/2

		if start < 0 {
			start = 0
		}

		if start+labelWidth > chartWidth {
			start = chartWidth - labelWidth
		}

		if start < 0 {
			continue
		}

		// Don't overwrite an already-present label.
		canPlace := true

		for i := 0; i < labelWidth; i++ {
			if labelRow[start+i] != ' ' {
				canPlace = false
				break
			}
		}

		if !canPlace {
			continue
		}

		for i, char := range []rune(label) {
			labelRow[start+i] = char
		}
	}

	return string(labelRow)
}


// chartPanel creates one compact history column.
func chartPanel(
	currency string,
	history []HistoryPoint,
) []string {
	panel := make([]string, 0)

	// Header.
	header := fmt.Sprintf(
		"%-18s",
		currency+"  3M",
	)

	panel = append(
		panel,
		headerStyle(header),
	)

	// Chart itself.
	chart := buildChart(history)

	for _, line := range chart {
		panel = append(
			panel,
			rateStyle(
				fmt.Sprintf(
					"%-18s",
					line,
				),
			),
		)
	}

	// Month tick row.
	panel = append(
		panel,
		rateStyle(
			buildMonthAxis(history),
		),
	)

	return panel
}

// printHistoryRow prints:
//
// [ TABLE ] [ EUR ] [ GBP ] [ USD ]
//
// Everything is kept in ONE horizontal row.
func printHistoryRow(
	targets []string,
	rates map[string]float64,
	amount float64,
	base string,
	histories map[string][]HistoryPoint,
) {
	const (
		tableWidth = 52
		panelWidth = chartWidth + 2
		gap        = "  "
	)

	title := fmt.Sprintf(
		"Currency Exchange converting %.2f %s",
		amount,
		strings.ToUpper(base),
	)

	fmt.Println()
	fmt.Println(titleStyle(title))
	fmt.Println()

	// ------------------------------------------------------------
	// TABLE
	// ------------------------------------------------------------

	tableLines := make([]string, 0)

	tableLines = append(
		tableLines,
		borderStyle(
			"┌───────────────────────────────┬────────────────────┐",
		),
	)

	col1Header := headerStyle(
		fmt.Sprintf(
			" %-30s",
			"Converted Value",
		),
	)

	col2Header := headerStyle(
		fmt.Sprintf(
			" %-19s",
			fmt.Sprintf(
				"Rate (1 %s = X)",
				strings.ToUpper(base),
			),
		),
	)

	tableLines = append(
		tableLines,
		borderStyle("│")+
			col1Header+
			borderStyle("│")+
			col2Header+
			borderStyle("│"),
	)

	tableLines = append(
		tableLines,
		borderStyle(
			"├───────────────────────────────┼────────────────────┤",
		),
	)

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)

		rate, ok := rates[currUpper]

		if !ok {
			col1 := errorStyle(
				fmt.Sprintf(
					" %-30s",
					"N/A",
				),
			)

			col2 := errorStyle(
				fmt.Sprintf(
					" %-19s",
					"Currency not found",
				),
			)

			tableLines = append(
				tableLines,
				borderStyle("│")+
					col1+
					borderStyle("│")+
					col2+
					borderStyle("│"),
			)

			continue
		}

		converted := amount * rate

		valueText := fmt.Sprintf(
			" %s %s",
			formatNumber(converted),
			currUpper,
		)

		rateText := fmt.Sprintf(
			" 1 %s = %.4f %s",
			strings.ToUpper(base),
			rate,
			currUpper,
		)

		col1 := valueStyle(
			fmt.Sprintf(
				"%-31s",
				valueText,
			),
		)

		col2 := rateStyle(
			fmt.Sprintf(
				"%-20s",
				rateText,
			),
		)

		tableLines = append(
			tableLines,
			borderStyle("│")+
				col1+
				borderStyle("│")+
				col2+
				borderStyle("│"),
		)
	}

	tableLines = append(
		tableLines,
		borderStyle(
			"└───────────────────────────────┴────────────────────┘",
		),
	)

	// ------------------------------------------------------------
	// THREE HISTORY COLUMNS
	// ------------------------------------------------------------

	historyTargets := make([]string, 0, 3)

	for _, target := range targets {
		if len(historyTargets) >= 3 {
			break
		}

		historyTargets = append(
			historyTargets,
			strings.ToUpper(target),
		)
	}

	panels := make([][]string, 0, 3)

	for _, currency := range historyTargets {
		panels = append(
			panels,
			chartPanel(
				currency,
				histories[currency],
			),
		)
	}

	panelHeight := chartHeight + 3

	totalRows := len(tableLines)

	if panelHeight > totalRows {
		totalRows = panelHeight
	}

	// ------------------------------------------------------------
	// ONE HORIZONTAL ROW
	// ------------------------------------------------------------

	for row := 0; row < totalRows; row++ {
		// Table.
		if row < len(tableLines) {
			fmt.Print(tableLines[row])
		} else {
			fmt.Print(
				strings.Repeat(
					" ",
					tableWidth,
				),
			)
		}

		fmt.Print(gap)

		// EUR / GBP / USD panels.
		for i, panel := range panels {
			if i > 0 {
				fmt.Print(gap)
			}

			switch {
			case row == 0:
				fmt.Print(
					borderStyle(
						"┌"+
							strings.Repeat(
								"─",
								chartWidth,
							)+
							"┐",
					),
				)

			case row == panelHeight-1:
				fmt.Print(
					borderStyle(
						"└"+
							strings.Repeat(
								"─",
								chartWidth,
							)+
							"┘",
					),
				)

			default:
				contentRow := row - 1

				content := ""

				if contentRow >= 0 &&
					contentRow < len(panel) {
					content = panel[contentRow]
				}

				content = padANSI(
					content,
					chartWidth,
				)

				fmt.Print(
					borderStyle("│")+
						content+
						borderStyle("│"),
				)
			}
		}

		fmt.Println()
	}

	fmt.Println()
}

// printTable is the normal output.
func printTable(
	targets []string,
	rates map[string]float64,
	amount float64,
	base string,
) {
	title := fmt.Sprintf(
		"Currency Exchange converting %.2f %s",
		amount,
		strings.ToUpper(base),
	)

	fmt.Println()
	fmt.Println(titleStyle(title))
	fmt.Println()

	topBorder := borderStyle(
		"┌───────────────────────────────┬────────────────────┐",
	)

	headerBorder := borderStyle(
		"├───────────────────────────────┼────────────────────┤",
	)

	bottomBorder := borderStyle(
		"└───────────────────────────────┴────────────────────┘",
	)

	leftPipe := borderStyle("│")
	rightPipe := borderStyle("│")

	fmt.Println(topBorder)

	col1Header := headerStyle(
		fmt.Sprintf(
			" %-30s",
			"Converted Value",
		),
	)

	col2Header := headerStyle(
		fmt.Sprintf(
			" %-19s",
			fmt.Sprintf(
				"Rate (1 %s = X)",
				strings.ToUpper(base),
			),
		),
	)

	fmt.Printf(
		"%s%s%s%s%s\n",
		leftPipe,
		col1Header,
		leftPipe,
		col2Header,
		rightPipe,
	)

	fmt.Println(headerBorder)

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)

		rate, ok := rates[currUpper]

		if !ok {
			col1 := errorStyle(
				fmt.Sprintf(
					" %-30s",
					"N/A",
				),
			)

			col2 := errorStyle(
				fmt.Sprintf(
					" %-19s",
					"Currency not found",
				),
			)

			fmt.Printf(
				"%s%s%s%s%s\n",
				leftPipe,
				col1,
				leftPipe,
				col2,
				rightPipe,
			)

			continue
		}

		converted := amount * rate

		valueWithCurrency := fmt.Sprintf(
			"%s %s",
			formatNumber(converted),
			currUpper,
		)

		rateStr := fmt.Sprintf(
			"1 %s = %.4f %s",
			strings.ToUpper(base),
			rate,
			currUpper,
		)

		col1 := valueStyle(
			fmt.Sprintf(
				" %-30s",
				valueWithCurrency,
			),
		)

		col2 := rateStyle(
			fmt.Sprintf(
				" %-19s",
				rateStr,
			),
		)

		fmt.Printf(
			"%s%s%s%s%s\n",
			leftPipe,
			col1,
			leftPipe,
			col2,
			rightPipe,
		)
	}

	fmt.Println(bottomBorder)
	fmt.Println()
}

// printCards preserves the existing cards format.
func printCards(
	targets []string,
	rates map[string]float64,
	amount float64,
	base string,
) {
	title := fmt.Sprintf(
		"Currency Exchange converting %.2f %s",
		amount,
		strings.ToUpper(base),
	)

	fmt.Println()
	fmt.Println(titleStyle(title))
	fmt.Println()

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)

		rate, ok := rates[currUpper]

		if !ok {
			fmt.Println(
				borderStyle(
					"┌────────────────────────────────┐",
				),
			)

			fmt.Printf(
				"%s %s%s\n",
				borderStyle("│"),
				errorStyle(
					fmt.Sprintf(
						"Not found: %s",
						currUpper,
					),
				),
				borderStyle(" │"),
			)

			fmt.Println(
				borderStyle(
					"└────────────────────────────────┘",
				),
			)

			fmt.Println()
			continue
		}

		converted := amount * rate

		fmt.Println(
			borderStyle(
				"┌────────────────────────────────┐",
			),
		)

		fmt.Printf(
			"%s %s%-10s%s %s\n",
			borderStyle("│"),
			headerStyle(""),
			currUpper,
			headerStyle(""),
			borderStyle(" │"),
		)

		fmt.Printf(
			"%s Converted: %s%-15s%s %s\n",
			borderStyle("│"),
			valueStyle(""),
			formatNumber(converted),
			valueStyle(""),
			borderStyle(" │"),
		)

		fmt.Printf(
			"%s Rate: %s1 %s = %.4f %s%s %s\n",
			borderStyle("│"),
			rateStyle(""),
			strings.ToUpper(base),
			rate,
			currUpper,
			strings.Repeat(
				" ",
				maxInt(1, 11-len(currUpper)),
			),
			borderStyle(" │"),
		)

		fmt.Println(
			borderStyle(
				"└────────────────────────────────┘",
			),
		)

		fmt.Println()
	}
}

// printList preserves the existing list format.
func printList(
	targets []string,
	rates map[string]float64,
	amount float64,
	base string,
) {
	title := fmt.Sprintf(
		"Currency Exchange converting %.2f %s",
		amount,
		strings.ToUpper(base),
	)

	fmt.Println()
	fmt.Println(titleStyle(title))
	fmt.Println()

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)

		rate, ok := rates[currUpper]

		if !ok {
			fmt.Printf(
				"  %s%-6s%s → %s\n",
				headerStyle(""),
				currUpper,
				headerStyle(""),
				errorStyle("N/A"),
			)

			continue
		}

		converted := amount * rate

		fmt.Printf(
			"  %s%-6s%s → %s%-15s%s  %s\n",
			headerStyle(""),
			currUpper,
			headerStyle(""),
			valueStyle(""),
			formatNumber(converted),
			valueStyle(""),
			rateStyle(
				fmt.Sprintf(
					"(1 %s = %.4f %s)",
					strings.ToUpper(base),
					rate,
					currUpper,
				),
			),
		)
	}

	fmt.Println()
}

// printMinimal preserves the existing minimal format.
func printMinimal(
	targets []string,
	rates map[string]float64,
	amount float64,
	base string,
) {
	titleLine := fmt.Sprintf(
		"%.2f %s →",
		amount,
		strings.ToUpper(base),
	)

	fmt.Println()
	fmt.Println(titleStyle(titleLine))

	var results []string

	for _, curr := range targets {
		currUpper := strings.ToUpper(curr)

		rate, ok := rates[currUpper]

		if !ok {
			results = append(
				results,
				errorStyle(currUpper+"=N/A"),
			)

			continue
		}

		converted := amount * rate

		results = append(
			results,
			valueStyle(currUpper)+
				"="+
				formatNumber(converted),
		)
	}

	fmt.Printf(
		"  %s\n\n",
		strings.Join(
			results,
			"  |  ",
		),
	)
}

// printResults chooses the output format.
func printResults(
	targets []string,
	rates map[string]float64,
	amount float64,
	base string,
	format string,
	history bool,
) {
	// History is specifically designed for the table layout.
	if history && format == "table" {
		histories, err := getHistoricalRates(
			base,
			targets,
		)

		if err != nil {
			fmt.Println()
			fmt.Println(
				errorStyle(
					"Unable to load history: "+err.Error(),
				),
			)

			// Still show the normal conversion.
			printTable(
				targets,
				rates,
				amount,
				base,
			)

			return
		}

		printHistoryRow(
			targets,
			rates,
			amount,
			base,
			histories,
		)

		return
	}

	switch format {
	case "cards":
		printCards(
			targets,
			rates,
			amount,
			base,
		)

	case "list":
		printList(
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
	fmt.Println("USAGE:")
	fmt.Println("  cex [options] <amount> <from_currency> to <to_currency1> <to_currency2> ...")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  cex 100 czk to pln eur usd")
	fmt.Println("  cex 700k chf to czk")
	fmt.Println("  cex -theme sunset 50 usd to gbp jpy cad")
	fmt.Println("  cex -h 100 czk to eur gbp usd")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -h")
	fmt.Println("      Show 3-month rate history.")
	fmt.Println("      History is displayed as 3 compact charts")
	fmt.Println("      next to the conversion table.")
	fmt.Println()
	fmt.Println("  -theme string")
	fmt.Println("      Color theme (default \"ocean\")")
	fmt.Println("      Available: ocean, sunset, forest, neon")
	fmt.Println()
	fmt.Println("  -format string")
	fmt.Println("      Display format (default \"table\")")
	fmt.Println("      Available: table, cards, list, minimal")
	fmt.Println()
	fmt.Println("  -help")
	fmt.Println("      Show this help.")
	fmt.Println()
	fmt.Println("  --help")
	fmt.Println("      Show this help.")
	fmt.Println()
	fmt.Println("HISTORY:")
	fmt.Println("  Uses the free Frankfurter API.")
	fmt.Println("  No API key or account is required.")
	fmt.Println("  Historical rates cover approximately the")
	fmt.Println("  previous 3 months.")
	fmt.Println()
	fmt.Println("  History is daily reference-rate data, so")
	fmt.Println("  weekends and holidays may have no new point.")
	fmt.Println()
	fmt.Println("THEME PREVIEW:")
	fmt.Println("  Ocean:   Blues and cyans - professional look")
	fmt.Println("  Sunset:  Warm tones - yellows and oranges")
	fmt.Println("  Forest:  Greens and earthy tones")
	fmt.Println("  Neon:    Bright high-contrast colors")
	fmt.Println()
	fmt.Println("AMOUNT FORMAT:")
	fmt.Println("  Supports 'k' suffix for thousands")
	fmt.Println("  (e.g., 700k = 700,000)")
}

func main() {
	// IMPORTANT:
	//
	// We use -h for HISTORY instead of Go's default help flag.
	// Handle help aliases before flag.Parse().
	for _, arg := range os.Args[1:] {
		if arg == "-help" || arg == "--help" {
			printHelp()
			return
		}
	}

	flag.BoolVar(
		&showHistory,
		"h",
		false,
		"Show 3-month rate history",
	)

	flag.StringVar(
		&themeName,
		"theme",
		"ocean",
		"Color theme: ocean, sunset, forest, neon",
	)

	flag.StringVar(
		&displayFormat,
		"format",
		"table",
		"Display format: table, cards, list, minimal",
	)

	flag.Usage = printHelp

	if err := flag.CommandLine.Parse(
		os.Args[1:],
	); err != nil {
		os.Exit(1)
	}

	initStyles(themeName)

	args := flag.Args()

	if len(args) < 4 {
		printHelp()
		os.Exit(1)
	}

	amount, err := parseAmount(args[0])

	if err != nil {
		fmt.Printf(
			"Invalid amount: %s\n",
			args[0],
		)

		os.Exit(1)
	}

	base := args[1]

	if strings.ToLower(args[2]) != "to" {
		fmt.Println(
			"Error: expected 'to' as third argument",
		)

		fmt.Println(
			"Usage: cex <amount> <from_currency> to <to_currency1> <to_currency2> ...",
		)

		os.Exit(1)
	}

	var targets []string

	for i := 3; i < len(args); i++ {
		targets = append(
			targets,
			args[i],
		)
	}

	rates, err := getRates(base)

	if err != nil {
		fmt.Printf(
			"Error: %v\n",
			err,
		)

		os.Exit(1)
	}

	printResults(
		targets,
		rates,
		amount,
		base,
		displayFormat,
		showHistory,
	)
}