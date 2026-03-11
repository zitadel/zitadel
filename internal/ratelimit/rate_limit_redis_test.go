package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestRedisRateLimiter_Check(t *testing.T) {
	t.Parallel()

	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	defer client.Close()

	limiter := NewRedisRateLimiter(client)
	ctx := context.Background()

	count, allowed := limiter.Check(ctx, "tenant:rule:key", 5*time.Second, 1)
	if count != 1 || !allowed {
		t.Fatalf("first check = (%d, %v), want (1, true)", count, allowed)
	}

	count, allowed = limiter.Check(ctx, "tenant:rule:key", 5*time.Second, 1)
	if count != 2 || allowed {
		t.Fatalf("second check = (%d, %v), want (2, false)", count, allowed)
	}

	server.FastForward(6 * time.Second)

	count, allowed = limiter.Check(ctx, "tenant:rule:key", 5*time.Second, 1)
	if count != 1 || !allowed {
		t.Fatalf("post-expiry check = (%d, %v), want (1, true)", count, allowed)
	}
}
