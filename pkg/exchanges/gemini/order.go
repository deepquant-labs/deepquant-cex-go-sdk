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

// OrderAPI handles order management related operations
type OrderAPI struct {
	gemini *Gemini
}

// NewOrderAPI creates a new order API instance
func NewOrderAPI(g *Gemini) *OrderAPI {
	return &OrderAPI{
		gemini: g,
	}
}

// OrderSide represents the side of an order
type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

// OrderType represents the type of an order
type OrderType string

const (
	OrderTypeExchangeLimit        OrderType = "exchange limit"
	OrderTypeAuctionOnly          OrderType = "auction-only"
	OrderTypeMarketBuy            OrderType = "market buy"
	OrderTypeMarketSell           OrderType = "market sell"
	OrderTypeImmediateOrCancel    OrderType = "immediate-or-cancel"
	OrderTypeFillOrKill           OrderType = "fill-or-kill"
	OrderTypeIndicationOfInterest OrderType = "indication-of-interest"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusOpen      OrderStatus = "open"
	OrderStatusClosed    OrderStatus = "closed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// NewOrderRequest represents a new order request
type NewOrderRequest struct {
	Request       string    `json:"request"`
	Nonce         string    `json:"nonce"`
	ClientOrderID string    `json:"client_order_id,omitempty"`
	Symbol        string    `json:"symbol"`
	Amount        string    `json:"amount"`
	Price         string    `json:"price,omitempty"`
	Side          OrderSide `json:"side"`
	Type          OrderType `json:"type"`
	Options       []string  `json:"options,omitempty"`
	Account       string    `json:"account,omitempty"`
}

// Order represents an order
type Order struct {
	OrderID           string    `json:"order_id"`
	ID                string    `json:"id"`
	Symbol            string    `json:"symbol"`
	Exchange          string    `json:"exchange"`
	AvgExecutionPrice string    `json:"avg_execution_price"`
	Side              OrderSide `json:"side"`
	Type              OrderType `json:"type"`
	Timestamp         string    `json:"timestamp"`
	Timestampms       int64     `json:"timestampms"`
	IsLive            bool      `json:"is_live"`
	IsCancelled       bool      `json:"is_cancelled"`
	IsHidden          bool      `json:"is_hidden"`
	WasForced         bool      `json:"was_forced"`
	ExecutedAmount    string    `json:"executed_amount"`
	RemainingAmount   string    `json:"remaining_amount"`
	Options           []string  `json:"options"`
	Price             string    `json:"price"`
	OriginalAmount    string    `json:"original_amount"`
	ClientOrderID     string    `json:"client_order_id,omitempty"`
}

// PlaceOrder places a new order
func (o *OrderAPI) PlaceOrder(ctx context.Context, req *NewOrderRequest) (*Order, error) {
	if o.gemini.apiKey == "" || o.gemini.apiSecret == "" {
		return nil, errors.New(errors.ErrInvalidInput, "API key and secret are required for private endpoints")
	}

	endpoint := "/v1/order/new"
	url := fmt.Sprintf("%s%s", o.gemini.baseURL, endpoint)

	// Set request endpoint and nonce
	req.Request = endpoint
	req.Nonce = strconv.FormatInt(time.Now().UnixNano(), 10)

	// Marshal request to JSON
	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to marshal order request", err)
	}

	// Encode payload to base64
	payload := base64.StdEncoding.EncodeToString(payloadBytes)

	// Create HMAC-SHA384 signature
	mac := hmac.New(sha512.New384, []byte(o.gemini.apiSecret))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Set required headers for private API
	headers := map[string]string{
		"X-GEMINI-APIKEY":    o.gemini.apiKey,
		"X-GEMINI-PAYLOAD":   payload,
		"X-GEMINI-SIGNATURE": signature,
		"Content-Type":       "text/plain",
		"Content-Length":     "0",
		"Cache-Control":      "no-cache",
	}

	o.gemini.logger.Debug().Str("url", url).Str("symbol", req.Symbol).Str("side", string(req.Side)).Str("type", string(req.Type)).Msg("Placing order")

	// Make POST request with authentication headers
	response, err := o.gemini.client.PostWithHeaders(ctx, url, nil, headers, client.APITypePrivate)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNetworkError, "failed to place order", err)
	}

	// Check for API error response
	var errorResp ErrorResponse
	if err := json.Unmarshal(response, &errorResp); err == nil && errorResp.Result == errorStatus {
		return nil, errors.Newf(errors.ErrAPIError, "Gemini API error: %s - %s", errorResp.Reason, errorResp.Message)
	}

	var order Order
	if err := json.Unmarshal(response, &order); err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to parse order response", err)
	}

	o.gemini.logger.Debug().Str("order_id", order.OrderID).Msg("Successfully placed order")
	return &order, nil
}

