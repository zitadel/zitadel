package instrumentation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"github.com/zitadel/zitadel/cmd/build"
)

//go:generate enumer -type=ExporterType -trimprefix=ExporterType -text -linecomment
type ExporterType int

// ExporterType defines the type of exporter to use.
// - metrics supports all types
// - tracing supports all types except Prometheus
// - logging supports all types except Google and Prometheus
// - profiling supports only Google and None
const (
	// Empty line comment sets empty string of unspecified value
	ExporterTypeUnspecified ExporterType = iota //
	ExporterTypeNone
	ExporterTypeStdOut
	ExporterTypeStdErr
	ExporterTypeGRPC
	ExporterTypeHTTP
	ExporterTypeGoogle
	ExporterTypePrometheus
)

// isNone returns true if the ExporterType is either Unspecified or None.
func (e ExporterType) isNone() bool {
	return e == ExporterTypeUnspecified || e == ExporterTypeNone
}

type ExporterConfig struct {
	Type            ExporterType
	Endpoint        string
	Insecure        bool
	BatchDuration   time.Duration
	GoogleProjectID string // only for tracing, metrics and profiler
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context, cfg Config) (*log.LoggerProvider, ShutdownFunc, error) {
	var (
		shutdownFuncs []ShutdownFunc
		err           error
	)

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	resource, err := resourceWithService("ZITADEL")
	if err != nil {
		return nil, shutdown, err
	}

	// Set up trace provider.
	tracerProvider, err := newTracerProvider(ctx, cfg.Trace, resource)
	if err != nil {
		handleErr(err)
		return nil, shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider(ctx, cfg.Metric, resource)
	if err != nil {
		handleErr(err)
		return nil, shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)
	// Set up logger provider.
	loggerProvider, err := newLoggerProvider(ctx, cfg.Log.Exporter, resource)
	if err != nil {
		handleErr(err)
		return nil, shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return loggerProvider, shutdown, err
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func resourceWithService(serviceName string) (*resource.Resource, error) {
	attributes := []attribute.KeyValue{
		semconv.ServiceNameKey.String(serviceName),
	}
	if build.Version() != "" {
		attributes = append(attributes, semconv.ServiceVersionKey.String(build.Version()))
	}
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes("", attributes...),
	)
}

func errExporterType(typ ExporterType, instrument string) error {
	return fmt.Errorf("exporter type \"%v\" unsupported for %s", typ, instrument)
}
