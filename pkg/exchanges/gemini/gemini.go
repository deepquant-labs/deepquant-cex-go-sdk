package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/client"
	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/errors"
	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/exchange"
	"github.com/rs/zerolog"
)

const (
	// API endpoints
	baseURLProd    = "https://api.gemini.com"
	baseURLSandbox = "https://api.sandbox.gemini.com"
	// Exchange name
	exchangeName = "gemini"
	// API response status
	errorStatus = "error"
)

// Gemini represents the Gemini exchange
type Gemini struct {
	client    *client.HTTPClient
	baseURL   string
	apiKey    string
	apiSecret string
	sandbox   bool
	userAgent string
	logger    zerolog.Logger

	// API categories
	Market *MarketAPI
	Order  *OrderAPI
	Fund   *FundAPI
}

// NewGemini creates a new Gemini exchange instance
func NewGemini(config *exchange.Config) *Gemini {
	baseURL := baseURLProd
	if config != nil && config.Testnet {
		baseURL = baseURLSandbox
	}

	timeout := 30 * time.Second
	if config != nil && config.Timeout > 0 {
		timeout = config.Timeout
	}

	g := &Gemini{
		client:    client.NewHTTPClient(timeout),
		baseURL:   baseURL,
		userAgent: "CEX-SDK/1.0",
		logger:    zerolog.Nop(), // Default no-op logger
	}

	if config != nil {
		g.apiKey = config.APIKey
		g.apiSecret = config.SecretKey
		g.sandbox = config.Testnet
		// UserAgent can be set via headers

		// Set custom logger if provided
		if config.Logger != nil {
			g.logger = *config.Logger
			g.client.SetLogger(*config.Logger)
		}
		// Set custom HTTP client if provided
		if config.HTTPClient != nil {
			g.client.SetCustomHTTPClient(config.HTTPClient)
		}
		// Set rate limits
		if config.RateLimit.Public.Requests > 0 {
			g.client.SetRateLimit(client.APITypePublic, config.RateLimit.Public.Requests, config.RateLimit.Public.Interval)
		} else {
			// Default public API rate limit: 120 requests per minute
			g.client.SetRateLimit(client.APITypePublic, 120, time.Minute)
		}
		if config.RateLimit.Private.Requests > 0 {
			g.client.SetRateLimit(client.APITypePrivate, config.RateLimit.Private.Requests, config.RateLimit.Private.Interval)
		} else {
			// Default private API rate limit: 600 requests per minute
			g.client.SetRateLimit(client.APITypePrivate, 600, time.Minute)
		}
	}

	// Set default headers
	headers := map[string]string{
		"User-Agent":   g.userAgent,
		"Content-Type": "application/json",
	}
	g.client.SetHeaders(headers)

	// Initialize API categories
	g.Market = NewMarketAPI(g)
	g.Order = NewOrderAPI(g)
	g.Fund = NewFundAPI(g)

	g.logger.Info().Str("baseURL", g.baseURL).Msg("Gemini exchange initialized")
	return g
}

// GetName returns the exchange name
func (g *Gemini) GetName() string {
	return exchangeName
}

// GetTradingPairs fetches all available trading pairs from Gemini
func (g *Gemini) GetTradingPairs(ctx context.Context) ([]exchange.TradingPair, error) {
	symbolsURL := fmt.Sprintf("%s/v1/symbols", g.baseURL)

	// Fetch symbols
	response, err := g.client.Get(ctx, symbolsURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNetworkError, "failed to fetch symbols", err)
	}

	var symbols []string
	if err := json.Unmarshal(response, &symbols); err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to parse symbols response", err)
	}

	// Get detailed symbol information
	detailsURL := fmt.Sprintf("%s/v1/symbols/details", g.baseURL)
	detailsResp, err := g.client.Get(ctx, detailsURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNetworkError, "failed to fetch symbol details", err)
	}

	var symbolDetails []Symbol
	if err := json.Unmarshal(detailsResp, &symbolDetails); err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to parse symbol details", err)
	}

	// Create a map for quick lookup
	detailsMap := make(map[string]Symbol)
	for _, detail := range symbolDetails {
		detailsMap[strings.ToLower(detail.Symbol)] = detail
	}

	// Fetch ticker data for each symbol
	pairs := make([]exchange.TradingPair, 0, len(symbols))
	for _, symbol := range symbols {
		detail, exists := detailsMap[strings.ToLower(symbol)]
		if !exists {
			// If no details available, create basic pair info
			pair := exchange.TradingPair{
				Symbol:     strings.ToUpper(symbol),
				BaseAsset:  extractBaseCurrency(symbol),
				QuoteAsset: extractQuoteCurrency(symbol),
				Status:     "TRADING",
				MinQty:     0,
				MaxQty:     0,
				StepSize:   0,
				TickSize:   0,
			}
			pairs = append(pairs, pair)
			continue
		}

		minOrderSize, _ := parseFloatFromString(detail.MinOrderSize)

		pair := exchange.TradingPair{
			Symbol:     strings.ToUpper(detail.Symbol),
			BaseAsset:  strings.ToUpper(detail.BaseCurrency),
			QuoteAsset: strings.ToUpper(detail.QuoteCurrency),
			Status:     detail.Status,
			MinQty:     minOrderSize,
			MaxQty:     0, // Gemini doesn't provide max order size in this endpoint
			StepSize:   0,
			TickSize:   detail.TickSize,
		}
		pairs = append(pairs, pair)
	}

	return pairs, nil
}

