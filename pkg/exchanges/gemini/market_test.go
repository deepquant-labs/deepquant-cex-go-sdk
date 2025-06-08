package gemini

import (
	"context"
	"testing"
	"time"

	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/exchange"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarketAPI_ListSymbols(t *testing.T) {
	// Create a test configuration
	config := &exchange.Config{
		Testnet: true, // Use sandbox for testing
		Timeout: 30 * time.Second,
		Logger:  &zerolog.Logger{},
	}

	// Create Gemini instance
	gemini := NewGemini(config)
	require.NotNil(t, gemini)
	require.NotNil(t, gemini.Market)

	// Test ListSymbols
	ctx := context.Background()
	symbols, err := gemini.Market.ListSymbols(ctx)

	// Assertions
	require.NoError(t, err, "ListSymbols should not return an error")
	require.NotNil(t, symbols, "Symbols should not be nil")
	assert.Greater(t, len(symbols), 0, "Should return at least one symbol")

	// Check that common symbols exist
	commonSymbols := []string{"btcusd", "ethusd", "ltcusd"}
	for _, commonSymbol := range commonSymbols {
		found := false
		for _, symbol := range symbols {
			if symbol == commonSymbol {
				found = true
				break
			}
		}
		assert.True(t, found, "Common symbol %s should be present", commonSymbol)
	}

	t.Logf("Found %d symbols", len(symbols))
	t.Logf("First few symbols: %v", symbols[:min(5, len(symbols))])
}

func TestMarketAPI_GetSymbolDetails(t *testing.T) {
	// Create a test configuration
	config := &exchange.Config{
		Testnet: true, // Use sandbox for testing
		Timeout: 30 * time.Second,
		Logger:  &zerolog.Logger{},
	}

	// Create Gemini instance
	gemini := NewGemini(config)
	require.NotNil(t, gemini)
	require.NotNil(t, gemini.Market)

	// Test GetSymbolDetails for BTCUSD
	ctx := context.Background()
	details, err := gemini.Market.GetSymbolDetails(ctx, "btcusd")

	// Assertions
	require.NoError(t, err, "GetSymbolDetails should not return an error")
	require.NotNil(t, details, "Details should not be nil")
	assert.Equal(t, "BTCUSD", details.Symbol, "Symbol should match")
	assert.NotEmpty(t, details.BaseCurrency, "BaseCurrency should not be empty")
	assert.NotEmpty(t, details.QuoteCurrency, "QuoteCurrency should not be empty")
	assert.NotEmpty(t, details.Status, "Status should not be empty")
	assert.GreaterOrEqual(t, details.TickSize, 0.0, "TickSize should be non-negative")
	assert.NotEmpty(t, details.ProductType, "ProductType should not be empty")

	t.Logf("Symbol details for BTCUSD: %+v", details)
}

func TestMarketAPI_GetTickerV2(t *testing.T) {
	// Create a test configuration
	config := &exchange.Config{
		Testnet: true, // Use sandbox for testing
		Timeout: 30 * time.Second,
		Logger:  &zerolog.Logger{},
	}

	// Create Gemini instance
	gemini := NewGemini(config)
	require.NotNil(t, gemini)
	require.NotNil(t, gemini.Market)

	// Test GetTickerV2 for BTCUSD
	ctx := context.Background()
	ticker, err := gemini.Market.GetTickerV2(ctx, "btcusd")

	// Assertions
	require.NoError(t, err, "GetTickerV2 should not return an error")
	require.NotNil(t, ticker, "Ticker should not be nil")
	assert.Equal(t, "BTCUSD", ticker.Symbol, "Symbol should match")
	assert.NotEmpty(t, ticker.Bid, "Bid should not be empty")
	assert.NotEmpty(t, ticker.Ask, "Ask should not be empty")

	t.Logf("Ticker for BTCUSD: %+v", ticker)
}

// Helper function for min (Go 1.21+)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
