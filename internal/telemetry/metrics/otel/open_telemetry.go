package otel

import (
	"context"
	"fmt"
	"github.com/caos/zitadel/internal/telemetry/metrics"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/label"
	"net/http"
)

type Metrics struct {
	Exporter      *prometheus.Exporter
	Meter         metric.Meter
	Counters      map[string]metric.Int64Counter
	UpDownCounter map[string]metric.Int64UpDownCounter
}

func NewMetrics() (metrics.Metrics, error) {
	exporter, err := prometheus.NewExportPipeline(
		prometheus.Config{},
	)
	if err != nil {
		return &Metrics{}, err
	}
	return &Metrics{
		Exporter:      exporter,
		Meter:         exporter.MeterProvider().Meter("hodor"),
		Counters:      make(map[string]metric.Int64Counter),
		UpDownCounter: make(map[string]metric.Int64UpDownCounter),
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
	m.Counters[name].Add(ctx, value, mapToKeyValue(labels)...)
	fmt.Printf("Counters: %v", m.Counters)
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
