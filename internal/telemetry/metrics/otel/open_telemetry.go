package otel

import (
	"context"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
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
	Histograms        sync.Map
}

func NewMetrics(meterName string) (metrics.Metrics, error) {
	resource, err := otel_resource.ResourceWithService("ZITADEL")
	if err != nil {
		return nil, err
	}
	exporter, err := prometheus.New(prometheus.WithoutScopeInfo())
	if err != nil {
		return &Metrics{}, err
	}
	// create a view to filter out unwanted attributes
	view := sdk_metric.NewView(
		sdk_metric.Instrument{
			Scope: instrumentation.Scope{Name: otelhttp.ScopeName},
		},
		sdk_metric.Stream{
			AttributeFilter: attribute.NewAllowKeysFilter("http.method", "http.status_code", "http.target"),
		},
	)
	meterProvider := sdk_metric.NewMeterProvider(
		sdk_metric.WithReader(exporter),
		sdk_metric.WithResource(resource),
		sdk_metric.WithView(view),
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

func (m *Metrics) AddHistogramMeasurement(ctx context.Context, name string, value float64, labels map[string]attribute.Value) error {
	histogram, exists := m.Histograms.Load(name)
	if !exists {
		return zerrors.ThrowNotFound(nil, "METER-5wwb1", "Errors.Metrics.Histogram.NotFound")
	}
	histogram.(metric.Float64Histogram).Record(ctx, value, MapToRecordOption(labels)...)
	return nil
}

func (m *Metrics) RegisterHistogram(name, description, unit string, buckets []float64) error {
	if _, exists := m.Histograms.Load(name); exists {
		return nil
	}

	histogram, err := m.Meter.Float64Histogram(name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
		metric.WithExplicitBucketBoundaries(buckets...),
	)
	if err != nil {
		return err
	}

	m.Histograms.Store(name, histogram)
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
	return []metric.AddOption{metric.WithAttributes(labelsToAttributes(labels)...)}
}

func MapToRecordOption(labels map[string]attribute.Value) []metric.RecordOption {
	return []metric.RecordOption{metric.WithAttributes(labelsToAttributes(labels)...)}
}

func labelsToAttributes(labels map[string]attribute.Value) []attribute.KeyValue {
	if labels == nil {
		return nil
	}
	attributes := make([]attribute.KeyValue, 0, len(labels))
	for key, value := range labels {
		attributes = append(attributes, attribute.KeyValue{
			Key:   attribute.Key(key),
			Value: value,
		})
	}
	return attributes
}
