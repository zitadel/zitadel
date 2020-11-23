package metrics

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"go.opentelemetry.io/otel/api/metric"
)

type Meter interface {
	GetMetricsProvider() metric.MeterProvider
	RegisterCounter(name string) error
	AddCount(ctx context.Context, name string, value int64, labels map[string]interface{}) error
	PrintCounters()
}

type Config interface {
	NewMetrics() error
}

var M Meter

func GetMetricsProvider(name string) metric.MeterProvider {
	if M == nil {
		return nil
	}
	return M.GetMetricsProvider()
}

func RegisterCounter(name string) error {
	if M == nil {
		return errors.ThrowPreconditionFailed(nil, "METER-3m9si", "No Meter implemented")
	}
	return M.RegisterCounter(name)
}

func AddCount(ctx context.Context, name string, value int64, labels map[string]interface{}) error {
	if M == nil {
		return errors.ThrowPreconditionFailed(nil, "METER-3m9si", "No Meter implemented")
	}
	return M.AddCount(ctx, name, value, labels)
}

func PrintCounters() {
	if M == nil {
		return
	}
	M.PrintCounters()
}
