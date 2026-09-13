package database

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
)

func Test_registerPoolMetrics(t *testing.T) {
	cfg, err := pgxpool.ParseConfig("postgres://user:pass@localhost:5432/db")
	require.NoError(t, err)

	// pgxpool.NewWithConfig does not establish any connection, so this works
	// without a reachable database.
	pool, err := pgxpool.NewWithConfig(t.Context(), cfg)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	m := instrumentation.NewMeter(t.Name())
	registerPoolMetrics(t.Context(), m, pool)

	gauges := []string{
		PoolAcquiredConns,
		PoolConstructingConns,
		PoolIdleConns,
		PoolMaxConns,
		PoolTotalConns,
	}
	for _, name := range gauges {
		_, exists := m.UpDownSumObserver.Load(name)
		assert.Truef(t, exists, "expected gauge %q to be registered", name)
	}

	counters := []string{
		PoolAcquireCount,
		PoolAcquireDurationMs,
		PoolCanceledAcquireCount,
		PoolEmptyAcquireCount,
		PoolMaxIdleDestroyCount,
		PoolMaxLifetimeDestroyCount,
		PoolNewConnsCount,
	}
	for _, name := range counters {
		_, exists := m.CounterObservers.Load(name)
		assert.Truef(t, exists, "expected counter %q to be registered", name)
	}
}
