package otel

import (
	"context"
	"github.com/caos/zitadel/internal/metrics"
	"go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/label"
)

type Meter struct {
	meter    metric.Meter
	counters map[string]metric.Int64Counter
}

func NewMeter(name string) metrics.Meter {
	return &Meter{
		meter: global.MeterProvider().Meter(
			name,
			metric.WithInstrumentationVersion(contrib.SemVersion()),
		),
	}
}

func (m *Meter) RegisterCounter(name string) error {
	counter, err := m.meter.NewInt64Counter(name)
	if err != nil {
		return err
	}
	m.counters[name] = counter
	return nil
}

func (m *Meter) AddCount(ctx context.Context, name string, value int64, labels map[string]interface{}) error {
	m.counters[name].Add(ctx, value, mapToKeyValue(labels)...)
	return nil
}

func mapToKeyValue(labels map[string]interface{}) []label.KeyValue {
	return nil
}
