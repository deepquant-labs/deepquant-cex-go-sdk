# CEX SDK

[![CI](https://github.com/sekfung/cex-sdk/workflows/CI/badge.svg)](https://github.com/sekfung/cex-sdk/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/codecov/c/github/sekfung/cex-sdk)](https://codecov.io/gh/sekfung/cex-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/sekfung/cex-sdk)](https://goreportcard.com/report/github.com/sekfung/cex-sdk)
[![GoDoc](https://godoc.org/github.com/sekfung/cex-sdk?status.svg)](https://godoc.org/github.com/sekfung/cex-sdk)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A comprehensive Go SDK for cryptocurrency exchange APIs, providing unified interfaces for multiple exchanges.

## Features

- 🏗️ **Unified Interface**: Consistent API across different exchanges
- 🔒 **Type Safety**: Strongly typed Go interfaces
- 🚀 **High Performance**: Built with performance in mind
- 🛡️ **Security**: Secure credential handling
- 📊 **Market Data**: Real-time and historical market data
- 💰 **Trading**: Spot and futures trading support
- 📈 **Portfolio**: Account and portfolio management
- 🔄 **Rate Limiting**: Built-in rate limiting
- 🧪 **Testing**: Comprehensive test coverage
- 📚 **Documentation**: Well-documented APIs

## Supported Exchanges

- [x] **Gemini** - Full support for market data and trading APIs
- [ ] **Binance** - Coming soon
- [ ] **Coinbase** - Coming soon
- [ ] **Kraken** - Coming soon

## Installation

```bash
go get github.com/deepquant-labs/deepquant-cex-go-sdk
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/exchange"
    "github.com/deepquant-labs/deepquant-cex-go-sdk"
)

func main() {
    // Create SDK instance
    sdk := cexsdk.NewSDK()

    // Configure exchange
    config := exchange.Config{
        APIKey:    "your-api-key",
        APISecret: "your-api-secret",
        Testnet:   true, // Use testnet for development
        Timeout:   30 * time.Second,
    }

    // Create exchange instance
    gemini, err := sdk.NewExchange("gemini", config)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Get market data
    symbols, err := gemini.Market.ListSymbols(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Available symbols: %v\n", symbols)

    // Get ticker
    ticker, err := gemini.Market.GetTickerV2(ctx, "btcusd")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("BTC/USD Ticker: %+v\n", ticker)
}
```

### Advanced Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/exchange"
    "github.com/deepquant-labs/deepquant-cex-go-sdk"
    "github.com/rs/zerolog"
)

func main() {
    // Create logger
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

    // Create SDK with custom configuration
    sdk := cexsdk.NewSDK()

    // Advanced configuration
    config := exchange.Config{
        APIKey:     "your-api-key",
        APISecret:  "your-api-secret",
        Passphrase: "your-passphrase", // For exchanges that require it
        Testnet:    false,
        Timeout:    60 * time.Second,
        Logger:     &logger,
        RateLimit: exchange.RateLimit{
            RequestsPerSecond: 10,
            BurstSize:        20,
        },
    }

    // Create exchange instance
    gemini, err := sdk.NewExchange("gemini", config)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Get account balances
    balances, err := gemini.Fund.GetAvailableBalances(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Account balances: %+v\n", balances)

    // Get symbol details
    details, err := gemini.Market.GetSymbolDetails(ctx, "btcusd")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("BTC/USD Details: %+v\n", details)
}
```

## API Reference

### Market Data

- `ListSymbols(ctx)` - Get all available trading symbols
- `GetTickerV2(ctx, symbol)` - Get ticker data for a symbol
- `GetSymbolDetails(ctx, symbol)` - Get detailed information about a symbol
- `GetAllSymbolDetails(ctx)` - Get details for all symbols

### Account & Funds

- `GetAvailableBalances(ctx)` - Get account balances

## Configuration

### Exchange Configuration

```go
type Config struct {
    APIKey     string        // API key
    APISecret  string        // API secret
    Passphrase string        // Passphrase (for some exchanges)
    Testnet    bool          // Use testnet/sandbox
    Timeout    time.Duration // Request timeout
    Logger     *zerolog.Logger // Logger instance
    RateLimit  RateLimit     // Rate limiting configuration
}
```

### Rate Limiting

```go
type RateLimit struct {
    RequestsPerSecond int // Requests per second
    BurstSize        int // Burst size
}
```

## Error Handling

The SDK provides structured error handling:

```go
import "github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/errors"

// Check for specific error types
if errors.Is(err, errors.ErrInvalidCredentials) {
    // Handle authentication error
}

if errors.Is(err, errors.ErrRateLimitExceeded) {
    // Handle rate limit error
}

if errors.Is(err, errors.ErrNetworkError) {
    // Handle network error
}
```

## Testing

Run tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Write tests for new features
- Follow Go best practices
- Update documentation
- Ensure CI passes

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you have any questions or need help, please:

1. Check the [documentation](https://godoc.org/github.com/sekfung/cex-sdk)
2. Search [existing issues](https://github.com/sekfung/cex-sdk/issues)
3. Create a [new issue](https://github.com/sekfung/cex-sdk/issues/new)

## Roadmap

- [ ] Add more exchanges (Binance, Coinbase, Kraken)
- [ ] WebSocket support for real-time data
- [ ] Order management APIs
- [ ] Historical data APIs
- [ ] Portfolio analytics
- [ ] Paper trading mode
- [ ] CLI tool

---

**Disclaimer**: This software is for educational and development purposes. Use at your own risk when trading with real funds.