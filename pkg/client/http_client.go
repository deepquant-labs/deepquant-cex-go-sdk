package client

import (
	"context"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/deepquant-labs/deepquant-cex-go-sdk/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

// APIType represents the type of API endpoint
type APIType string

const (
	APITypePublic  APIType = "public"
	APITypePrivate APIType = "private"
)

// HTTPClient HTTP client wrapper with rate limiting and proxy support
type HTTPClient struct {
	client         *fasthttp.Client
	customClient   *http.Client
	publicLimiter  *RateLimiter
	privateLimiter *RateLimiter
	headers        map[string]string
	proxies        []string
	logger         zerolog.Logger
	mu             sync.RWMutex
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &fasthttp.Client{
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
		},
		headers: make(map[string]string),
		proxies: make([]string, 0),
		logger:  zerolog.Nop(), // Default no-op logger
	}
}

// SetRateLimit sets rate limiting configuration for specific API type
func (c *HTTPClient) SetRateLimit(apiType APIType, requests int, interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch apiType {
	case APITypePublic:
		c.publicLimiter = NewRateLimiter(requests, interval)
	case APITypePrivate:
		c.privateLimiter = NewRateLimiter(requests, interval)
	}
}

// SetLogger sets custom logger
func (c *HTTPClient) SetLogger(logger zerolog.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = logger
}

// SetCustomHTTPClient sets custom HTTP client
func (c *HTTPClient) SetCustomHTTPClient(client *http.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.customClient = client
}

// SetHeaders sets custom request headers
func (c *HTTPClient) SetHeaders(headers map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range headers {
		c.headers[k] = v
	}
}

// SetProxies sets proxy list for multi-IP requests
func (c *HTTPClient) SetProxies(proxies []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.proxies = make([]string, len(proxies))
	copy(c.proxies, proxies)
}

// Get sends a GET request (public API by default)
func (c *HTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	return c.RequestWithType(ctx, "GET", url, nil, APITypePublic)
}

// Post sends a POST request (private API by default)
func (c *HTTPClient) Post(ctx context.Context, url string, body []byte) ([]byte, error) {
	return c.RequestWithType(ctx, "POST", url, body, APITypePrivate)
}

// GetWithType sends a GET request with specified API type
func (c *HTTPClient) GetWithType(ctx context.Context, url string, apiType APIType) ([]byte, error) {
	return c.RequestWithType(ctx, "GET", url, nil, apiType)
}

// PostWithHeaders sends a POST request with custom headers
func (c *HTTPClient) PostWithHeaders(ctx context.Context, url string, body []byte, headers map[string]string, apiType APIType) ([]byte, error) {
	return c.requestWithHeaders(ctx, "POST", url, body, headers, apiType)
}

// RequestWithType sends HTTP request with specified API type
func (c *HTTPClient) RequestWithType(ctx context.Context, method, url string, body []byte, apiType APIType) ([]byte, error) {
	return c.request(ctx, method, url, body, apiType)
}

