package metrics

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	ActiveSessionCounter            = "zitadel.active_session_counter"
	ActiveSessionCounterDescription = "Active session counter"
	SpoolerDivCounter               = "zitadel.spooler_div_milliseconds"
	SpoolerDivCounterDescription    = "Spooler div from last successful run to now in milliseconds"
	Database                        = "database"
	ViewName                        = "view_name"
)

type Metrics interface {
	GetExporter() http.Handler
	GetMetricsProvider() metric.MeterProvider
	RegisterCounter(name, description string) error
	AddCount(ctx context.Context, name string, value int64, labels map[string]attribute.Value) error
	AddHistogramMeasurement(ctx context.Context, name string, value float64, labels map[string]attribute.Value) error
	RegisterUpDownSumObserver(name, description string, callbackFunc metric.Int64Callback) error
	RegisterValueObserver(name, description string, callbackFunc metric.Int64Callback) error
	RegisterHistogram(name, description, unit string, buckets []float64) error
}

var M Metrics

func GetExporter() http.Handler {
	if M == nil {
		return nil
	}
	return M.GetExporter()
}

func GetMetricsProvider() metric.MeterProvider {
	if M == nil {
		return nil
	}
	return M.GetMetricsProvider()
}

func RegisterCounter(name, description string) error {
	if M == nil {
		return nil
	}
	return M.RegisterCounter(name, description)
}

func AddCount(ctx context.Context, name string, value int64, labels map[string]attribute.Value) error {
	if M == nil {
		return nil
	}
	return M.AddCount(ctx, name, value, labels)
}

func AddHistogramMeasurement(ctx context.Context, name string, value float64, labels map[string]attribute.Value) error {
	if M == nil {
		return nil
	}
	return M.AddHistogramMeasurement(ctx, name, value, labels)
}

func RegisterHistogram(name, description, unit string, buckets []float64) error {
	if M == nil {
		return nil
	}
	return M.RegisterHistogram(name, description, unit, buckets)
}

func RegisterUpDownSumObserver(name, description string, callbackFunc metric.Int64Callback) error {
	if M == nil {
		return nil
	}
	return M.RegisterUpDownSumObserver(name, description, callbackFunc)
}

func RegisterValueObserver(name, description string, callbackFunc metric.Int64Callback) error {
	if M == nil {
		return nil
	}
	return M.RegisterValueObserver(name, description, callbackFunc)
}
