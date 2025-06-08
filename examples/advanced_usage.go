package examples

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	cexsdk "github.com/deepquant-labs/deepquant-cex-go-sdk"
	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/exchange"
	"github.com/rs/zerolog"
)

func AdvancedUsageExample() {
	// Example 1: Using custom logger
	customLoggerExample()

	// Example 2: Using custom HTTP client
	customHTTPClientExample()

	// Example 3: Using different rate limits for public and private APIs
	rateLimitExample()

	// Example 4: Using sandbox/testnet environment
	sandboxExample()
}

func customLoggerExample() {
	fmt.Println("=== Custom Logger Example ===")

	// Create a custom logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Create SDK with custom logger
	sdk := cexsdk.New()
	config := exchange.Config{
		APIKey:    "your-api-key",
		SecretKey: "your-secret-key",
		Testnet:   true,
		Logger:    &logger,
		RateLimit: exchange.RateLimitConfig{
			Public: exchange.RateLimit{
				Requests: 120, // 120 requests per minute for public APIs
				Interval: time.Minute,
			},
			Private: exchange.RateLimit{
				Requests: 600, // 600 requests per minute for private APIs
				Interval: time.Minute,
			},
		},
	}

	exch, err := sdk.NewExchange("gemini", config)
	if err != nil {
		log.Printf("Error creating exchange: %v", err)
		return
	}

	// Get trading pairs (this will use public API rate limit)
	ctx := context.Background()
	tradingPairs, err := exch.GetTradingPairs(ctx)
	if err != nil {
		log.Printf("Error getting trading pairs: %v", err)
		return
	}

	fmt.Printf("Found %d trading pairs\n", len(tradingPairs))
}

func customHTTPClientExample() {
	fmt.Println("\n=== Custom HTTP Client Example ===")

	// Create a custom HTTP client with specific timeout
	customClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	// Create SDK with custom HTTP client
	sdk := cexsdk.New()
	config := exchange.Config{
		Testnet:    true,
		HTTPClient: customClient,
	}

	exch, err := sdk.NewExchange("gemini", config)
	if err != nil {
		log.Printf("Error creating exchange: %v", err)
		return
	}

	ctx := context.Background()
	tradingPairs, err := exch.GetTradingPairs(ctx)
	if err != nil {
		log.Printf("Error getting trading pairs: %v", err)
		return
	}

	fmt.Printf("Found %d trading pairs using custom HTTP client\n", len(tradingPairs))
}

func rateLimitExample() {
	fmt.Println("\n=== Rate Limit Example ===")

	sdk := cexsdk.New()
	config := exchange.Config{
		Testnet: true,
		RateLimit: exchange.RateLimitConfig{
			Public: exchange.RateLimit{
				Requests: 120, // 120 requests per minute for public APIs
				Interval: time.Minute,
			},
			Private: exchange.RateLimit{
				Requests: 600, // 600 requests per minute for private APIs
				Interval: time.Minute,
			},
		},
	}

	exch, err := sdk.NewExchange("gemini", config)
	if err != nil {
		log.Printf("Error creating exchange: %v", err)
		return
	}

	// You can also update rate limits after creation
	publicLimit := exchange.RateLimit{
		Requests: 100,
		Interval: time.Minute,
	}
	privateLimit := exchange.RateLimit{
		Requests: 500,
		Interval: time.Minute,
	}

	exch.SetRateLimit(exchange.APITypePublic, publicLimit)
	exch.SetRateLimit(exchange.APITypePrivate, privateLimit)

	fmt.Println("Rate limits configured successfully")
}

func sandboxExample() {
	fmt.Println("\n=== Sandbox/Testnet Example ===")

	sdk := cexsdk.New()

	// Using Testnet flag
	testnetConfig := exchange.Config{
		APIKey:    "testnet-api-key",
		SecretKey: "testnet-secret-key",
		Testnet:   true, // This will use testnet endpoints
	}

	exch1, err := sdk.NewExchange("gemini", testnetConfig)
	if err != nil {
		log.Printf("Error creating testnet exchange: %v", err)
		return
	}

	// Using Sandbox flag (alias for Testnet)
	sandboxConfig := exchange.Config{
		APIKey:    "sandbox-api-key",
		SecretKey: "sandbox-secret-key",
		Sandbox:   true, // This will also use testnet endpoints
	}

	exch2, err := sdk.NewExchange("gemini", sandboxConfig)
	if err != nil {
		log.Printf("Error creating sandbox exchange: %v", err)
		return
	}

	fmt.Printf("Testnet exchange: %s\n", exch1.GetName())
	fmt.Printf("Sandbox exchange: %s\n", exch2.GetName())
}
