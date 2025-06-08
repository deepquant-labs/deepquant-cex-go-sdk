package exchange

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// APIType represents the type of API endpoint
type APIType string

const (
	APITypePublic  APIType = "public"
	APITypePrivate APIType = "private"
)

// Exchange defines the interface for cryptocurrency exchanges
type Exchange interface {
	// GetName returns the exchange name
	GetName() string

	// GetTradingPairs fetches all trading pairs
	GetTradingPairs(ctx context.Context) ([]TradingPair, error)

	// SetRateLimit sets rate limiting configuration for specific API type
	SetRateLimit(apiType APIType, limit RateLimit)

	// SetHeaders sets custom request headers
	SetHeaders(headers map[string]string)

	// SetProxies sets proxy list for multi-IP requests
	SetProxies(proxies []string)

	// SetLogger sets custom logger
	SetLogger(logger zerolog.Logger)

	// SetHTTPClient sets custom HTTP client
	SetHTTPClient(client *http.Client)
}

// TradingPair represents a trading pair information
type TradingPair struct {
	Symbol     string  `json:"symbol"`      // Trading pair symbol
	BaseAsset  string  `json:"base_asset"`  // Base asset
	QuoteAsset string  `json:"quote_asset"` // Quote asset
	Status     string  `json:"status"`      // Trading status
	MinQty     float64 `json:"min_qty"`     // Minimum quantity
	MaxQty     float64 `json:"max_qty"`     // Maximum quantity
	StepSize   float64 `json:"step_size"`   // Quantity step size
	TickSize   float64 `json:"tick_size"`   // Price tick size
}

// RateLimit represents rate limiting configuration
type RateLimit struct {
	Requests int           `json:"requests"` // Number of requests
	Interval time.Duration `json:"interval"` // Time interval
}

// RateLimitConfig represents rate limiting configuration for different API types
type RateLimitConfig struct {
	Public  RateLimit `json:"public"`  // Rate limit for public APIs
	Private RateLimit `json:"private"` // Rate limit for private APIs
}

// Config represents exchange configuration
type Config struct {
	APIKey     string            `json:"api_key"`    // API key
	SecretKey  string            `json:"secret_key"` // Secret key
	BaseURL    string            `json:"base_url"`   // Base URL
	Timeout    time.Duration     `json:"timeout"`    // Request timeout
	RateLimit  RateLimitConfig   `json:"rate_limit"` // Rate limiting configuration
	Headers    map[string]string `json:"headers"`    // Custom headers
	Proxies    []string          `json:"proxies"`    // Proxy list
	Testnet    bool              `json:"testnet"`    // Testnet flag
	Sandbox    bool              `json:"sandbox"`    // Sandbox flag (alias for Testnet)
	Logger     *zerolog.Logger   `json:"-"`          // Custom logger (not serialized)
	HTTPClient *http.Client      `json:"-"`          // Custom HTTP client (not serialized)
}