// CancelOrderRequest represents a cancel order request
type CancelOrderRequest struct {
	Request string `json:"request"`
	Nonce   string `json:"nonce"`
	OrderID string `json:"order_id"`
	Account string `json:"account,omitempty"`
}

// CancelOrder cancels an existing order
func (o *OrderAPI) CancelOrder(ctx context.Context, orderID string, account string) (*Order, error) {
	if o.gemini.apiKey == "" || o.gemini.apiSecret == "" {
		return nil, errors.New(errors.ErrInvalidInput, "API key and secret are required for private endpoints")
	}

	endpoint := "/v1/order/cancel"
	url := fmt.Sprintf("%s%s", o.gemini.baseURL, endpoint)

	// Create request payload
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	request := CancelOrderRequest{
		Request: endpoint,
		Nonce:   nonce,
		OrderID: orderID,
		Account: account,
	}

	// Marshal request to JSON
	payloadBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to marshal cancel request", err)
	}

	// Encode payload to base64
	payload := base64.StdEncoding.EncodeToString(payloadBytes)

	// Create HMAC-SHA384 signature
	mac := hmac.New(sha512.New384, []byte(o.gemini.apiSecret))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Set required headers for private API
	headers := map[string]string{
		"X-GEMINI-APIKEY":    o.gemini.apiKey,
		"X-GEMINI-PAYLOAD":   payload,
		"X-GEMINI-SIGNATURE": signature,
		"Content-Type":       "text/plain",
		"Content-Length":     "0",
		"Cache-Control":      "no-cache",
	}

	o.gemini.logger.Debug().Str("url", url).Str("order_id", orderID).Msg("Cancelling order")

	// Make POST request with authentication headers
	response, err := o.gemini.client.PostWithHeaders(ctx, url, nil, headers, client.APITypePrivate)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNetworkError, "failed to cancel order", err)
	}

	// Check for API error response
	var errorResp ErrorResponse
	if err := json.Unmarshal(response, &errorResp); err == nil && errorResp.Result == errorStatus {
		return nil, errors.Newf(errors.ErrAPIError, "Gemini API error: %s - %s", errorResp.Reason, errorResp.Message)
	}

	var order Order
	if err := json.Unmarshal(response, &order); err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to parse cancel order response", err)
	}

	o.gemini.logger.Debug().Str("order_id", orderID).Msg("Successfully cancelled order")
	return &order, nil
}

// GetActiveOrdersRequest represents a request to get active orders
type GetActiveOrdersRequest struct {
	Request string `json:"request"`
	Nonce   string `json:"nonce"`
	Account string `json:"account,omitempty"`
}

