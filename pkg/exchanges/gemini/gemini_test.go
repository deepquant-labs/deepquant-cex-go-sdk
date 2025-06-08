package gemini

import (
	"context"
	"testing"
	"time"

	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/exchange"
)

func TestNewGemini(t *testing.T) {
	// Test with nil config
	g := NewGemini(nil)
	if g == nil {
		t.Error("Expected non-nil Gemini instance")
	}
	if g.GetName() != "gemini" {
		t.Errorf("Expected name 'gemini', got '%s'", g.GetName())
	}
	if g.baseURL != "https://api.gemini.com" {
		t.Errorf("Expected production URL, got '%s'", g.baseURL)
	}

	// Test with testnet config
	config := &exchange.Config{
		Testnet:   true,
		APIKey:    "test-key",
		SecretKey: "test-secret",
		// UserAgent can be set via headers
	}
	g = NewGemini(config)
	if g.baseURL != "https://api.sandbox.gemini.com" {
		t.Errorf("Expected sandbox URL, got '%s'", g.baseURL)
	}
	if g.apiKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got '%s'", g.apiKey)
	}
	// UserAgent is now set to default value since it's not configurable via Config
	if g.userAgent != "CEX-SDK/1.0" {
		t.Errorf("Expected default user agent 'CEX-SDK/1.0', got '%s'", g.userAgent)
	}
}

func TestGemini_GetName(t *testing.T) {
	g := NewGemini(nil)
	name := g.GetName()
	if name != "gemini" {
		t.Errorf("Expected 'gemini', got '%s'", name)
	}
}

func TestGemini_SetRateLimit(t *testing.T) {
	g := NewGemini(nil)

	// Should not panic
	rateLimit := exchange.RateLimit{
		Requests: 10,
		Interval: time.Second,
	}
	g.SetRateLimit(exchange.APITypePublic, rateLimit)
	g.SetRateLimit(exchange.APITypePrivate, rateLimit)
}

func TestGemini_SetHeaders(t *testing.T) {
	g := NewGemini(nil)

	headers := map[string]string{
		"X-Custom-Header": "test-value",
	}

	// Should not panic
	g.SetHeaders(headers)
}

func TestGemini_SetProxies(t *testing.T) {
	g := NewGemini(nil)

	proxies := []string{"http://proxy1:8080", "http://proxy2:8080"}

	// Should not panic
	g.SetProxies(proxies)
}

func TestGemini_SetAPICredentials(t *testing.T) {
	g := NewGemini(nil)

	g.SetAPICredentials("new-key", "new-secret")

	if g.apiKey != "new-key" {
		t.Errorf("Expected API key 'new-key', got '%s'", g.apiKey)
	}
	if g.apiSecret != "new-secret" {
		t.Errorf("Expected API secret 'new-secret', got '%s'", g.apiSecret)
	}
}

func TestGemini_SetSandbox(t *testing.T) {
	g := NewGemini(nil)

	// Test enabling sandbox
	g.SetSandbox(true)
	if g.baseURL != "https://api.sandbox.gemini.com" {
		t.Errorf("Expected sandbox URL, got '%s'", g.baseURL)
	}
	if !g.sandbox {
		t.Error("Expected sandbox to be true")
	}

	// Test disabling sandbox
	g.SetSandbox(false)
	if g.baseURL != "https://api.gemini.com" {
		t.Errorf("Expected production URL, got '%s'", g.baseURL)
	}
	if g.sandbox {
		t.Error("Expected sandbox to be false")
	}
}

func TestGemini_ValidateConfig(t *testing.T) {
	// Test with valid config
	g := NewGemini(nil)
	err := g.ValidateConfig()
	// This might fail due to network issues, so we just check it doesn't panic
	_ = err

	// Test with invalid URL
	g.baseURL = "invalid-url"
	err = g.ValidateConfig()
	if err == nil {
		t.Error("Expected error for invalid URL")
	}

	// Test with empty URL
	g.baseURL = ""
	err = g.ValidateConfig()
	if err == nil {
		t.Error("Expected error for empty URL")
	}
}

func TestExtractBaseCurrency(t *testing.T) {
	tests := []struct {
		symbol   string
		expected string
	}{
		{"btcusd", "BTC"},
		{"ethusd", "ETH"},
		{"ltcbtc", "LTC"},
		{"ethbtc", "ETH"},
		{"dogusd", "DOG"},
		{"adausd", "ADA"},
		{"BTCUSD", "BTC"},
		{"short", "SHORT"}, // fallback case
	}

	for _, test := range tests {
		result := extractBaseCurrency(test.symbol)
		if result != test.expected {
			t.Errorf("extractBaseCurrency(%s) = %s, expected %s", test.symbol, result, test.expected)
		}
	}
}

func TestExtractQuoteCurrency(t *testing.T) {
	tests := []struct {
		symbol   string
		expected string
	}{
		{"btcusd", "USD"},
		{"ethusd", "USD"},
		{"ltcbtc", "BTC"},
		{"ethbtc", "BTC"},
		{"dogusd", "USD"},
		{"adaeth", "ETH"},
		{"BTCUSD", "USD"},
		{"short", "USD"}, // fallback case
	}

	for _, test := range tests {
		result := extractQuoteCurrency(test.symbol)
		if result != test.expected {
			t.Errorf("extractQuoteCurrency(%s) = %s, expected %s", test.symbol, result, test.expected)
		}
	}
}

func TestParseFloatFromString(t *testing.T) {
	tests := []struct {
		input     string
		expected  float64
		shouldErr bool
	}{
		{"123.45", 123.45, false},
		{"0", 0, false},
		{"", 0, false},
		{"  123.45  ", 123.45, false},
		{"invalid", 0, true},
	}

	for _, test := range tests {
		result, err := parseFloatFromString(test.input)
		if test.shouldErr {
			if err == nil {
				t.Errorf("parseFloatFromString(%s) expected error but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("parseFloatFromString(%s) unexpected error: %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("parseFloatFromString(%s) = %f, expected %f", test.input, result, test.expected)
			}
		}
	}
}

// Integration test - skip by default to avoid network dependency
func TestGemini_GetTradingPairs_Integration(t *testing.T) {
	t.Skip("Skipping integration test to avoid network dependency")

	g := NewGemini(nil)
	ctx := context.Background()
	pairs, err := g.GetTradingPairs(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if len(pairs) == 0 {
		t.Error("Expected at least one trading pair")
	}

	// Check if we have some expected pairs
	foundBTCUSD := false
	for _, pair := range pairs {
		if pair.Symbol == "BTCUSD" {
			foundBTCUSD = true
			if pair.BaseAsset != "BTC" {
				t.Errorf("Expected base currency BTC, got %s", pair.BaseAsset)
			}
			if pair.QuoteAsset != "USD" {
				t.Errorf("Expected quote currency USD, got %s", pair.QuoteAsset)
			}
			break
		}
	}

	if !foundBTCUSD {
		t.Error("Expected to find BTCUSD pair")
	}
}
