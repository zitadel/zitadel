package otel

import (
	"context"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdk_metric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/zitadel/zitadel/internal/telemetry/metrics"
	otel_resource "github.com/zitadel/zitadel/internal/telemetry/otel"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Metrics struct {
	Provider          metric.MeterProvider
	Meter             metric.Meter
	Counters          sync.Map
	UpDownSumObserver sync.Map
	ValueObservers    sync.Map
}

func NewMetrics(meterName string) (metrics.Metrics, error) {
	resource, err := otel_resource.ResourceWithService()
	if err != nil {
		return nil, err
	}
	exporter, err := prometheus.New()
	if err != nil {
		return &Metrics{}, err
	}
	meterProvider := sdk_metric.NewMeterProvider(
		sdk_metric.WithReader(exporter),
		sdk_metric.WithResource(resource),
	)
	return &Metrics{
		Provider: meterProvider,
		Meter:    meterProvider.Meter(meterName),
	}, nil
}

func (m *Metrics) GetExporter() http.Handler {
	return promhttp.Handler()
}

func (m *Metrics) GetMetricsProvider() metric.MeterProvider {
	return m.Provider
}

func (m *Metrics) RegisterCounter(name, description string) error {
	if _, exists := m.Counters.Load(name); exists {
		return nil
	}
	counter, err := m.Meter.Int64Counter(name, metric.WithDescription(description))
	if err != nil {
		return err
	}
	m.Counters.Store(name, counter)
	return nil
}

func (m *Metrics) AddCount(ctx context.Context, name string, value int64, labels map[string]attribute.Value) error {
	counter, exists := m.Counters.Load(name)
	if !exists {
		return zerrors.ThrowNotFound(nil, "METER-4u8fs", "Errors.Metrics.Counter.NotFound")
	}
	counter.(metric.Int64Counter).Add(ctx, value, MapToAddOption(labels)...)
	return nil
}

func (m *Metrics) RegisterUpDownSumObserver(name, description string, callbackFunc metric.Int64Callback) error {
	if _, exists := m.UpDownSumObserver.Load(name); exists {
		return nil
	}

	counter, err := m.Meter.Int64ObservableUpDownCounter(name, metric.WithInt64Callback(callbackFunc), metric.WithDescription(description))
	if err != nil {
		return err
	}

	m.UpDownSumObserver.Store(name, counter)
	return nil
}

func (m *Metrics) RegisterValueObserver(name, description string, callbackFunc metric.Int64Callback) error {
	if _, exists := m.UpDownSumObserver.Load(name); exists {
		return nil
	}

	gauge, err := m.Meter.Int64ObservableGauge(name, metric.WithInt64Callback(callbackFunc), metric.WithDescription(description))
	if err != nil {
		return err
	}

	m.UpDownSumObserver.Store(name, gauge)
	return nil
}

func MapToAddOption(labels map[string]attribute.Value) []metric.AddOption {
	if labels == nil {
		return nil
	}
	keyValues := make([]attribute.KeyValue, 0, len(labels))
	for key, value := range labels {
		keyValues = append(keyValues, attribute.KeyValue{
			Key:   attribute.Key(key),
			Value: value,
		})
	}
	return []metric.AddOption{metric.WithAttributes(keyValues...)}
}
