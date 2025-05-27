package metrics

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type NoopMetrics struct{}

var _ Metrics = new(NoopMetrics)

func (n *NoopMetrics) GetExporter() http.Handler {
	return nil
}

func (n *NoopMetrics) GetMetricsProvider() metric.MeterProvider {
	return nil
}

func (n *NoopMetrics) RegisterCounter(name, description string) error {
	return nil
}

func (n *NoopMetrics) AddCount(ctx context.Context, name string, value int64, labels map[string]attribute.Value) error {
	return nil
}

func (n *NoopMetrics) AddHistogramMeasurement(ctx context.Context, name string, value float64, labels map[string]attribute.Value) error {
	return nil
}

func (n *NoopMetrics) RegisterUpDownSumObserver(name, description string, callbackFunc metric.Int64Callback) error {
	return nil
}

func (n *NoopMetrics) RegisterValueObserver(name, description string, callbackFunc metric.Int64Callback) error {
	return nil
}

func (n *NoopMetrics) RegisterHistogram(name, description, unit string, buckets []float64) error {
	return nil
}
