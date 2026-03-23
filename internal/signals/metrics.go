package signals

import (
	"context"

	"go.opentelemetry.io/otel/attribute"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/metrics"
)

const (
	streamLabel = "stream"

	signalsIngested            = "identity_signals.ingested_total"
	signalsIngestedDescription = "Total number of signals ingested into DuckLake"

	signalsDropped            = "identity_signals.dropped_total"
	signalsDroppedDescription = "Total number of signals dropped due to full channel"

	signalsBatchWriteDuration            = "identity_signals.batch_write_duration_seconds"
	signalsBatchWriteDurationDescription = "Time taken to write a signal batch to DuckLake"

	signalsCompactionDuration            = "identity_signals.compaction_duration_seconds"
	signalsCompactionDurationDescription = "Time taken to compact Parquet files"

	signalsPruned            = "identity_signals.retention_pruned_total"
	signalsPrunedDescription = "Total number of signals pruned by the retention worker"
)

// SignalMetrics provides OTEL metric instrumentation for the signals subsystem.
type SignalMetrics struct {
	provider metrics.Metrics
}

// NewSignalMetrics registers all signal-related OTEL metrics.
func NewSignalMetrics(ctx context.Context) *SignalMetrics {
	m := &SignalMetrics{provider: metrics.GlobalMeter()}

	err := m.provider.RegisterCounter(signalsIngested, signalsIngestedDescription)
	logging.OnError(ctx, err).Error("failed to register signals ingested counter")

	err = m.provider.RegisterCounter(signalsDropped, signalsDroppedDescription)
	logging.OnError(ctx, err).Error("failed to register signals dropped counter")

	err = m.provider.RegisterCounter(signalsPruned, signalsPrunedDescription)
	logging.OnError(ctx, err).Error("failed to register signals pruned counter")

	err = m.provider.RegisterHistogram(
		signalsBatchWriteDuration,
		signalsBatchWriteDurationDescription,
		"s",
		[]float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 10},
	)
	logging.OnError(ctx, err).Error("failed to register signals batch write histogram")

	err = m.provider.RegisterHistogram(
		signalsCompactionDuration,
		signalsCompactionDurationDescription,
		"s",
		[]float64{0.1, 0.5, 1, 5, 10, 30, 60, 120, 300},
	)
	logging.OnError(ctx, err).Error("failed to register signals compaction histogram")

	return m
}

// RecordIngested increments the ingested counter by count for the given stream.
func (m *SignalMetrics) RecordIngested(ctx context.Context, stream string, count int64) {
	if m == nil {
		return
	}
	err := m.provider.AddCount(ctx, signalsIngested, count, map[string]attribute.Value{
		streamLabel: attribute.StringValue(stream),
	})
	logging.OnError(ctx, err).Error("failed to record signals ingested")
}

// RecordDropped increments the dropped counter.
func (m *SignalMetrics) RecordDropped(ctx context.Context, count int64) {
	if m == nil {
		return
	}
	err := m.provider.AddCount(ctx, signalsDropped, count, nil)
	logging.OnError(ctx, err).Error("failed to record signals dropped")
}

// RecordBatchWriteDuration records the time taken to write a batch.
func (m *SignalMetrics) RecordBatchWriteDuration(ctx context.Context, durationSeconds float64) {
	if m == nil {
		return
	}
	err := m.provider.AddHistogramMeasurement(ctx, signalsBatchWriteDuration, durationSeconds, nil)
	logging.OnError(ctx, err).Error("failed to record batch write duration")
}

// RecordCompactionDuration records the time taken for a compaction cycle.
func (m *SignalMetrics) RecordCompactionDuration(ctx context.Context, durationSeconds float64) {
	if m == nil {
		return
	}
	err := m.provider.AddHistogramMeasurement(ctx, signalsCompactionDuration, durationSeconds, nil)
	logging.OnError(ctx, err).Error("failed to record compaction duration")
}

// RecordPruned increments the pruned counter by count for the given stream.
func (m *SignalMetrics) RecordPruned(ctx context.Context, stream string, count int64) {
	if m == nil {
		return
	}
	err := m.provider.AddCount(ctx, signalsPruned, count, map[string]attribute.Value{
		streamLabel: attribute.StringValue(stream),
	})
	logging.OnError(ctx, err).Error("failed to record signals pruned")
}
