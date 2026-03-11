package ratelimit

import "fmt"

// Mode selects the rate limiter backend.
type Mode string

const (
	// ModeMemory uses in-memory sharded counters (default).
	// Fast but not shared across instances.
	ModeMemory Mode = "memory"
	// ModeRedis uses Redis for shared rate limit counters.
	// Requires a Redis connection. Counters expire via TTL.
	ModeRedis Mode = "redis"
	// ModePG uses an UNLOGGED PostgreSQL table for shared counters.
	// Works without Redis, suitable for small multi-instance deployments.
	ModePG Mode = "pg"
)

// Config configures the rate limiter backend used by rate_limit rules.
type Config struct {
	// Mode selects the backend: "memory" (default), "redis", or "pg".
	Mode Mode
}

// EffectiveMode returns the configured mode, defaulting to memory.
func (c Config) EffectiveMode() Mode {
	if c.Mode == "" {
		return ModeMemory
	}
	return c.Mode
}

func (c Config) Validate() error {
	switch c.Mode {
	case "", ModeMemory, ModeRedis, ModePG:
		return nil
	default:
		return fmt.Errorf(
			"risk rate limit mode must be one of %q, %q or %q",
			ModeMemory,
			ModeRedis,
			ModePG,
		)
	}
}
