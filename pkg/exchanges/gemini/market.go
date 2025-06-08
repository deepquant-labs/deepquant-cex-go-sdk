package gemini

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/client"
	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/errors"
)

// MarketAPI handles market data related operations
type MarketAPI struct {
	gemini *Gemini
}

// NewMarketAPI creates a new market API instance
func NewMarketAPI(g *Gemini) *MarketAPI {
	return &MarketAPI{
		gemini: g,
	}
}

// ListSymbolsResponse represents the response from list symbols API
type ListSymbolsResponse []string

// ListSymbols fetches all available trading symbols from Gemini
// This implements the public API: https://docs.gemini.com/rest/market-data#list-symbols
func (m *MarketAPI) ListSymbols(ctx context.Context) (ListSymbolsResponse, error) {
	url := fmt.Sprintf("%s/v1/symbols", m.gemini.baseURL)

	m.gemini.logger.Debug().Str("url", url).Msg("Fetching symbols")

	// This is a public API, no authentication required
	response, err := m.gemini.client.GetWithType(ctx, url, client.APITypePublic)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNetworkError, "failed to fetch symbols", err)
	}

	var symbols ListSymbolsResponse
	if err := json.Unmarshal(response, &symbols); err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to parse symbols response", err)
	}

	m.gemini.logger.Debug().Int("count", len(symbols)).Msg("Successfully fetched symbols")
	return symbols, nil
}

// SymbolDetails represents detailed information about a trading symbol
type SymbolDetails struct {
	Symbol                string  `json:"symbol"`
	BaseCurrency          string  `json:"base_currency"`
	QuoteCurrency         string  `json:"quote_currency"`
	TickSize              float64 `json:"tick_size"`
	QuoteIncrement        float64 `json:"quote_increment"`
	MinOrderSize          string  `json:"min_order_size"`
	Status                string  `json:"status"`
	WrapEnabled           bool    `json:"wrap_enabled"`
	ProductType           string  `json:"product_type"`
	ContractType          string  `json:"contract_type"`
	ContractPriceCurrency string  `json:"contract_price_currency"`
}

// GetSymbolDetails fetches detailed information for a specific symbol
func (m *MarketAPI) GetSymbolDetails(ctx context.Context, symbol string) (*SymbolDetails, error) {
	url := fmt.Sprintf("%s/v1/symbols/details/%s", m.gemini.baseURL, symbol)

	m.gemini.logger.Debug().Str("url", url).Str("symbol", symbol).Msg("Fetching symbol details")

	// This is a public API, no authentication required
	response, err := m.gemini.client.GetWithType(ctx, url, client.APITypePublic)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNetworkError, "failed to fetch symbol details", err)
	}

	var details SymbolDetails
	if err := json.Unmarshal(response, &details); err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to parse symbol details response", err)
	}

	m.gemini.logger.Debug().Str("symbol", symbol).Msg("Successfully fetched symbol details")
	return &details, nil
}

// GetAllSymbolDetails fetches detailed information for all symbols
func (m *MarketAPI) GetAllSymbolDetails(ctx context.Context) ([]SymbolDetails, error) {
	// First get all symbols
	symbols, err := m.ListSymbols(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNetworkError, "failed to fetch symbols list", err)
	}

	allDetails := make([]SymbolDetails, 0, len(symbols))
	for _, symbol := range symbols {
		details, err := m.GetSymbolDetails(ctx, symbol)
		if err != nil {
			m.gemini.logger.Warn().Str("symbol", symbol).Err(err).Msg("Failed to fetch details for symbol")
			continue
		}
		allDetails = append(allDetails, *details)
	}

	m.gemini.logger.Debug().Int("count", len(allDetails)).Msg("Successfully fetched all symbol details")
	return allDetails, nil
}

// GetTickerV2 fetches ticker data for a specific symbol
func (m *MarketAPI) GetTickerV2(ctx context.Context, symbol string) (*TickerV2, error) {
	url := fmt.Sprintf("%s/v2/ticker/%s", m.gemini.baseURL, symbol)

	m.gemini.logger.Debug().Str("url", url).Str("symbol", symbol).Msg("Fetching ticker data")

	// This is a public API, no authentication required
	response, err := m.gemini.client.GetWithType(ctx, url, client.APITypePublic)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNetworkError, "failed to fetch ticker data", err)
	}

	var ticker TickerV2
	if err := json.Unmarshal(response, &ticker); err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to parse ticker response", err)
	}

	m.gemini.logger.Debug().Str("symbol", symbol).Msg("Successfully fetched ticker data")
	return &ticker, nil
}
