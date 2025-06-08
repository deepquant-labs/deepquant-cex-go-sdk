package gemini

import (
	"strconv"
	"strings"
)

// Symbol represents a trading symbol from Gemini API
type Symbol struct {
	Symbol         string  `json:"symbol"`
	BaseCurrency   string  `json:"base_currency"`
	QuoteCurrency  string  `json:"quote_currency"`
	TickSize       float64 `json:"tick_size"`
	QuoteIncrement float64 `json:"quote_increment"`
	MinOrderSize   string  `json:"min_order_size"`
	Status         string  `json:"status"`
	WrapEnabled    bool    `json:"wrap_enabled"`
}

// TickerV2 represents ticker data from Gemini API v2
type TickerV2 struct {
	Symbol  string   `json:"symbol"`
	Open    string   `json:"open"`
	High    string   `json:"high"`
	Low     string   `json:"low"`
	Close   string   `json:"close"`
	Changes []string `json:"changes"`
	Bid     string   `json:"bid"`
	Ask     string   `json:"ask"`
}

// ErrorResponse represents an error response from Gemini API
type ErrorResponse struct {
	Result  string `json:"result"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// parseFloatFromString safely converts string to float64 with error handling
func parseFloatFromString(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}

	// Remove any whitespace
	s = strings.TrimSpace(s)

	return strconv.ParseFloat(s, 64)
}
