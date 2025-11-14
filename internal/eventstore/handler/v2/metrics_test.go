package handler

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

func TestNewProjectionMetrics(t *testing.T) {
	mockMetrics := metrics.NewMockMetrics()
	metrics.M = mockMetrics

	metrics := NewProjectionMetrics()
	require.NotNil(t, metrics)
	assert.NotNil(t, metrics.provider)
}

func TestProjectionMetrics_ProjectionUpdateTiming(t *testing.T) {

	mockMetrics := metrics.NewMockMetrics()
	metrics.M = mockMetrics
	projectionMetrics := NewProjectionMetrics()

	ctx := context.Background()
	projection := "test_projection"
	duration := 0.5

	projectionMetrics.ProjectionUpdateTiming(ctx, projection, duration)

	values := mockMetrics.GetHistogramValues(ProjectionHandleTimerMetric)
	require.Len(t, values, 1)
	assert.Equal(t, duration, values[0])

	labels := mockMetrics.GetHistogramLabels(ProjectionHandleTimerMetric)
	require.Len(t, labels, 1)
	assert.Equal(t, projection, labels[0][ProjectionLabel].AsString())
}

func TestProjectionMetrics_ProjectionEventsProcessed(t *testing.T) {

	mockMetrics := metrics.NewMockMetrics()
	metrics.M = mockMetrics
	projectionMetrics := NewProjectionMetrics()

	ctx := context.Background()
	projection := "test_projection"
	count := int64(5)
	success := true

	projectionMetrics.ProjectionEventsProcessed(ctx, projection, count, success)

	value := mockMetrics.GetCounterValue(ProjectionEventsProcessed)
	assert.Equal(t, count, value)

	labels := mockMetrics.GetCounterLabels(ProjectionEventsProcessed)
	require.Len(t, labels, 1)
	assert.Equal(t, projection, labels[0][ProjectionLabel].AsString())
	assert.Equal(t, success, labels[0][SuccessLabel].AsBool())
}

func TestProjectionMetrics_ProjectionStateLatency(t *testing.T) {

	mockMetrics := metrics.NewMockMetrics()
	metrics.M = mockMetrics
	projectionMetrics := NewProjectionMetrics()

	ctx := context.Background()
	projection := "test_projection"
	latency := 10.0

	projectionMetrics.ProjectionStateLatency(ctx, projection, latency)

	values := mockMetrics.GetHistogramValues(ProjectionStateLatencyMetric)
	require.Len(t, values, 1)
	assert.Equal(t, latency, values[0])

	labels := mockMetrics.GetHistogramLabels(ProjectionStateLatencyMetric)
	require.Len(t, labels, 1)
	assert.Equal(t, projection, labels[0][ProjectionLabel].AsString())
}

func TestProjectionMetrics_Integration(t *testing.T) {

	mockMetrics := metrics.NewMockMetrics()
	metrics.M = mockMetrics
	projectionMetrics := NewProjectionMetrics()

	ctx := context.Background()
	projection := "test_projection"

	start := time.Now()

	projectionMetrics.ProjectionEventsProcessed(ctx, projection, 3, true)
	projectionMetrics.ProjectionEventsProcessed(ctx, projection, 1, false)

	duration := time.Since(start).Seconds()
	projectionMetrics.ProjectionUpdateTiming(ctx, projection, duration)

	latency := 5.0
	projectionMetrics.ProjectionStateLatency(ctx, projection, latency)

	value := mockMetrics.GetCounterValue(ProjectionEventsProcessed)
	assert.Equal(t, int64(4), value)

	timingValues := mockMetrics.GetHistogramValues(ProjectionHandleTimerMetric)
	require.Len(t, timingValues, 1)
	assert.Equal(t, duration, timingValues[0])

	latencyValues := mockMetrics.GetHistogramValues(ProjectionStateLatencyMetric)
	require.Len(t, latencyValues, 1)
	assert.Equal(t, latency, latencyValues[0])

	eventsLabels := mockMetrics.GetCounterLabels(ProjectionEventsProcessed)
	require.Len(t, eventsLabels, 2)
	assert.Equal(t, projection, eventsLabels[0][ProjectionLabel].AsString())
	assert.Equal(t, true, eventsLabels[0][SuccessLabel].AsBool())
	assert.Equal(t, projection, eventsLabels[1][ProjectionLabel].AsString())
	assert.Equal(t, false, eventsLabels[1][SuccessLabel].AsBool())

	timingLabels := mockMetrics.GetHistogramLabels(ProjectionHandleTimerMetric)
	require.Len(t, timingLabels, 1)
	assert.Equal(t, projection, timingLabels[0][ProjectionLabel].AsString())

	latencyLabels := mockMetrics.GetHistogramLabels(ProjectionStateLatencyMetric)
	require.Len(t, latencyLabels, 1)
	assert.Equal(t, projection, latencyLabels[0][ProjectionLabel].AsString())
}
