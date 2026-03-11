package ratelimit

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMemoryRateLimiter_WindowExpiryBoundary(t *testing.T) {
	t.Parallel()

	mu := sync.Mutex{}
	current := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	rl := NewMemoryRateLimiter()
	rl.now = func() time.Time {
		mu.Lock()
		defer mu.Unlock()
		return current
	}

	ctx := context.Background()
	window := 10 * time.Second
	max := 3

	for i := 1; i <= max; i++ {
		count, allowed := rl.Check(ctx, "key", window, max)
		if count != i || !allowed {
			t.Fatalf("request %d: Check() = (%d, %v), want (%d, true)", i, count, allowed, i)
		}
	}

	// Max+1 should be denied.
	count, allowed := rl.Check(ctx, "key", window, max)
	if count != max+1 || allowed {
		t.Fatalf("over-limit: Check() = (%d, %v), want (%d, false)", count, allowed, max+1)
	}

	// Advance clock to exactly the window duration – counter must reset.
	mu.Lock()
	current = current.Add(window)
	mu.Unlock()

	count, allowed = rl.Check(ctx, "key", window, max)
	if count != 1 || !allowed {
		t.Fatalf("post-expiry: Check() = (%d, %v), want (1, true)", count, allowed)
	}
}

func TestMemoryRateLimiter_WindowNotExpired(t *testing.T) {
	t.Parallel()

	mu := sync.Mutex{}
	current := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	rl := NewMemoryRateLimiter()
	rl.now = func() time.Time {
		mu.Lock()
		defer mu.Unlock()
		return current
	}

	ctx := context.Background()
	window := 10 * time.Second
	max := 1

	count, allowed := rl.Check(ctx, "key", window, max)
	if count != 1 || !allowed {
		t.Fatalf("first: Check() = (%d, %v), want (1, true)", count, allowed)
	}

	// Advance to 1ns before window expiry – counter must NOT reset.
	mu.Lock()
	current = current.Add(window - time.Nanosecond)
	mu.Unlock()

	count, allowed = rl.Check(ctx, "key", window, max)
	if count != 2 || allowed {
		t.Fatalf("before-expiry: Check() = (%d, %v), want (2, false)", count, allowed)
	}
}

func TestMemoryRateLimiter_ConcurrentChecks(t *testing.T) {
	t.Parallel()

	rl := NewMemoryRateLimiter()
	ctx := context.Background()
	window := time.Minute
	max := 50
	n := 100

	var allowedCount atomic.Int64
	var deniedCount atomic.Int64

	var wg sync.WaitGroup
	wg.Add(n)
	start := make(chan struct{})
	for range n {
		go func() {
			defer wg.Done()
			<-start
			_, allowed := rl.Check(ctx, "concurrent-key", window, max)
			if allowed {
				allowedCount.Add(1)
			} else {
				deniedCount.Add(1)
			}
		}()
	}
	close(start)
	wg.Wait()

	if got := allowedCount.Load(); got != int64(max) {
		t.Fatalf("allowed = %d, want %d", got, max)
	}
	if got := deniedCount.Load(); got != int64(n-max) {
		t.Fatalf("denied = %d, want %d", got, n-max)
	}
}

func TestMemoryRateLimiter_Prune(t *testing.T) {
	t.Parallel()

	mu := sync.Mutex{}
	current := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	rl := NewMemoryRateLimiter()
	rl.now = func() time.Time {
		mu.Lock()
		defer mu.Unlock()
		return current
	}

	ctx := context.Background()
	window := 5 * time.Second

	rl.Check(ctx, "a", window, 10)
	rl.Check(ctx, "b", window, 10)
	rl.Check(ctx, "c", window, 10)

	if got := rl.Len(); got != 3 {
		t.Fatalf("Len() before prune = %d, want 3", got)
	}

	// Advance past window and prune.
	mu.Lock()
	current = current.Add(window)
	mu.Unlock()

	rl.Prune(ctx)

	if got := rl.Len(); got != 0 {
		t.Fatalf("Len() after prune = %d, want 0", got)
	}
}

func TestMemoryRateLimiter_DifferentKeys(t *testing.T) {
	t.Parallel()

	rl := NewMemoryRateLimiter()
	ctx := context.Background()
	window := time.Minute
	max := 1

	count1, allowed1 := rl.Check(ctx, "key-alpha", window, max)
	if count1 != 1 || !allowed1 {
		t.Fatalf("key-alpha first: Check() = (%d, %v), want (1, true)", count1, allowed1)
	}

	count2, allowed2 := rl.Check(ctx, "key-beta", window, max)
	if count2 != 1 || !allowed2 {
		t.Fatalf("key-beta first: Check() = (%d, %v), want (1, true)", count2, allowed2)
	}

	// Exhaust key-alpha.
	count1, allowed1 = rl.Check(ctx, "key-alpha", window, max)
	if count1 != 2 || allowed1 {
		t.Fatalf("key-alpha second: Check() = (%d, %v), want (2, false)", count1, allowed1)
	}

	// key-beta must still be at count 1, so second call gives count 2, denied.
	count2, allowed2 = rl.Check(ctx, "key-beta", window, max)
	if count2 != 2 || allowed2 {
		t.Fatalf("key-beta second: Check() = (%d, %v), want (2, false)", count2, allowed2)
	}
}

func TestCanonicalRateLimitKey_Deterministic(t *testing.T) {
	t.Parallel()

	key1 := CanonicalRateLimitKey("rule-1", "inst-1", 5*time.Minute, "user:alice")
	key2 := CanonicalRateLimitKey("rule-1", "inst-1", 5*time.Minute, "user:alice")

	if key1 != key2 {
		t.Fatalf("CanonicalRateLimitKey not deterministic: %q != %q", key1, key2)
	}
	if key1 == "" {
		t.Fatal("CanonicalRateLimitKey returned empty string")
	}
}

func TestCanonicalRateLimitKey_EmptyRenderedKey(t *testing.T) {
	t.Parallel()

	key := CanonicalRateLimitKey("rule-2", "inst-2", time.Minute, "")
	if key == "" {
		t.Fatal("CanonicalRateLimitKey returned empty string for empty rendered key")
	}

	// Must be deterministic even for empty rendered key.
	key2 := CanonicalRateLimitKey("rule-2", "inst-2", time.Minute, "")
	if key != key2 {
		t.Fatalf("empty rendered key not deterministic: %q != %q", key, key2)
	}
}
