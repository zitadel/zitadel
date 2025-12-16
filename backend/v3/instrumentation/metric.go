package instrumentation

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	google_metric "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdk_metric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type MetricConfig struct {
	Exporter ExporterConfig
}

type Meter struct {
	Meter             metric.Meter
	Counters          sync.Map
	UpDownSumObserver sync.Map
	ValueObservers    sync.Map
	Histograms        sync.Map
}

// TODO: remove for v5 release
type LegacyMetricConfig struct {
	Type string
}

func (c *MetricConfig) SetLegacyConfig(lc *LegacyMetricConfig) {
	typ := c.Exporter.Type
	if lc == nil || typ.isNone() {
		return
	}
	if lc.Type == "otel" {
		c.Exporter.Type = ExporterTypePrometheus
	}
}

var hasPrometheusExporter bool

func (m *Meter) GetExporter() http.Handler {
	if hasPrometheusExporter {
		return promhttp.Handler()
	}
	return http.NotFoundHandler()
}

func (m *Meter) RegisterCounter(name, description string) error {
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

func (m *Meter) AddCount(ctx context.Context, name string, value int64, labels map[string]attribute.Value) error {
	counter, exists := m.Counters.Load(name)
	if !exists {
		return zerrors.ThrowNotFound(nil, "METER-4u8fs", "Errors.Metrics.Counter.NotFound")
	}
	counter.(metric.Int64Counter).Add(ctx, value, MapToAddOption(labels)...)
	return nil
}

func (m *Meter) AddHistogramMeasurement(ctx context.Context, name string, value float64, labels map[string]attribute.Value) error {
	histogram, exists := m.Histograms.Load(name)
	if !exists {
		return zerrors.ThrowNotFound(nil, "METER-5wwb1", "Errors.Metrics.Histogram.NotFound")
	}
	histogram.(metric.Float64Histogram).Record(ctx, value, MapToRecordOption(labels)...)
	return nil
}

func (m *Meter) RegisterHistogram(name, description, unit string, buckets []float64) error {
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

func (m *Meter) RegisterUpDownSumObserver(name, description string, callbackFunc metric.Int64Callback) error {
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

func (m *Meter) RegisterValueObserver(name, description string, callbackFunc metric.Int64Callback) error {
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

func newMeterProvider(ctx context.Context, cfg MetricConfig, resource *resource.Resource) (_ *sdk_metric.MeterProvider, err error) {
	var readerOption sdk_metric.Option
	switch cfg.Exporter.Type {
	case ExporterTypeUnspecified, ExporterTypeNone:
		// no reader option
	case ExporterTypeStdOut, ExporterTypeStdErr:
		readerOption, err = metricStdOutOption(cfg.Exporter)
	case ExporterTypeGRPC:
		readerOption, err = metricGrpcOption(ctx, cfg.Exporter)
	case ExporterTypeHTTP:
		readerOption, err = metricHttpOption(ctx, cfg.Exporter)
	case ExporterTypeGoogle:
		readerOption, err = metricGoogleOption(cfg.Exporter)
	case ExporterTypePrometheus:
		readerOption, err = metricPrometheusOption()
	default:
		err = errExporterType(cfg.Exporter.Type, "metrics")
	}
	if err != nil {
		return nil, fmt.Errorf("meter reader: %w", err)
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

	opts := []sdk_metric.Option{
		sdk_metric.WithResource(resource),
		sdk_metric.WithView(view),
	}
	if readerOption != nil {
		opts = append(opts, readerOption)
	}
	meterProvider := sdk_metric.NewMeterProvider(opts...)
	return meterProvider, nil
}

func metricStdOutOption(cfg ExporterConfig) (sdk_metric.Option, error) {
	options := []stdoutmetric.Option{
		stdoutmetric.WithPrettyPrint(),
	}
	if cfg.Type == ExporterTypeStdErr {
		options = append(options, stdoutmetric.WithWriter(os.Stderr))
	}
	exporter, err := stdoutmetric.New(options...)
	if err != nil {
		return nil, err
	}
	return sdk_metric.WithReader(
		sdk_metric.NewPeriodicReader(
			exporter,
			sdk_metric.WithInterval(cfg.BatchDuration),
		),
	), nil
}

func metricGrpcOption(ctx context.Context, cfg ExporterConfig) (sdk_metric.Option, error) {
	var grpcOpts []otlpmetricgrpc.Option
	if cfg.Endpoint != "" {
		grpcOpts = append(grpcOpts, otlpmetricgrpc.WithEndpoint(cfg.Endpoint))
	}
	if cfg.Insecure {
		grpcOpts = append(grpcOpts, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(ctx, grpcOpts...)
	if err != nil {
		return nil, err
	}
	return sdk_metric.WithReader(
		sdk_metric.NewPeriodicReader(
			exporter,
			sdk_metric.WithInterval(cfg.BatchDuration),
		),
	), nil
}

func metricHttpOption(ctx context.Context, cfg ExporterConfig) (sdk_metric.Option, error) {
	var httpOpts []otlpmetrichttp.Option
	if cfg.Endpoint != "" {
		httpOpts = append(httpOpts, otlpmetrichttp.WithEndpoint(cfg.Endpoint))
	}
	if cfg.Insecure {
		httpOpts = append(httpOpts, otlpmetrichttp.WithInsecure())
	}

	exporter, err := otlpmetrichttp.New(ctx, httpOpts...)
	if err != nil {
		return nil, err
	}
	return sdk_metric.WithReader(
		sdk_metric.NewPeriodicReader(
			exporter,
			sdk_metric.WithInterval(cfg.BatchDuration),
		),
	), nil
}

func metricGoogleOption(cfg ExporterConfig) (sdk_metric.Option, error) {
	exporter, err := google_metric.New(
		google_metric.WithProjectID(cfg.GoogleProjectID),
	)
	if err != nil {
		return nil, err
	}
	return sdk_metric.WithReader(
		sdk_metric.NewPeriodicReader(
			exporter,
			sdk_metric.WithInterval(cfg.BatchDuration),
		),
	), nil
}

func metricPrometheusOption() (sdk_metric.Option, error) {
	prom, err := prometheus.New(prometheus.WithoutScopeInfo())
	if err != nil {
		return nil, err
	}
	hasPrometheusExporter = true
	return sdk_metric.WithReader(prom), nil
}