// SetRateLimit sets the rate limiting for the HTTP client
func (g *Gemini) SetRateLimit(apiType exchange.APIType, limit exchange.RateLimit) {
	g.client.SetRateLimit(client.APIType(apiType), limit.Requests, limit.Interval)
	g.logger.Info().Str("apiType", string(apiType)).Int("requests", limit.Requests).Dur("interval", limit.Interval).Msg("Rate limit updated")
}

// SetLogger sets custom logger
func (g *Gemini) SetLogger(logger zerolog.Logger) {
	g.logger = logger
	g.client.SetLogger(logger)
	g.logger.Info().Msg("Logger updated")
}

// SetHTTPClient sets custom HTTP client
func (g *Gemini) SetHTTPClient(client *http.Client) {
	g.client.SetCustomHTTPClient(client)
	g.logger.Info().Msg("Custom HTTP client set")
}

// SetHeaders sets custom headers for the HTTP client
func (g *Gemini) SetHeaders(headers map[string]string) {
	// Preserve essential headers
	if headers["User-Agent"] == "" {
		headers["User-Agent"] = g.userAgent
	}
	if headers["Content-Type"] == "" {
		headers["Content-Type"] = "application/json"
	}
	g.client.SetHeaders(headers)
}

// SetProxies sets proxy configuration for the HTTP client
func (g *Gemini) SetProxies(proxies []string) {
	g.client.SetProxies(proxies)
}

// SetAPICredentials sets the API credentials
func (g *Gemini) SetAPICredentials(apiKey, apiSecret string) {
	g.apiKey = apiKey
	g.apiSecret = apiSecret
}

// SetSandbox enables or disables sandbox mode
func (g *Gemini) SetSandbox(sandbox bool) {
	g.sandbox = sandbox
	if sandbox {
		g.baseURL = baseURLSandbox
	} else {
		g.baseURL = baseURLProd
	}
}

// ValidateConfig validates the exchange configuration
func (g *Gemini) ValidateConfig() error {
	// Basic validation
	if g.baseURL == "" {
		return errors.New(errors.ErrInvalidInput, "base URL is required")
	}

	// Validate URL format
	if !strings.HasPrefix(g.baseURL, "http://") && !strings.HasPrefix(g.baseURL, "https://") {
		return errors.New(errors.ErrInvalidInput, "invalid base URL format")
	}

	// Test connectivity
	testURL := fmt.Sprintf("%s/v1/symbols", g.baseURL)
	ctx := context.Background()
	_, err := g.client.Get(ctx, testURL)
	if err != nil {
		return errors.Wrap(errors.ErrNetworkError, "failed to connect to Gemini API", err)
	}

	return nil
}

// Helper functions

// extractBaseCurrency extracts base currency from symbol
// For Gemini, symbols are typically like "btcusd", "ethusd", etc.
func extractBaseCurrency(symbol string) string {
	symbol = strings.ToLower(symbol)

	// Common quote currencies in Gemini
	quoteCurrencies := []string{"usd", "btc", "eth", "eur", "gbp", "sgd", "gusd", "dai"}

	for _, quote := range quoteCurrencies {
		if strings.HasSuffix(symbol, quote) {
			return strings.ToUpper(symbol[:len(symbol)-len(quote)])
		}
	}

	// Default fallback - assume first 3 characters are base
	if len(symbol) >= 6 {
		return strings.ToUpper(symbol[:3])
	}

	return strings.ToUpper(symbol)
}

// extractQuoteCurrency extracts quote currency from symbol
func extractQuoteCurrency(symbol string) string {
	symbol = strings.ToLower(symbol)

	// Common quote currencies in Gemini
	quoteCurrencies := []string{"usd", "btc", "eth", "eur", "gbp", "sgd", "gusd", "dai"}

	for _, quote := range quoteCurrencies {
		if strings.HasSuffix(symbol, quote) {
			return strings.ToUpper(quote)
		}
	}

	// Default fallback - assume last 3 characters are quote
	if len(symbol) >= 6 {
		return strings.ToUpper(symbol[len(symbol)-3:])
	}

	return "USD" // Default to USD
}
