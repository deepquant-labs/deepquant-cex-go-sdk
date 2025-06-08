package gemini

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/client"
	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/errors"
)

// FundAPI handles fund management related operations
type FundAPI struct {
	gemini *Gemini
}

// NewFundAPI creates a new fund API instance
func NewFundAPI(g *Gemini) *FundAPI {
	return &FundAPI{
		gemini: g,
	}
}

// Balance represents account balance information
type Balance struct {
	Type                   string `json:"type"`
	Currency               string `json:"currency"`
	Amount                 string `json:"amount"`
	Available              string `json:"available"`
	AvailableForWithdrawal string `json:"availableForWithdrawal"`
}

// GetAvailableBalancesRequest represents the request payload for getting available balances
type GetAvailableBalancesRequest struct {
	Request string `json:"request"`
	Nonce   string `json:"nonce"`
	Account string `json:"account,omitempty"`
}

// GetAvailableBalances fetches available balances for the account
// This implements the private API: https://docs.gemini.com/rest/fund-management#get-available-balances
func (f *FundAPI) GetAvailableBalances(ctx context.Context, account string) ([]Balance, error) {
	if f.gemini.apiKey == "" || f.gemini.apiSecret == "" {
		return nil, errors.New(errors.ErrInvalidInput, "API key and secret are required for private endpoints")
	}

	endpoint := "/v1/balances"
	url := fmt.Sprintf("%s%s", f.gemini.baseURL, endpoint)

	// Create request payload
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	request := GetAvailableBalancesRequest{
		Request: endpoint,
		Nonce:   nonce,
		Account: account,
	}

	// Marshal request to JSON
	payloadBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to marshal request payload", err)
	}

	// Encode payload to base64
	payload := base64.StdEncoding.EncodeToString(payloadBytes)

	// Create HMAC-SHA384 signature
	mac := hmac.New(sha512.New384, []byte(f.gemini.apiSecret))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Set required headers for private API
	headers := map[string]string{
		"X-GEMINI-APIKEY":    f.gemini.apiKey,
		"X-GEMINI-PAYLOAD":   payload,
		"X-GEMINI-SIGNATURE": signature,
		"Content-Type":       "text/plain",
		"Content-Length":     "0",
		"Cache-Control":      "no-cache",
	}

	f.gemini.logger.Debug().Str("url", url).Str("account", account).Msg("Fetching available balances")

	// Make POST request with authentication headers
	response, err := f.gemini.client.PostWithHeaders(ctx, url, nil, headers, client.APITypePrivate)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNetworkError, "failed to fetch available balances", err)
	}

	// Check for API error response
	var errorResp ErrorResponse
	if err := json.Unmarshal(response, &errorResp); err == nil && errorResp.Result == errorStatus {
		return nil, errors.Newf(errors.ErrAPIError, "Gemini API error: %s - %s", errorResp.Reason, errorResp.Message)
	}

	var balances []Balance
	if err := json.Unmarshal(response, &balances); err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to parse balances response", err)
	}

	f.gemini.logger.Debug().Int("count", len(balances)).Msg("Successfully fetched available balances")
	return balances, nil
}

// NotionalBalance represents notional balance information
type NotionalBalance struct {
	Currency                       string `json:"currency"`
	Amount                         string `json:"amount"`
	AmountNotional                 string `json:"amountNotional"`
	Available                      string `json:"available"`
	AvailableNotional              string `json:"availableNotional"`
	AvailableForWithdrawal         string `json:"availableForWithdrawal"`
	AvailableForWithdrawalNotional string `json:"availableForWithdrawalNotional"`
}

// GetNotionalBalancesRequest represents the request payload for getting notional balances
type GetNotionalBalancesRequest struct {
	Request string `json:"request"`
	Nonce   string `json:"nonce"`
	Account string `json:"account,omitempty"`
}

