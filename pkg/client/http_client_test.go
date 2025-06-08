package client

import (
	"testing"
	"time"
)

func TestNewHTTPClient(t *testing.T) {
	timeout := 10 * time.Second
	client := NewHTTPClient(timeout)

	if client == nil {
		t.Error("Expected non-nil HTTP client")
		return
	}

	if client.client == nil {
		t.Error("Expected non-nil underlying client")
		return
	}

	if client.client.ReadTimeout != timeout {
		t.Errorf("Expected read timeout %v, got %v", timeout, client.client.ReadTimeout)
	}

	if client.client.WriteTimeout != timeout {
		t.Errorf("Expected write timeout %v, got %v", timeout, client.client.WriteTimeout)
	}
}

func TestHTTPClient_SetRateLimit(t *testing.T) {
	client := NewHTTPClient(10 * time.Second)

	// Test setting rate limit for public API
	client.SetRateLimit(APITypePublic, 10, time.Minute)

	if client.publicLimiter == nil {
		t.Error("Expected public rate limiter to be set")
	}

	// Test setting rate limit for private API
	client.SetRateLimit(APITypePrivate, 20, time.Minute)

	if client.privateLimiter == nil {
		t.Error("Expected private rate limiter to be set")
	}
}

func TestHTTPClient_SetHeaders(t *testing.T) {
	client := NewHTTPClient(10 * time.Second)

	headers := map[string]string{
		"User-Agent": "test-agent",
		"X-Custom":   "test-value",
	}

	client.SetHeaders(headers)

	client.mu.RLock()
	defer client.mu.RUnlock()

	for k, v := range headers {
		if client.headers[k] != v {
			t.Errorf("Expected header %s=%s, got %s", k, v, client.headers[k])
		}
	}
}

func TestHTTPClient_SetProxies(t *testing.T) {
	client := NewHTTPClient(10 * time.Second)

	proxies := []string{
		"http://proxy1:8080",
		"http://proxy2:8080",
	}

	client.SetProxies(proxies)

	client.mu.RLock()
	defer client.mu.RUnlock()

	if len(client.proxies) != len(proxies) {
		t.Errorf("Expected %d proxies, got %d", len(proxies), len(client.proxies))
	}

	for i, proxy := range proxies {
		if client.proxies[i] != proxy {
			t.Errorf("Expected proxy %s, got %s", proxy, client.proxies[i])
		}
	}
}

// TestHTTPClient_Get is skipped to avoid network dependencies in unit tests
// Integration tests should be run separately
func TestHTTPClient_Get(t *testing.T) {
	t.Skip("Skipping network-dependent test")
}

// TestHTTPClient_RateLimitIntegration is skipped to avoid network dependencies
// Rate limiting is tested separately in rate_limiter_test.go
func TestHTTPClient_RateLimitIntegration(t *testing.T) {
	t.Skip("Skipping network-dependent test")
}