// requestWithHeaders sends HTTP request with custom headers
func (c *HTTPClient) requestWithHeaders(ctx context.Context, method, url string, body []byte, headers map[string]string, apiType APIType) ([]byte, error) {
	c.mu.RLock()
	logger := c.logger
	c.mu.RUnlock()

	// Log request
	logger.Debug().Str("method", method).Str("url", url).Str("apiType", string(apiType)).Msg("Sending HTTP request with custom headers")

	// Apply rate limiting based on API type
	var rateLimiter *RateLimiter
	c.mu.RLock()
	switch apiType {
	case APITypePublic:
		rateLimiter = c.publicLimiter
	case APITypePrivate:
		rateLimiter = c.privateLimiter
	}
	c.mu.RUnlock()

	if rateLimiter != nil {
		if err := rateLimiter.Wait(ctx); err != nil {
			logger.Error().Err(err).Msg("Rate limit error")
			return nil, errors.Wrap(errors.ErrRateLimit, "rate limit error", err)
		}
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// Set request method and URL
	req.Header.SetMethod(method)
	req.SetRequestURI(url)

	// Set request body
	if body != nil {
		req.SetBody(body)
	}

	// Set default headers first
	c.mu.RLock()
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	proxies := make([]string, len(c.proxies))
	copy(proxies, c.proxies)
	c.mu.RUnlock()

	// Override with custom headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Select client (with or without proxy)
	client := c.client
	if len(proxies) > 0 {
		proxy := proxies[rand.Intn(len(proxies))]
		client = &fasthttp.Client{
			ReadTimeout:  c.client.ReadTimeout,
			WriteTimeout: c.client.WriteTimeout,
			Dial: func(addr string) (net.Conn, error) {
				return fasthttp.DialTimeout(proxy, time.Second*10)
			},
		}
	}

	// Send request
	start := time.Now()
	err := client.DoTimeout(req, resp, c.client.ReadTimeout)
	duration := time.Since(start)

	if err != nil {
		logger.Error().Err(err).Dur("duration", duration).Msg("Request failed")
		return nil, errors.Wrap(errors.ErrNetworkError, "request failed", err)
	}

	// Log response
	logger.Debug().Int("status", resp.StatusCode()).Dur("duration", duration).Msg("Received HTTP response")

	// Check response status
	if resp.StatusCode() != fasthttp.StatusOK {
		logger.Error().Int("status", resp.StatusCode()).Bytes("body", resp.Body()).Msg("HTTP error response")
		return nil, errors.Newf(errors.ErrNetworkError, "HTTP error: %d %s", resp.StatusCode(), resp.Body())
	}

	logger.Debug().Int("bodySize", len(resp.Body())).Msg("Request completed successfully")
	return resp.Body(), nil
}

// request sends HTTP request with rate limiting and proxy support
func (c *HTTPClient) request(ctx context.Context, method, url string, body []byte, apiType APIType) ([]byte, error) {
	c.mu.RLock()
	logger := c.logger
	c.mu.RUnlock()

	// Log request
	logger.Debug().Str("method", method).Str("url", url).Str("apiType", string(apiType)).Msg("Sending HTTP request")

	// Apply rate limiting based on API type
	var rateLimiter *RateLimiter
	c.mu.RLock()
	switch apiType {
	case APITypePublic:
		rateLimiter = c.publicLimiter
	case APITypePrivate:
		rateLimiter = c.privateLimiter
	}
	c.mu.RUnlock()

	if rateLimiter != nil {
		if err := rateLimiter.Wait(ctx); err != nil {
			logger.Error().Err(err).Msg("Rate limit error")
			return nil, errors.Wrap(errors.ErrRateLimit, "rate limit error", err)
		}
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// Set request method and URL
	req.Header.SetMethod(method)
	req.SetRequestURI(url)

	// Set request body
	if body != nil {
		req.SetBody(body)
		req.Header.SetContentType("application/json")
	}

	// Set custom headers
	c.mu.RLock()
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	proxies := make([]string, len(c.proxies))
	copy(proxies, c.proxies)
	c.mu.RUnlock()

	// Select client (with or without proxy)
	client := c.client
	if len(proxies) > 0 {
		proxy := proxies[rand.Intn(len(proxies))]
		client = &fasthttp.Client{
			ReadTimeout:  c.client.ReadTimeout,
			WriteTimeout: c.client.WriteTimeout,
			Dial: func(addr string) (net.Conn, error) {
				return fasthttp.DialTimeout(proxy, time.Second*10)
			},
		}
	}

	// Send request
	start := time.Now()
	err := client.DoTimeout(req, resp, c.client.ReadTimeout)
	duration := time.Since(start)

	if err != nil {
		logger.Error().Err(err).Dur("duration", duration).Msg("Request failed")
		return nil, errors.Wrap(errors.ErrNetworkError, "request failed", err)
	}

	// Log response
	logger.Debug().Int("status", resp.StatusCode()).Dur("duration", duration).Msg("Received HTTP response")

	// Check response status
	if resp.StatusCode() != fasthttp.StatusOK {
		logger.Error().Int("status", resp.StatusCode()).Bytes("body", resp.Body()).Msg("HTTP error response")
		return nil, errors.Newf(errors.ErrNetworkError, "HTTP error: %d %s", resp.StatusCode(), resp.Body())
	}

	logger.Debug().Int("bodySize", len(resp.Body())).Msg("Request completed successfully")
	return resp.Body(), nil
}
