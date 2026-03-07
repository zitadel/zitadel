package risk

import (
	"sync"
	"time"
)

// RateLimiter provides a thread-safe in-memory sliding-window rate limiter.
// Each key tracks a count within a fixed time window. When the window expires
// the counter resets.
type RateLimiter struct {
	mu       sync.Mutex
	counters map[string]*windowCounter
}

type windowCounter struct {
	count       int
	windowStart time.Time
	window      time.Duration
	max         int
}

// NewRateLimiter creates a new in-memory rate limiter.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		counters: make(map[string]*windowCounter),
	}
}

// Check increments the counter for key and returns the current count and
// whether the request is still within the limit. The window and max are
// per-rule and provided at check time. Expired counters encountered during
// lookup are lazily evicted to bound memory growth.
func (rl *RateLimiter) Check(key string, window time.Duration, max int, now time.Time) (count int, allowed bool) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	wc, ok := rl.counters[key]
	if !ok || now.Sub(wc.windowStart) >= wc.window {
		// New window — lazily evict the old entry if present.
		rl.counters[key] = &windowCounter{
			count:       1,
			windowStart: now,
			window:      window,
			max:         max,
		}
		return 1, true
	}

	wc.count++
	return wc.count, wc.count <= max
}

// Len returns the number of tracked counter keys (for testing/monitoring).
func (rl *RateLimiter) Len() int {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return len(rl.counters)
}

// Prune removes expired counters. Call periodically to prevent unbounded growth.
func (rl *RateLimiter) Prune(now time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for key, wc := range rl.counters {
		if now.Sub(wc.windowStart) >= wc.window {
			delete(rl.counters, key)
		}
	}
}
