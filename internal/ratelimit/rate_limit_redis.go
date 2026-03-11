package ratelimit

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

const rateLimitRedisPrefix = "zitadel:ratelimit:"

// rateLimitScript is a Lua script that atomically increments a fixed-window
// counter in Redis. It sets the key with a TTL equal to the window duration
// on first access and increments on subsequent calls within the same window.
//
// KEYS[1] = counter key
// ARGV[1] = window duration in seconds (used as TTL)
// Returns: current count
var rateLimitScript = redis.NewScript(`
local key = KEYS[1]
local ttl = tonumber(ARGV[1])
local count = redis.call('INCR', key)
if count == 1 then
    redis.call('EXPIRE', key, ttl)
end
return count
`)

// RedisRateLimiter is a [RateLimiterStore] backed by Redis. Counters are
// shared across all ZITADEL instances, making this suitable for multi-node
// deployments. Each key gets a TTL equal to the rate limit window, so Redis
// handles expiry automatically — Prune is a no-op.
type RedisRateLimiter struct {
	client *redis.Client
}

// NewRedisRateLimiter creates a Redis-backed rate limiter.
func NewRedisRateLimiter(client *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{client: client}
}

// Check implements [RateLimiterStore].
func (rl *RedisRateLimiter) Check(ctx context.Context, key string, window time.Duration, max int) (count int, allowed bool) {
	redisKey := rateLimitRedisPrefix + key
	ttlSec := int(window.Seconds())
	if ttlSec < 1 {
		ttlSec = 1
	}

	result, err := rateLimitScript.Run(ctx, rl.client, []string{redisKey}, ttlSec).Int()
	if err != nil {
		logging.WithError(ctx, err).Warn("risk.ratelimit.redis_failed",
			slog.String("key", key),
		)
		// Fail open: allow the request if Redis is unavailable.
		return 0, true
	}

	return result, result <= max
}

// Prune implements [RateLimiterStore]. No-op for Redis since keys expire via TTL.
func (rl *RedisRateLimiter) Prune(_ context.Context) {}

// Len returns the approximate number of rate limit keys in Redis (for monitoring).
func (rl *RedisRateLimiter) Len(ctx context.Context) (int, error) {
	var cursor uint64
	count := 0
	for {
		keys, next, err := rl.client.Scan(ctx, cursor, rateLimitRedisPrefix+"*", 100).Result()
		if err != nil {
			return 0, fmt.Errorf("redis scan rate limit keys: %w", err)
		}
		count += len(keys)
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return count, nil
}
