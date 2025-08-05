package handler

import (
	"context"

	"github.com/zitadel/logging"
	"go.opentelemetry.io/otel/attribute"

	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

const (
	ProjectionLabel = "projection"
	SuccessLabel    = "success"

	ProjectionEventsProcessed    = "projection_events_processed"
	ProjectionHandleTimerMetric  = "projection_handle_timer"
	ProjectionStateLatencyMetric = "projection_state_latency"
)

type ProjectionMetrics struct {
	provider metrics.Metrics
}

func NewProjectionMetrics() *ProjectionMetrics {
	projectionMetrics := &ProjectionMetrics{provider: metrics.M}

	err := projectionMetrics.provider.RegisterCounter(
		ProjectionEventsProcessed,
		"Number of events reduced to process projection updates",
	)
	logging.OnError(err).Error("failed to register projection events processed counter")
	err = projectionMetrics.provider.RegisterHistogram(
		ProjectionHandleTimerMetric,
		"Time taken to process a projection update",
		"s",
		[]float64{0.005, 0.01, 0.05, 0.1, 1, 5, 10, 30, 60, 120},
	)
	logging.OnError(err).Error("failed to register projection handle timer metric")
	err = projectionMetrics.provider.RegisterHistogram(
		ProjectionStateLatencyMetric,
		"When finishing processing a batch of events, this track the age of the last events seen from current time",
		"s",
		[]float64{0.1, 0.5, 1, 5, 10, 30, 60, 300, 600, 1800},
	)
	logging.OnError(err).Error("failed to register projection state latency metric")
	return projectionMetrics
}

func (m *ProjectionMetrics) ProjectionUpdateTiming(ctx context.Context, projection string, duration float64) {
	err := m.provider.AddHistogramMeasurement(ctx, ProjectionHandleTimerMetric, duration, map[string]attribute.Value{
		ProjectionLabel: attribute.StringValue(projection),
	})
	logging.OnError(err).Error("failed to add projection trigger timing")
}

func (m *ProjectionMetrics) ProjectionEventsProcessed(ctx context.Context, projection string, count int64, success bool) {
	err := m.provider.AddCount(ctx, ProjectionEventsProcessed, count, map[string]attribute.Value{
		ProjectionLabel: attribute.StringValue(projection),
		SuccessLabel:    attribute.BoolValue(success),
	})
	logging.OnError(err).Error("failed to add projection events processed metric")
}

func (m *ProjectionMetrics) ProjectionStateLatency(ctx context.Context, projection string, latency float64) {
	err := m.provider.AddHistogramMeasurement(ctx, ProjectionStateLatencyMetric, latency, map[string]attribute.Value{
		ProjectionLabel: attribute.StringValue(projection),
	})
	logging.OnError(err).Error("failed to add projection state latency metric")
}
