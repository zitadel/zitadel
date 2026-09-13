package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zitadel/logging"
	"go.opentelemetry.io/otel/metric"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/metrics"
)

const (
	PoolAcquiredConns     = "zitadel.database.pool.acquired_conns"
	PoolConstructingConns = "zitadel.database.pool.constructing_conns"
	PoolIdleConns         = "zitadel.database.pool.idle_conns"
	PoolMaxConns          = "zitadel.database.pool.max_conns"
	PoolTotalConns        = "zitadel.database.pool.total_conns"

	PoolAcquireCount            = "zitadel.database.pool.acquire_count"
	PoolAcquireDurationMs       = "zitadel.database.pool.acquire_duration_milliseconds"
	PoolCanceledAcquireCount    = "zitadel.database.pool.canceled_acquire_count"
	PoolEmptyAcquireCount       = "zitadel.database.pool.empty_acquire_count"
	PoolMaxIdleDestroyCount     = "zitadel.database.pool.max_idle_destroy_count"
	PoolMaxLifetimeDestroyCount = "zitadel.database.pool.max_lifetime_destroy_count"
	PoolNewConnsCount           = "zitadel.database.pool.new_conns_count"
)

// RegisterPoolMetrics registers observable gauges and counters reporting pool's
// underlying *pgxpool.Pool state (see pgxpool.Pool.Stat) with the global metrics
// pipeline.
func RegisterPoolMetrics(ctx context.Context, pool *pgxpool.Pool) {
	registerPoolMetrics(ctx, metrics.GlobalMeter(), pool)
}

func registerPoolMetrics(ctx context.Context, m metrics.Metrics, pool *pgxpool.Pool) {
	gauge := func(name, description string, value func(*pgxpool.Stat) int64) {
		err := m.RegisterValueObserver(name, description, func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(value(pool.Stat()))
			return nil
		})
		logging.OnError(err).WithField("metric", name).Error("failed to register database pool gauge metric")
	}
	counter := func(name, description string, value func(*pgxpool.Stat) int64) {
		err := m.RegisterCounterObserver(name, description, func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(value(pool.Stat()))
			return nil
		})
		logging.OnError(err).WithField("metric", name).Error("failed to register database pool counter metric")
	}

	gauge(PoolAcquiredConns, "Number of currently acquired connections in the pool", func(s *pgxpool.Stat) int64 { return int64(s.AcquiredConns()) })
	gauge(PoolConstructingConns, "Number of connections currently being established", func(s *pgxpool.Stat) int64 { return int64(s.ConstructingConns()) })
	gauge(PoolIdleConns, "Number of currently idle connections in the pool", func(s *pgxpool.Stat) int64 { return int64(s.IdleConns()) })
	gauge(PoolMaxConns, "Maximum size of the pool", func(s *pgxpool.Stat) int64 { return int64(s.MaxConns()) })
	gauge(PoolTotalConns, "Total number of connections currently in the pool", func(s *pgxpool.Stat) int64 { return int64(s.TotalConns()) })

	counter(PoolAcquireCount, "Cumulative count of successful acquires from the pool", func(s *pgxpool.Stat) int64 { return s.AcquireCount() })
	counter(PoolAcquireDurationMs, "Cumulative time spent waiting for successful acquires from the pool, in milliseconds", func(s *pgxpool.Stat) int64 { return s.AcquireDuration().Milliseconds() })
	counter(PoolCanceledAcquireCount, "Cumulative count of acquires from the pool that were canceled by a context", func(s *pgxpool.Stat) int64 { return s.CanceledAcquireCount() })
	counter(PoolEmptyAcquireCount, "Cumulative count of successful acquires that waited for a resource to be released or constructed because the pool was empty", func(s *pgxpool.Stat) int64 { return s.EmptyAcquireCount() })
	counter(PoolMaxIdleDestroyCount, "Cumulative count of connections destroyed because they exceeded MaxConnIdleTime", func(s *pgxpool.Stat) int64 { return s.MaxIdleDestroyCount() })
	counter(PoolMaxLifetimeDestroyCount, "Cumulative count of connections destroyed because they exceeded MaxConnLifetime", func(s *pgxpool.Stat) int64 { return s.MaxLifetimeDestroyCount() })
	counter(PoolNewConnsCount, "Cumulative count of new connections opened", func(s *pgxpool.Stat) int64 { return s.NewConnsCount() })
}
