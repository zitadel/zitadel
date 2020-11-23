package otel

import (
	"context"
	"fmt"
	"github.com/caos/zitadel/internal/telemetry/metrics"
	"go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/label"
)

type Meter struct {
	Meter         metric.Meter
	Counters      map[string]metric.Int64Counter
	UpDownCounter map[string]metric.Int64UpDownCounter
}

func NewMeter(name string) metrics.Meter {
	return &Meter{
		Meter: global.MeterProvider().Meter(
			name,
			metric.WithInstrumentationVersion(contrib.SemVersion()),
		),
	}
}

func (m *Meter) GetMetricsProvider() metric.MeterProvider {
	return global.MeterProvider()
}

func (m *Meter) RegisterCounter(name string) error {
	counter, err := m.Meter.NewInt64Counter(name)
	if err != nil {
		return err
	}
	m.Counters[name] = counter
	return nil
}

func (m *Meter) AddCount(ctx context.Context, name string, value int64, labels map[string]interface{}) error {
	m.Counters[name].Add(ctx, value, mapToKeyValue(labels)...)
	return nil
}

func (m *Meter) PrintCounters() {
	fmt.Print(m.Counters)
}

func mapToKeyValue(labels map[string]interface{}) []label.KeyValue {
	keyValues := make([]label.KeyValue, 0)
	for key, value := range labels {
		keyValues = append(keyValues, label.Any(key, value))
	}
	return keyValues
}
