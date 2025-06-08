package client

import (
	"context"
	"sync"
	"time"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	tokens     int           // current available tokens
	maxTokens  int           // maximum tokens
	interval   time.Duration // refill interval
	lastRefill time.Time     // last refill time
	mu         sync.Mutex    // mutex for thread safety
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxTokens int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		interval:   interval,
		lastRefill: time.Now(),
	}
}

// Wait waits for a token to become available
func (rl *RateLimiter) Wait(ctx context.Context) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	if elapsed >= rl.interval {
		periods := int(elapsed / rl.interval)
		rl.tokens = min(rl.maxTokens, rl.tokens+periods)
		rl.lastRefill = now
	}

	// If no tokens available, wait
	if rl.tokens <= 0 {
		waitTime := rl.interval - (now.Sub(rl.lastRefill) % rl.interval)
		rl.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			// Continue after waiting
		}

		rl.mu.Lock()
		rl.tokens = 1
		rl.lastRefill = time.Now()
	}

	// Consume a token
	rl.tokens--
	return nil
}

// TryAcquire attempts to acquire a token without waiting
func (rl *RateLimiter) TryAcquire() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	if elapsed >= rl.interval {
		periods := int(elapsed / rl.interval)
		rl.tokens = min(rl.maxTokens, rl.tokens+periods)
		rl.lastRefill = now
	}

	// Check if tokens are available
	if rl.tokens <= 0 {
		return false
	}

	// Consume a token
	rl.tokens--
	return true
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
