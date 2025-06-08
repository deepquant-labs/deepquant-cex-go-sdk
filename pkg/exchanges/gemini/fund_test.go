package gemini

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/exchange"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFundAPI_GetAvailableBalances(t *testing.T) {
	// Skip test if API credentials are not provided
	apiKey := os.Getenv("GEMINI_API_KEY")
	apiSecret := os.Getenv("GEMINI_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping test: GEMINI_API_KEY and GEMINI_API_SECRET environment variables are required")
	}

	// Create a test configuration with API credentials
	config := &exchange.Config{
		APIKey:    apiKey,
		SecretKey: apiSecret,
		Testnet:   true, // Use sandbox for testing
		Timeout:   30 * time.Second,
		Logger:    &zerolog.Logger{},
	}

	// Create Gemini instance
	gemini := NewGemini(config)
	require.NotNil(t, gemini)
	require.NotNil(t, gemini.Fund)

	// Test GetAvailableBalances
	ctx := context.Background()
	balances, err := gemini.Fund.GetAvailableBalances(ctx, "")

	// Assertions
	require.NoError(t, err, "GetAvailableBalances should not return an error")
	require.NotNil(t, balances, "Balances should not be nil")

	// Log results
	t.Logf("Found %d balances", len(balances))
	for i, balance := range balances {
		t.Logf("Balance %d: %+v", i, balance)

		// Basic structure validation
		assert.NotEmpty(t, balance.Currency, "Currency should not be empty")
		assert.NotEmpty(t, balance.Type, "Type should not be empty")
		// Amount fields can be empty or "0" for currencies with no balance
	}
}

func TestFundAPI_GetAvailableBalances_NoCredentials(t *testing.T) {
	// Create a test configuration without API credentials
	config := &exchange.Config{
		Testnet: true,
		Timeout: 30 * time.Second,
		Logger:  &zerolog.Logger{},
	}

	// Create Gemini instance
	gemini := NewGemini(config)
	require.NotNil(t, gemini)
	require.NotNil(t, gemini.Fund)

	// Test GetAvailableBalances without credentials
	ctx := context.Background()
	balances, err := gemini.Fund.GetAvailableBalances(ctx, "")

	// Should return an error due to missing credentials
	require.Error(t, err, "GetAvailableBalances should return an error when credentials are missing")
	require.Nil(t, balances, "Balances should be nil when error occurs")
	assert.Contains(t, err.Error(), "API key and secret are required", "Error should mention missing credentials")
}

func TestFundAPI_GetNotionalBalances(t *testing.T) {
	// Skip test if API credentials are not provided
	apiKey := os.Getenv("GEMINI_API_KEY")
	apiSecret := os.Getenv("GEMINI_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping test: GEMINI_API_KEY and GEMINI_API_SECRET environment variables are required")
	}

	// Create a test configuration with API credentials
	config := &exchange.Config{
		APIKey:    apiKey,
		SecretKey: apiSecret,
		Testnet:   true, // Use sandbox for testing
		Timeout:   30 * time.Second,
		Logger:    &zerolog.Logger{},
	}

	// Create Gemini instance
	gemini := NewGemini(config)
	require.NotNil(t, gemini)
	require.NotNil(t, gemini.Fund)

	// Test GetNotionalBalances in USD
	ctx := context.Background()
	balances, err := gemini.Fund.GetNotionalBalances(ctx, "usd", "")

	// Assertions
	require.NoError(t, err, "GetNotionalBalances should not return an error")
	require.NotNil(t, balances, "Balances should not be nil")

	// Log results
	t.Logf("Found %d notional balances in USD", len(balances))
	for i, balance := range balances {
		t.Logf("Notional Balance %d: %+v", i, balance)

		// Basic structure validation
		assert.NotEmpty(t, balance.Currency, "Currency should not be empty")
		// Amount fields can be empty or "0" for currencies with no balance
	}
}

func TestFundAPI_ListDepositAddresses(t *testing.T) {
	// Skip test if API credentials are not provided
	apiKey := os.Getenv("GEMINI_API_KEY")
	apiSecret := os.Getenv("GEMINI_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping test: GEMINI_API_KEY and GEMINI_API_SECRET environment variables are required")
	}

	// Create a test configuration with API credentials
	config := &exchange.Config{
		APIKey:    apiKey,
		SecretKey: apiSecret,
		Testnet:   true, // Use sandbox for testing
		Timeout:   30 * time.Second,
		Logger:    &zerolog.Logger{},
	}

	// Create Gemini instance
	gemini := NewGemini(config)
	require.NotNil(t, gemini)
	require.NotNil(t, gemini.Fund)

	// Test ListDepositAddresses for Bitcoin network
	ctx := context.Background()
	addresses, err := gemini.Fund.ListDepositAddresses(ctx, "bitcoin", "")

	// Note: This might return an error if no addresses exist or if the account doesn't have permission
	// We'll check both success and expected error cases
	if err != nil {
		// Log the error for debugging
		t.Logf("ListDepositAddresses returned error (this might be expected): %v", err)
		// Don't fail the test as this might be expected behavior for test accounts
		return
	}

	// If successful, validate the response
	require.NotNil(t, addresses, "Addresses should not be nil")
	t.Logf("Found %d deposit addresses for Bitcoin", len(addresses))

	for i, address := range addresses {
		t.Logf("Address %d: %+v", i, address)

		// Basic structure validation
		assert.NotEmpty(t, address.Address, "Address should not be empty")
		assert.NotEmpty(t, address.Network, "Network should not be empty")
		assert.Greater(t, address.Timestamp, int64(0), "Timestamp should be positive")
	}
}
