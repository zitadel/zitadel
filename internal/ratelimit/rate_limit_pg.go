package ratelimit

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

// PGRateLimiter is a [RateLimiterStore] backed by an UNLOGGED PostgreSQL table.
// Counters are shared across all ZITADEL instances via the database, making
// this suitable for multi-node deployments without Redis. The table uses
// INSERT ON CONFLICT for atomic upsert and Prune deletes expired windows.
//
// The table must exist (created by migration 70.sql):
//
//	CREATE UNLOGGED TABLE IF NOT EXISTS signals.rate_limit_counters (
//	    key          TEXT        NOT NULL,
//	    count        INTEGER     NOT NULL DEFAULT 1,
//	    window_start TIMESTAMPTZ NOT NULL,
//	    window_secs  INTEGER     NOT NULL,
//	    PRIMARY KEY (key)
//	);
type PGRateLimiter struct {
	db *sql.DB
}

// NewPGRateLimiter creates a PG-backed rate limiter.
func NewPGRateLimiter(db *sql.DB) *PGRateLimiter {
	return &PGRateLimiter{db: db}
}

// Check implements [RateLimiterStore]. Uses INSERT ... ON CONFLICT to
// atomically create or increment the counter for the given key.
func (rl *PGRateLimiter) Check(ctx context.Context, key string, window time.Duration, max int) (count int, allowed bool) {
	now := time.Now().UTC()
	windowSecs := int(window.Seconds())
	if windowSecs < 1 {
		windowSecs = 1
	}

	// Atomic upsert: insert if key doesn't exist, or increment if the window
	// is still active. If the window has expired, reset the counter.
	err := rl.db.QueryRowContext(ctx,
		`INSERT INTO signals.rate_limit_counters (key, count, window_start, window_secs)
		 VALUES ($1, 1, $2, $3)
		 ON CONFLICT (key) DO UPDATE SET
		     count = CASE
		         WHEN signals.rate_limit_counters.window_start + (signals.rate_limit_counters.window_secs || ' seconds')::INTERVAL <= $2
		         THEN 1
		         ELSE signals.rate_limit_counters.count + 1
		     END,
		     window_start = CASE
		         WHEN signals.rate_limit_counters.window_start + (signals.rate_limit_counters.window_secs || ' seconds')::INTERVAL <= $2
		         THEN $2
		         ELSE signals.rate_limit_counters.window_start
		     END,
		     window_secs = $3
		 RETURNING count`,
		key, now, windowSecs,
	).Scan(&count)
	if err != nil {
		logging.WithError(ctx, err).Warn("risk.ratelimit.pg_failed",
			slog.String("key", key),
		)
		// Fail open: allow the request if PG is unavailable.
		return 0, true
	}

	return count, count <= max
}

// Prune implements [RateLimiterStore]. Deletes counters whose window has expired.
func (rl *PGRateLimiter) Prune(ctx context.Context) {
	now := time.Now().UTC()
	result, err := rl.db.ExecContext(ctx,
		`DELETE FROM signals.rate_limit_counters
		 WHERE window_start + (window_secs || ' seconds')::INTERVAL <= $1`,
		now,
	)
	if err != nil {
		logging.WithError(ctx, err).Warn("risk.ratelimit.pg_prune_failed")
		return
	}
	affected, _ := result.RowsAffected()
	if affected > 0 {
		logging.Debug(ctx, "risk.ratelimit.pg_pruned",
			slog.Int64("deleted", affected),
		)
	}
}

// Len returns the number of tracked counter keys (for monitoring).
func (rl *PGRateLimiter) Len(ctx context.Context) (int, error) {
	var count int
	err := rl.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM signals.rate_limit_counters`,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count rate limit keys: %w", err)
	}
	return count, nil
}
