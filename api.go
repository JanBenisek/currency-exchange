package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// APIProvider represents a currency rate API provider
type APIProvider string

const (
	// APIExchangeRates is the default Exchange Rates API (free, no key required)
	APIExchangeRates APIProvider = "exchangerates"
	// APIFixer is the Fixer.io API (requires API key)
	APIFixer APIProvider = "fixer"
	// APICurrencyLayer is the CurrencyLayer API (requires API key)
	APICurrencyLayer APIProvider = "currencylayer"
	// APIOpenExchangeRates is the Open Exchange Rates API (requires API key)
	APIOpenExchangeRates APIProvider = "openexchangerates"
)

// RateProvider defines the interface for fetching exchange rates
type RateProvider interface {
	// FetchRates fetches rates for the given base currency
	FetchRates(base string) (map[string]float64, error)
	// Name returns the provider name
	Name() string
}

// ExchangeRatesProvider implements the Exchange Rates API (v4 or v6)
type ExchangeRatesProvider struct {
	client *http.Client
	apiKey string
	useV6  bool
}

// FixerProvider implements the Fixer.io API
type FixerProvider struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

// CurrencyLayerProvider implements the CurrencyLayer API
type CurrencyLayerProvider struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

// OpenExchangeRatesProvider implements the Open Exchange Rates API
type OpenExchangeRatesProvider struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

// Response structures for different APIs

// ExchangeRatesResponse represents Exchange Rates API response
type ExchangeRatesResponse struct {
	Rates map[string]float64 `json:"rates"`
}

// FixerResponse represents Fixer.io API response
type FixerResponse struct {
	Success bool               `json:"success"`
	Rates   map[string]float64 `json:"rates"`
	Error   *FixerError        `json:"error,omitempty"`
}

type FixerError struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}

// CurrencyLayerResponse represents CurrencyLayer API response
type CurrencyLayerResponse struct {
	Success bool                `json:"success"`
	Quotes  map[string]float64  `json:"quotes"`
	Error   *CurrencyLayerError `json:"error,omitempty"`
}

type CurrencyLayerError struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}

// OpenExchangeRatesResponse represents Open Exchange Rates API response
type OpenExchangeRatesResponse struct {
	Rates map[string]float64 `json:"rates"`
}

// NewProvider creates a new rate provider based on the selected API
func NewProvider(provider APIProvider) (RateProvider, error) {
	client := &http.Client{
		Timeout: httpTimeout,
	}

	switch provider {
	case APIExchangeRates:
		apiKey := os.Getenv("EXCHANGERATE_API_KEY")
		useV6 := apiKey != ""
		return &ExchangeRatesProvider{
			client: client,
			apiKey: apiKey,
			useV6:  useV6,
		}, nil

	case APIFixer:
		apiKey := os.Getenv("FIXER_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("FIXER_API_KEY environment variable not set")
		}
		return &FixerProvider{
			client:  client,
			apiKey:  apiKey,
			baseURL: "https://api.fixer.io/latest",
		}, nil

	case APICurrencyLayer:
		apiKey := os.Getenv("CURRENCYLAYER_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("CURRENCYLAYER_API_KEY environment variable not set")
		}
		return &CurrencyLayerProvider{
			client:  client,
			apiKey:  apiKey,
			baseURL: "https://api.currencylayer.com/api/live",
		}, nil

	case APIOpenExchangeRates:
		apiKey := os.Getenv("OPENEXCHANGERATES_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OPENEXCHANGERATES_API_KEY environment variable not set")
		}
		return &OpenExchangeRatesProvider{
			client:  client,
			apiKey:  apiKey,
			baseURL: "https://openexchangerates.org/api/latest.json",
		}, nil

	default:
		return nil, fmt.Errorf("unknown API provider: %s", provider)
	}
}

// FetchRates implements RateProvider for ExchangeRatesProvider
func (p *ExchangeRatesProvider) FetchRates(base string) (map[string]float64, error) {
	var url string
	baseUpper := strings.ToUpper(base)

	if p.useV6 {
		// v6 requires API key: https://v6.exchangerate-api.com/v6/YOUR-API-KEY/latest/USD
		url = fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", p.apiKey, baseUpper)
	} else {
		// v4 is free, no API key required: https://api.exchangerate-api.com/v4/latest/USD
		url = fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/%s", baseUpper)
	}

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("currency '%s' not found", baseUpper)
		}
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rateResp ExchangeRatesResponse
	if err := json.Unmarshal(body, &rateResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return rateResp.Rates, nil
}

func (p *ExchangeRatesProvider) Name() string {
	if p.useV6 {
		return "Exchange Rates API (v6)"
	}
	return "Exchange Rates API (v4)"
}

// FetchRates implements RateProvider for FixerProvider
func (p *FixerProvider) FetchRates(base string) (map[string]float64, error) {
	url := fmt.Sprintf("%s?access_key=%s&base=%s", p.baseURL, p.apiKey, strings.ToUpper(base))

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rateResp FixerResponse
	if err := json.Unmarshal(body, &rateResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !rateResp.Success {
		if rateResp.Error != nil {
			return nil, fmt.Errorf("API error: %s", rateResp.Error.Info)
		}
		return nil, fmt.Errorf("API request failed")
	}

	return rateResp.Rates, nil
}

func (p *FixerProvider) Name() string {
	return "Fixer.io"
}

// FetchRates implements RateProvider for CurrencyLayerProvider
func (p *CurrencyLayerProvider) FetchRates(base string) (map[string]float64, error) {
	baseUpper := strings.ToUpper(base)
	url := fmt.Sprintf("%s?access_key=%s&source=%s", p.baseURL, p.apiKey, baseUpper)

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rateResp CurrencyLayerResponse
	if err := json.Unmarshal(body, &rateResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !rateResp.Success {
		if rateResp.Error != nil {
			return nil, fmt.Errorf("API error: %s", rateResp.Error.Info)
		}
		return nil, fmt.Errorf("API request failed")
	}

	// CurrencyLayer returns rates as concatenated currency codes (e.g., "USDEUR": 0.85)
	// We need to convert this to our standard format
	rates := make(map[string]float64)
	for quote, rate := range rateResp.Quotes {
		if len(quote) == 6 {
			// Extract target currency (last 3 characters)
			target := quote[3:6]
			rates[target] = rate
		}
	}

	// Add base currency as 1.0
	rates[baseUpper] = 1.0

	return rates, nil
}

func (p *CurrencyLayerProvider) Name() string {
	return "CurrencyLayer"
}

// FetchRates implements RateProvider for OpenExchangeRatesProvider
func (p *OpenExchangeRatesProvider) FetchRates(base string) (map[string]float64, error) {
	url := fmt.Sprintf("%s?app_id=%s&base=%s", p.baseURL, p.apiKey, strings.ToUpper(base))

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rateResp OpenExchangeRatesResponse
	if err := json.Unmarshal(body, &rateResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return rateResp.Rates, nil
}

func (p *OpenExchangeRatesProvider) Name() string {
	return "Open Exchange Rates"
}

// getRates gets the current rates using the specified provider
// This preserves your existing current-conversion behavior.
func getRates(base string, provider RateProvider) (map[string]float64, error) {
	return provider.FetchRates(base)
}
