package otel

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/metrics"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/label"
	"net/http"
	"sync"
)

type Metrics struct {
	Exporter          *prometheus.Exporter
	Meter             metric.Meter
	Counters          sync.Map
	UpDownSumObserver sync.Map
	ValueObservers    sync.Map
}

func NewMetrics(meterName string) (metrics.Metrics, error) {
	exporter, err := prometheus.NewExportPipeline(
		prometheus.Config{},
	)
	if err != nil {
		return &Metrics{}, err
	}
	return &Metrics{
		Exporter: exporter,
		Meter:    exporter.MeterProvider().Meter(meterName),
	}, nil
}

func (m *Metrics) GetExporter() http.Handler {
	return m.Exporter
}

func (m *Metrics) GetMetricsProvider() metric.MeterProvider {
	return m.Exporter.MeterProvider()
}

func (m *Metrics) RegisterCounter(name, description string) error {
	if _, exists := m.Counters.Load(name); exists {
		return nil
	}
	counter := metric.Must(m.Meter).NewInt64Counter(name, metric.WithDescription(description))
	m.Counters.Store(name, counter)
	return nil
}

func (m *Metrics) AddCount(ctx context.Context, name string, value int64, labels map[string]interface{}) error {
	counter, exists := m.Counters.Load(name)
	if !exists {
		return caos_errs.ThrowNotFound(nil, "METER-4u8fs", "Errors.Metrics.Counter.NotFound")
	}
	counter.(metric.Int64Counter).Add(ctx, value, MapToKeyValue(labels)...)
	return nil
}

func (m *Metrics) RegisterUpDownSumObserver(name, description string, callbackFunc metric.Int64ObserverFunc) error {
	if _, exists := m.UpDownSumObserver.Load(name); exists {
		return nil
	}
	sumObserver := metric.Must(m.Meter).NewInt64UpDownSumObserver(
		name, callbackFunc, metric.WithDescription(description))

	m.UpDownSumObserver.Store(name, sumObserver)
	return nil
}

func (m *Metrics) RegisterValueObserver(name, description string, callbackFunc metric.Int64ObserverFunc) error {
	if _, exists := m.UpDownSumObserver.Load(name); exists {
		return nil
	}
	sumObserver := metric.Must(m.Meter).NewInt64ValueObserver(
		name, callbackFunc, metric.WithDescription(description))

	m.UpDownSumObserver.Store(name, sumObserver)
	return nil
}

func MapToKeyValue(labels map[string]interface{}) []label.KeyValue {
	if labels == nil {
		return nil
	}
	keyValues := make([]label.KeyValue, 0, len(labels))
	for key, value := range labels {
		keyValues = append(keyValues, label.Any(key, value))
	}
	return keyValues
}
