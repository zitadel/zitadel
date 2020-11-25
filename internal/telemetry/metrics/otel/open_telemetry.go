package otel

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/metrics"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/label"
	"net/http"
)

type Metrics struct {
	Exporter          *prometheus.Exporter
	Meter             metric.Meter
	Counters          map[string]metric.Int64Counter
	UpDownSumObserver map[string]metric.Int64UpDownSumObserver
}

func NewMetrics() (metrics.Metrics, error) {
	exporter, err := prometheus.NewExportPipeline(
		prometheus.Config{},
	)
	if err != nil {
		return &Metrics{}, err
	}
	return &Metrics{
		Exporter:          exporter,
		Meter:             exporter.MeterProvider().Meter("hodor"),
		Counters:          make(map[string]metric.Int64Counter),
		UpDownSumObserver: make(map[string]metric.Int64UpDownSumObserver),
	}, nil
}

func (m *Metrics) GetExporter() http.Handler {
	return m.Exporter
}

func (m *Metrics) GetMetricsProvider() metric.MeterProvider {
	return m.Exporter.MeterProvider()
}

func (m *Metrics) RegisterCounter(name, description string) error {
	if _, exists := m.Counters[name]; exists {
		return nil
	}
	counter := metric.Must(m.Meter).NewInt64Counter(name, metric.WithDescription(description), metric.WithInstrumentationName("test"))

	m.Counters[name] = counter
	return nil
}

func (m *Metrics) AddCount(ctx context.Context, name string, value int64, labels map[string]interface{}) error {
	if _, exists := m.Counters[name]; !exists {
		return caos_errs.ThrowNotFound(nil, "METER-4u8fs", "Errors.Metrics.Counter.NotFound")
	}
	m.Counters[name].Add(ctx, value, mapToKeyValue(labels)...)
	return nil
}

func (m *Metrics) RegisterUpDownSumObserver(name, description string, callbackFunc metric.Int64ObserverFunc) error {
	//callbackFunc := func(_ context.Context, result metric.Int64ObserverResult) {
	//	result.Observe(data, mapToKeyValue(labels)...)
	//}
	if _, exists := m.UpDownSumObserver[name]; exists {
		return nil
	}
	sumObserver := metric.Must(m.Meter).NewInt64UpDownSumObserver(
		name, callbackFunc, metric.WithDescription(description), metric.WithInstrumentationName("test"))

	m.UpDownSumObserver[name] = sumObserver
	return nil
}

func mapToKeyValue(labels map[string]interface{}) []label.KeyValue {
	keyValues := make([]label.KeyValue, 0)
	if labels == nil {
		return keyValues
	}
	for key, value := range labels {
		keyValues = append(keyValues, label.Any(key, value))
	}
	return keyValues
}
