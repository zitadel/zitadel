package ratelimit

import (
	"context"
	"hash/fnv"
	"sync"
	"time"
)

// RateLimiterStore is the interface for rate limit counter backends.
// Implementations must be safe for concurrent use.
type RateLimiterStore interface {
	// Check atomically increments the counter for key within the given window
	// and returns the current count and whether the request is allowed (count <= max).
	Check(ctx context.Context, key string, window time.Duration, max int) (count int, allowed bool)
	// Prune removes expired counters. Backends with built-in TTL (Redis, PG) may no-op.
	Prune(ctx context.Context)
}

const rateLimiterShards = 64

// MemoryRateLimiter is an in-memory [RateLimiterStore] using sharded maps.
// Suitable for single-node deployments or when rate limits don't need to be
// shared across instances.
type MemoryRateLimiter struct {
	shards [rateLimiterShards]rateShard
	now    func() time.Time
}

type rateShard struct {
	mu       sync.Mutex
	counters map[string]*windowCounter
}

type windowCounter struct {
	count       int
	windowStart time.Time
	window      time.Duration
	max         int
}

// NewMemoryRateLimiter creates a new in-memory rate limiter.
func NewMemoryRateLimiter() *MemoryRateLimiter {
	rl := &MemoryRateLimiter{now: time.Now}
	for i := range rl.shards {
		rl.shards[i].counters = make(map[string]*windowCounter)
	}
	return rl
}

func (rl *MemoryRateLimiter) shard(key string) *rateShard {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return &rl.shards[h.Sum32()%rateLimiterShards]
}

// Check implements [RateLimiterStore].
func (rl *MemoryRateLimiter) Check(_ context.Context, key string, window time.Duration, max int) (count int, allowed bool) {
	now := rl.now()
	s := rl.shard(key)
	s.mu.Lock()
	defer s.mu.Unlock()

	wc, ok := s.counters[key]
	if !ok || now.Sub(wc.windowStart) >= wc.window {
		s.counters[key] = &windowCounter{
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

// Prune implements [RateLimiterStore].
func (rl *MemoryRateLimiter) Prune(_ context.Context) {
	now := rl.now()
	for i := range rl.shards {
		rl.shards[i].mu.Lock()
		for key, wc := range rl.shards[i].counters {
			if now.Sub(wc.windowStart) >= wc.window {
				delete(rl.shards[i].counters, key)
			}
		}
		rl.shards[i].mu.Unlock()
	}
}

// Len returns the number of tracked counter keys (for testing/monitoring).
func (rl *MemoryRateLimiter) Len() int {
	total := 0
	for i := range rl.shards {
		rl.shards[i].mu.Lock()
		total += len(rl.shards[i].counters)
		rl.shards[i].mu.Unlock()
	}
	return total
}
