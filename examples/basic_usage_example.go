package examples

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	cexsdk "github.com/deepquant-labs/deepquant-cex-go-sdk"
	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/exchange"
)

func BasicUsageExample() {

	// Example : Advanced usage with custom configuration
	fmt.Println("\n=== Advanced Usage Example ===")
	advancedUsageExample()

	// Example : Rate limiting demonstration
	fmt.Println("\n=== Rate Limiting Example ===")
	rateLimitingExample()

	// Example : Custom headers example
	fmt.Println("\n=== Custom Headers Example ===")
	customHeadersExample()
}

func advancedUsageExample() {
	// Create SDK instance
	sdk := cexsdk.New()

	// Configure exchange with custom settings
	config := exchange.Config{
		APIKey:    "", // Add your API key here
		SecretKey: "", // Add your secret key here
		Testnet:   true,
		Timeout:   30 * time.Second,
		RateLimit: exchange.RateLimitConfig{
			Public: exchange.RateLimit{
				Requests: 100, // 100 requests per minute
				Interval: time.Minute,
			},
			Private: exchange.RateLimit{
				Requests: 100, // 100 requests per minute
				Interval: time.Minute,
			},
		},
		Headers: map[string]string{
			"User-Agent": "MyTradingBot/1.0",
		},
	}

	// Create exchange instance
	exch, err := sdk.NewExchange("gemini", config)
	if err != nil {
		log.Printf("Error creating exchange: %v", err)
		return
	}

	fmt.Printf("Created exchange: %s\n", exch.GetName())

	// Get trading pairs
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tradingPairs, err := exch.GetTradingPairs(ctx)
	if err != nil {
		log.Printf("Error getting trading pairs: %v", err)
		return
	}

	fmt.Printf("Retrieved %d trading pairs\n", len(tradingPairs))

	// Find and display BTC pairs
	btcPairs := 0
	for _, pair := range tradingPairs {
		if pair.BaseAsset == "BTC" {
			btcPairs++
			if btcPairs <= 3 {
				fmt.Printf("BTC Pair: %s (Min: %.8f, Max: %.8f)\n",
					pair.Symbol, pair.MinQty, pair.MaxQty)
			}
		}
	}
	fmt.Printf("Total BTC pairs: %d\n", btcPairs)
}

func rateLimitingExample() {
	// Create exchange with strict rate limiting
	sdk := cexsdk.New()
	config := exchange.Config{
		Testnet: true,
		Timeout: 10 * time.Second,
		RateLimit: exchange.RateLimitConfig{
			Public: exchange.RateLimit{
				Requests: 2, // Only 2 requests per 10 seconds
				Interval: 10 * time.Second,
			},
			Private: exchange.RateLimit{
				Requests: 2, // Only 2 requests per 10 seconds
				Interval: 10 * time.Second,
			},
		},
	}

	exch, err := sdk.NewExchange("gemini", config)
	if err != nil {
		log.Printf("Error creating exchange: %v", err)
		return
	}

	ctx := context.Background()

	// Make multiple requests to demonstrate rate limiting
	for i := 1; i <= 3; i++ {
		start := time.Now()
		fmt.Printf("Request %d starting...\n", i)

		_, err := exch.GetTradingPairs(ctx)
		duration := time.Since(start)

		if err != nil {
			log.Printf("Request %d failed: %v", i, err)
		} else {
			fmt.Printf("Request %d completed in %v\n", i, duration)
		}
	}
}

func customHeadersExample() {
	// Create exchange and set custom headers
	sdk := cexsdk.New()
	config := exchange.Config{
		Testnet: true,
		Timeout: 15 * time.Second,
	}

	exch, err := sdk.NewExchange("gemini", config)
	if err != nil {
		log.Printf("Error creating exchange: %v", err)
		return
	}

	// Set custom headers
	customHeaders := map[string]string{
		"User-Agent":      "CustomBot/2.0",
		"X-Custom-Header": "MyValue",
		"Accept":          "application/json",
	}
	exch.SetHeaders(customHeaders)

	fmt.Println("Custom headers set successfully")

	// Make request with custom headers
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tradingPairs, err := exch.GetTradingPairs(ctx)
	if err != nil {
		log.Printf("Error with custom headers: %v", err)
		return
	}

	fmt.Printf("Successfully retrieved %d trading pairs with custom headers\n", len(tradingPairs))

	fmt.Println("\n=== Gemini Exchange Example ===")
	// Example 5: Gemini exchange
	gemini, err := sdk.NewExchange("gemini", exchange.Config{Testnet: true})
	if err != nil {
		log.Printf("Error creating Gemini exchange: %v", err)
		return
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel2()

	geminiPairs, err := gemini.GetTradingPairs(ctx2)
	if err != nil {
		log.Printf("Error fetching Gemini pairs: %v", err)
		return
	}

	fmt.Printf("Successfully retrieved %d trading pairs from Gemini\n", len(geminiPairs))

	// Show some Gemini pairs
	geminiCount := 0
	for _, pair := range geminiPairs {
		if strings.Contains(pair.Symbol, "BTC") && geminiCount < 5 {
			fmt.Printf("Gemini BTC Pair: %s (Base: %s, Quote: %s)\n",
				pair.Symbol, pair.BaseAsset, pair.QuoteAsset)
			geminiCount++
		}
	}

	fmt.Println("\n=== Multi-Exchange Comparison ===")
	// Example 6: Compare exchanges using SDK
	supportedExchanges := []string{"gemini"}
	fmt.Printf("Supported exchanges: %v\n", supportedExchanges)

	for _, exchangeName := range supportedExchanges {
		exch, err := sdk.NewExchange(exchangeName, exchange.Config{Testnet: true})
		if err != nil {
			log.Printf("Error creating %s exchange: %v", exchangeName, err)
			continue
		}

		ctx3, cancel3 := context.WithTimeout(context.Background(), 15*time.Second)
		pairs, err := exch.GetTradingPairs(ctx3)
		cancel3()
		if err != nil {
			log.Printf("Error fetching pairs from %s: %v", exchangeName, err)
			continue
		}

		fmt.Printf("%s: %d trading pairs\n", strings.ToUpper(exchangeName), len(pairs))
	}
}
