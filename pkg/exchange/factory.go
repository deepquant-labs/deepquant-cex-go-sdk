package exchange

import (
	"strings"

	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/errors"
)

// Factory creates exchange instances
type Factory struct {
	constructors map[string]func(Config) Exchange
}

// NewFactory creates a new exchange factory
func NewFactory() *Factory {
	return &Factory{
		constructors: make(map[string]func(Config) Exchange),
	}
}

// Register registers an exchange constructor
func (f *Factory) Register(exchangeName string, constructor func(Config) Exchange) {
	f.constructors[strings.ToLower(exchangeName)] = constructor
}

// Create creates an exchange instance by name
func (f *Factory) Create(exchangeName string, config Config) (Exchange, error) {
	name := strings.ToLower(exchangeName)
	constructor, exists := f.constructors[name]
	if !exists {
		return nil, errors.Newf(errors.ErrExchangeNotSupported, "exchange '%s' not supported", exchangeName)
	}
	return constructor(config), nil
}

// CreateByName creates an exchange instance by name (case-insensitive)
func (f *Factory) CreateByName(name string, config Config) (Exchange, error) {
	return f.Create(name, config)
}

// GetSupportedExchanges returns list of supported exchanges
func (f *Factory) GetSupportedExchanges() []string {
	exchanges := make([]string, 0, len(f.constructors))
	for exchangeName := range f.constructors {
		exchanges = append(exchanges, exchangeName)
	}
	return exchanges
}
