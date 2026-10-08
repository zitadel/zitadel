package database

import (
	"slices"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sdk_metric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

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

	reader := sdk_metric.NewManualReader()
	provider := sdk_metric.NewMeterProvider(sdk_metric.WithReader(reader))
	m := &instrumentation.Meter{Meter: provider.Meter(t.Name())}

	registerPoolMetrics(t.Context(), m, pool)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))

	gauges := map[string]int64{
		PoolAcquiredConns:     0,
		PoolConstructingConns: 0,
		PoolIdleConns:         0,
		PoolMaxConns:          int64(cfg.MaxConns),
		PoolTotalConns:        0,
	}
	for name, want := range gauges {
		t.Run(name, func(t *testing.T) {
			got := findMetric(t, rm, name)
			gauge, ok := got.Data.(metricdata.Gauge[int64])
			require.True(t, ok, "expected a gauge, got %T", got.Data)
			require.Len(t, gauge.DataPoints, 1)
			assert.Equal(t, want, gauge.DataPoints[0].Value)
		})
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
		t.Run(name, func(t *testing.T) {
			got := findMetric(t, rm, name)
			sum, ok := got.Data.(metricdata.Sum[int64])
			require.True(t, ok, "expected an observable counter (Sum[int64]), got %T", got.Data)
			assert.True(t, sum.IsMonotonic)
			require.Len(t, sum.DataPoints, 1)
			// A freshly created, unused pool has not acquired or destroyed any connections yet.
			assert.Equal(t, int64(0), sum.DataPoints[0].Value)
		})
	}
}

// findMetric locates a metric by name across all scopes collected from a reader.
func findMetric(t *testing.T, rm metricdata.ResourceMetrics, name string) metricdata.Metrics {
	t.Helper()
	for _, sm := range rm.ScopeMetrics {
		if i := slices.IndexFunc(sm.Metrics, func(m metricdata.Metrics) bool { return m.Name == name }); i >= 0 {
			return sm.Metrics[i]
		}
	}
	t.Fatalf("metric %q not found in collected data", name)
	return metricdata.Metrics{}
}
