package instrumentation

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	google_detector "go.opentelemetry.io/contrib/detectors/gcp"
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
	ExporterTypeAuto
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

//go:generate enumer -type=DetectorType -trimprefix=DetectorType -text -linecomment
type DetectorType int

// DetectorType defines a resource detector that describes the platform ZITADEL
// runs on. Detected attributes are attached to every trace, metric and log we
// export. Detectors are opt-in: each one probes its platform on startup, and the
// attributes it finds identify the account and host we run on to whichever
// backend we export to.
const (
	// Empty line comment sets empty string of unspecified value
	DetectorTypeUnspecified DetectorType = iota //
	DetectorTypeGoogle
)

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

	resource, err := newResource(ctx, cfg.ServiceName, cfg.Detectors)
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

// newResource builds the resource that identifies this process on every signal
// it emits. Detectors are applied in order and the later ones win, so the
// service attributes we are explicitly configured with always take precedence
// over anything a platform detector found.
func newResource(ctx context.Context, serviceName string, detectorTypes []DetectorType) (*resource.Resource, error) {
	// OTEL_SERVICE_NAME takes priority over programmatic config per the
	// OpenTelemetry specification.
	if envName := os.Getenv("OTEL_SERVICE_NAME"); envName != "" {
		serviceName = envName
	}
	attributes := []attribute.KeyValue{
		semconv.ServiceNameKey.String(serviceName),
	}
	if build.Version() != "" {
		attributes = append(attributes, semconv.ServiceVersionKey.String(build.Version()))
	}
	detectors, err := newDetectors(detectorTypes)
	if err != nil {
		return nil, err
	}
	return resource.New(ctx,
		resource.WithDetectors(detectors...),
		resource.WithTelemetrySDK(),
		resource.WithFromEnv(),
		resource.WithAttributes(attributes...),
	)
}

// newDetectors resolves the configured detector types, rejecting any value
// that does not name a detector. An empty list yields no detectors, which is
// the default: detection is opt-in because it costs a platform probe on startup
// and adds attributes that identify the account and host we run on to
// everything we export.
func newDetectors(types []DetectorType) ([]resource.Detector, error) {
	detectors := make([]resource.Detector, 0, len(types))
	for _, typ := range types {
		switch typ {
		case DetectorTypeGoogle:
			// Fills in faas.*, cloud.* and host.* attributes when we run on
			// Cloud Run, GKE, GCE, App Engine or Cloud Functions. Google Cloud
			// Monitoring derives the monitored resource a metric is written to
			// from exactly those attributes, so without them every process
			// reports as the same empty generic_node and the series overlap.
			detectors = append(detectors, google_detector.NewDetector())
		case DetectorTypeUnspecified:
			// The zero value is not a valid entry: an empty string in the list
			// is a typo, not a way to switch detection off. An empty list is.
			return nil, errDetectorType(typ)
		default:
			return nil, errDetectorType(typ)
		}
	}
	return detectors, nil
}

func errExporterType(typ ExporterType, instrument string) error {
	return fmt.Errorf("exporter type \"%v\" unsupported for %s", typ, instrument)
}

func errDetectorType(typ DetectorType) error {
	return fmt.Errorf("detector type \"%v\" unsupported", typ)
}