// GetActiveOrders fetches all active orders
func (o *OrderAPI) GetActiveOrders(ctx context.Context, account string) ([]Order, error) {
	if o.gemini.apiKey == "" || o.gemini.apiSecret == "" {
		return nil, errors.New(errors.ErrInvalidInput, "API key and secret are required for private endpoints")
	}

	endpoint := "/v1/orders"
	url := fmt.Sprintf("%s%s", o.gemini.baseURL, endpoint)

	// Create request payload
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	request := GetActiveOrdersRequest{
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
	mac := hmac.New(sha512.New384, []byte(o.gemini.apiSecret))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Set required headers for private API
	headers := map[string]string{
		"X-GEMINI-APIKEY":    o.gemini.apiKey,
		"X-GEMINI-PAYLOAD":   payload,
		"X-GEMINI-SIGNATURE": signature,
		"Content-Type":       "text/plain",
		"Content-Length":     "0",
		"Cache-Control":      "no-cache",
	}

	o.gemini.logger.Debug().Str("url", url).Str("account", account).Msg("Fetching active orders")

	// Make POST request with authentication headers
	response, err := o.gemini.client.PostWithHeaders(ctx, url, nil, headers, client.APITypePrivate)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNetworkError, "failed to fetch active orders", err)
	}

	// Check for API error response
	var errorResp ErrorResponse
	if err := json.Unmarshal(response, &errorResp); err == nil && errorResp.Result == errorStatus {
		return nil, errors.Newf(errors.ErrAPIError, "Gemini API error: %s - %s", errorResp.Reason, errorResp.Message)
	}

	var orders []Order
	if err := json.Unmarshal(response, &orders); err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to parse orders response", err)
	}

	o.gemini.logger.Debug().Int("count", len(orders)).Msg("Successfully fetched active orders")
	return orders, nil
}

// GetOrderStatusRequest represents a request to get order status
type GetOrderStatusRequest struct {
	Request       string `json:"request"`
	Nonce         string `json:"nonce"`
	OrderID       string `json:"order_id"`
	ClientOrderID string `json:"client_order_id,omitempty"`
	IncludeTrades bool   `json:"include_trades,omitempty"`
	Account       string `json:"account,omitempty"`
}

// GetOrderStatus fetches the status of a specific order
func (o *OrderAPI) GetOrderStatus(ctx context.Context, orderID string, clientOrderID string, includeTrades bool, account string) (*Order, error) {
	if o.gemini.apiKey == "" || o.gemini.apiSecret == "" {
		return nil, errors.New(errors.ErrInvalidInput, "API key and secret are required for private endpoints")
	}

	endpoint := "/v1/order/status"
	url := fmt.Sprintf("%s%s", o.gemini.baseURL, endpoint)

	// Create request payload
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	request := GetOrderStatusRequest{
		Request:       endpoint,
		Nonce:         nonce,
		OrderID:       orderID,
		ClientOrderID: clientOrderID,
		IncludeTrades: includeTrades,
		Account:       account,
	}

	// Marshal request to JSON
	payloadBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to marshal request payload", err)
	}

	// Encode payload to base64
	payload := base64.StdEncoding.EncodeToString(payloadBytes)

	// Create HMAC-SHA384 signature
	mac := hmac.New(sha512.New384, []byte(o.gemini.apiSecret))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Set required headers for private API
	headers := map[string]string{
		"X-GEMINI-APIKEY":    o.gemini.apiKey,
		"X-GEMINI-PAYLOAD":   payload,
		"X-GEMINI-SIGNATURE": signature,
		"Content-Type":       "text/plain",
		"Content-Length":     "0",
		"Cache-Control":      "no-cache",
	}

	o.gemini.logger.Debug().Str("url", url).Str("order_id", orderID).Str("client_order_id", clientOrderID).Msg("Fetching order status")

	// Make POST request with authentication headers
	response, err := o.gemini.client.PostWithHeaders(ctx, url, nil, headers, client.APITypePrivate)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNetworkError, "failed to fetch order status", err)
	}

	// Check for API error response
	var errorResp ErrorResponse
	if err := json.Unmarshal(response, &errorResp); err == nil && errorResp.Result == errorStatus {
		return nil, errors.Newf(errors.ErrAPIError, "Gemini API error: %s - %s", errorResp.Reason, errorResp.Message)
	}

	var order Order
	if err := json.Unmarshal(response, &order); err != nil {
		return nil, errors.Wrap(errors.ErrDataParsingError, "failed to parse order status response", err)
	}

	o.gemini.logger.Debug().Str("order_id", orderID).Msg("Successfully fetched order status")
	return &order, nil
}
