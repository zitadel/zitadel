package metrics

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"go.opentelemetry.io/otel/api/metric"
	"net/http"
)

const (
	ActiveSessionCounter            = "zitadel.active_session_counter"
	ActiveSessionCounterDescription = "Active session counter"
	SpoolerDivCounter               = "zitadel.spooler_div_nanoseconds"
	SpoolerDivCounterDescription    = "Spooler div from last successful run to now in nanoseconds"
	Database                        = "database"
	ViewName                        = "view_name"
)

type Metrics interface {
	GetExporter() http.Handler
	GetMetricsProvider() metric.MeterProvider
	RegisterCounter(name, description string) error
	AddCount(ctx context.Context, name string, value int64, labels map[string]interface{}) error
	RegisterUpDownSumObserver(name, description string, callbackFunc metric.Int64ObserverFunc) error
}

type Config interface {
	NewMetrics() error
}

var M Metrics

func GetMetricsProvider(name string) metric.MeterProvider {
	if M == nil {
		return nil
	}
	return M.GetMetricsProvider()
}

func RegisterCounter(name, description string) error {
	if M == nil {
		return errors.ThrowPreconditionFailed(nil, "METER-3m9si", "No Meter implemented")
	}
	return M.RegisterCounter(name, description)
}

func AddCount(ctx context.Context, name string, value int64, labels map[string]interface{}) error {
	if M == nil {
		return errors.ThrowPreconditionFailed(nil, "METER-3m9si", "No Meter implemented")
	}
	return M.AddCount(ctx, name, value, labels)
}
