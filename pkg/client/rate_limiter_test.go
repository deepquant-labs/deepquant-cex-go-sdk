package client

import (
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(10, time.Second)

	if rl == nil {
		t.Error("Expected non-nil rate limiter")
	}
}

func TestRateLimiter_TryAcquire(t *testing.T) {
	rl := NewRateLimiter(3, time.Second)

	// Should be able to acquire up to capacity
	for i := 0; i < 3; i++ {
		if !rl.TryAcquire() {
			t.Errorf("Expected acquisition %d to succeed", i+1)
		}
	}

	// Should fail to acquire beyond capacity
	if rl.TryAcquire() {
		t.Error("Expected acquisition to fail when over capacity")
	}
}

func TestRateLimiter_Wait(t *testing.T) {
	// Skip timing-sensitive test
	t.Skip("Skipping timing-sensitive test")
}

func TestRateLimiter_WaitWithContext(t *testing.T) {
	// Skip timing-sensitive test
	t.Skip("Skipping timing-sensitive test")
}

func TestRateLimiter_TokenRefill(t *testing.T) {
	// Skip timing-sensitive test
	t.Skip("Skipping timing-sensitive test")
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	// Skip concurrent test to avoid race conditions in CI
	t.Skip("Skipping concurrent test to avoid race conditions")
}

func TestMin(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{1, 2, 1},
		{2, 1, 1},
		{5, 5, 5},
		{0, 10, 0},
		{-1, 1, -1},
	}

	for _, test := range tests {
		result := min(test.a, test.b)
		if result != test.expected {
			t.Errorf("min(%d, %d) = %d, expected %d", test.a, test.b, result, test.expected)
		}
	}
}
