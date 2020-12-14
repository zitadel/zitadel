package metrics

import (
	"context"
	"go.opentelemetry.io/otel/api/metric"
	"net/http"
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
	AddCount(ctx context.Context, name string, value int64, labels map[string]interface{}) error
	RegisterUpDownSumObserver(name, description string, callbackFunc metric.Int64ObserverFunc) error
	RegisterValueObserver(name, description string, callbackFunc metric.Int64ObserverFunc) error
}

type Config interface {
	NewMetrics() error
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

func AddCount(ctx context.Context, name string, value int64, labels map[string]interface{}) error {
	if M == nil {
		return nil
	}
	return M.AddCount(ctx, name, value, labels)
}

func RegisterUpDownSumObserver(name, description string, callbackFunc metric.Int64ObserverFunc) error {
	if M == nil {
		return nil
	}
	return M.RegisterUpDownSumObserver(name, description, callbackFunc)
}

func RegisterValueObserver(name, description string, callbackFunc metric.Int64ObserverFunc) error {
	if M == nil {
		return nil
	}
	return M.RegisterValueObserver(name, description, callbackFunc)
}