// GetNotionalBalances fetches notional balances in the specified currency
func (f *FundAPI) GetNotionalBalances(ctx context.Context, currency string, account string) ([]NotionalBalance, error) {
	if f.gemini.apiKey == "" || f.gemini.apiSecret == "" {
		return nil, errors.New(errors.ErrInvalidInput, "API key and secret are required for private endpoints")
	}

	endpoint := fmt.Sprintf("/v1/notionalbalances/%s", currency)
	url := fmt.Sprintf("%s%s", f.gemini.baseURL, endpoint)

	// Create request payload
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	request := GetNotionalBalancesRequest{
		Request: endpoint,
		Nonce:   nonce,
		Account: account,
	}

	// Marshal request to JSON
	payloadBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to marshal request payload", err)
	}

	// Encode payload to base64
	payload := base64.StdEncoding.EncodeToString(payloadBytes)

	// Create HMAC-SHA384 signature
	mac := hmac.New(sha512.New384, []byte(f.gemini.apiSecret))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Set required headers for private API
	headers := map[string]string{
		"X-GEMINI-APIKEY":    f.gemini.apiKey,
		"X-GEMINI-PAYLOAD":   payload,
		"X-GEMINI-SIGNATURE": signature,
		"Content-Type":       "text/plain",
		"Content-Length":     "0",
		"Cache-Control":      "no-cache",
	}

	f.gemini.logger.Debug().Str("url", url).Str("currency", currency).Str("account", account).Msg("Fetching notional balances")

	// Make POST request with authentication headers
	response, err := f.gemini.client.PostWithHeaders(ctx, url, nil, headers, client.APITypePrivate)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNetworkError, "failed to fetch notional balances", err)
	}

	// Check for API error response
	var errorResp ErrorResponse
	if err := json.Unmarshal(response, &errorResp); err == nil && errorResp.Result == errorStatus {
		return nil, errors.Newf(errors.ErrAPIError, "Gemini API error: %s - %s", errorResp.Reason, errorResp.Message)
	}

	var balances []NotionalBalance
	if err := json.Unmarshal(response, &balances); err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to parse notional balances response", err)
	}

	f.gemini.logger.Debug().Int("count", len(balances)).Str("currency", currency).Msg("Successfully fetched notional balances")
	return balances, nil
}

// DepositAddress represents a deposit address
type DepositAddress struct {
	Address   string `json:"address"`
	Timestamp int64  `json:"timestamp"`
	Label     string `json:"label,omitempty"`
	Memo      string `json:"memo,omitempty"`
	Network   string `json:"network"`
}

// ListDepositAddressesRequest represents the request payload for listing deposit addresses
type ListDepositAddressesRequest struct {
	Request string `json:"request"`
	Nonce   string `json:"nonce"`
	Account string `json:"account,omitempty"`
}

// ListDepositAddresses fetches deposit addresses for the specified network
func (f *FundAPI) ListDepositAddresses(ctx context.Context, network string, account string) ([]DepositAddress, error) {
	if f.gemini.apiKey == "" || f.gemini.apiSecret == "" {
		return nil, errors.New(errors.ErrInvalidInput, "API key and secret are required for private endpoints")
	}

	endpoint := fmt.Sprintf("/v1/addresses/%s", network)
	url := fmt.Sprintf("%s%s", f.gemini.baseURL, endpoint)

	// Create request payload
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	request := ListDepositAddressesRequest{
		Request: endpoint,
		Nonce:   nonce,
		Account: account,
	}

	// Marshal request to JSON
	payloadBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to marshal request payload", err)
	}

	// Encode payload to base64
	payload := base64.StdEncoding.EncodeToString(payloadBytes)

	// Create HMAC-SHA384 signature
	mac := hmac.New(sha512.New384, []byte(f.gemini.apiSecret))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Set required headers for private API
	headers := map[string]string{
		"X-GEMINI-APIKEY":    f.gemini.apiKey,
		"X-GEMINI-PAYLOAD":   payload,
		"X-GEMINI-SIGNATURE": signature,
		"Content-Type":       "text/plain",
		"Content-Length":     "0",
		"Cache-Control":      "no-cache",
	}

	f.gemini.logger.Debug().Str("url", url).Str("network", network).Str("account", account).Msg("Listing deposit addresses")

	// Make POST request with authentication headers
	response, err := f.gemini.client.PostWithHeaders(ctx, url, nil, headers, client.APITypePrivate)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNetworkError, "failed to list deposit addresses", err)
	}

	// Check for API error response
	var errorResp ErrorResponse
	if err := json.Unmarshal(response, &errorResp); err == nil && errorResp.Result == errorStatus {
		return nil, errors.Newf(errors.ErrAPIError, "Gemini API error: %s - %s", errorResp.Reason, errorResp.Message)
	}

	var addresses []DepositAddress
	if err := json.Unmarshal(response, &addresses); err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to parse deposit addresses response", err)
	}

	f.gemini.logger.Debug().Int("count", len(addresses)).Str("network", network).Msg("Successfully listed deposit addresses")
	return addresses, nil
}
