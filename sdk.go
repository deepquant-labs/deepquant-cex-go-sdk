package cexsdk

import (
	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/exchange"
	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/exchanges/gemini"
)

// SDK main SDK struct
type SDK struct {
	factory *exchange.Factory
}

// New creates a new SDK instance with default factory
func New() *SDK {
	sdk := &SDK{
		factory: exchange.NewFactory(),
	}

	// Register supported exchanges
	sdk.registerExchanges()

	return sdk
}

// registerExchanges registers all supported exchanges
func (s *SDK) registerExchanges() {
	// Register Gemini
	s.factory.Register("gemini", func(config exchange.Config) exchange.Exchange {
		return gemini.NewGemini(&config)
	})
}

// NewExchange creates a new exchange instance
func (s *SDK) NewExchange(exchangeName string, config exchange.Config) (exchange.Exchange, error) {
	return s.factory.CreateByName(exchangeName, config)
}

// GetSupportedExchanges returns list of supported exchanges
func (s *SDK) GetSupportedExchanges() []string {
	return s.factory.GetSupportedExchanges()
}

// NewGemini creates a new Gemini exchange instance with default configuration
func NewGemini() exchange.Exchange {
	return gemini.NewGemini(nil)
}
