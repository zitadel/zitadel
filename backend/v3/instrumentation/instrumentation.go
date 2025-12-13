// Package instrumentation initializes and configures telemetry and logging instrumentation.
// After [Start] is called, the instrumentation providers can be accessed using global getters.
package instrumentation

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"sync"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	ServiceName string
	Trace       TraceConfig
	Metric      MetricConfig
	Log         LogConfig
	Profile     ProfileConfig
}

type ShutdownFunc func(context.Context) error

var (
	startOnce sync.Once
	shutdown  ShutdownFunc
	startErr  error
)

// Start initializes global instrumentation once, based on the provided configuration.
// It is safe to call multiple times; subsequent calls will return the same shutdown function and error as the first call.
func Start(ctx context.Context, cfg Config) (ShutdownFunc, error) {
	startOnce.Do(func() {
		startErr = startProfiler(cfg.Profile, cfg.ServiceName)
		if startErr != nil {
			return
		}
		var logProvider *log.LoggerProvider
		logProvider, shutdown, startErr = setupOTelSDK(ctx, cfg)
		if startErr != nil {
			return
		}
		setLogger(logProvider, cfg.Log)
	})
	return shutdown, startErr
}

// NewTracer creates a new tracer from the global provider.
// See [trace.TracerProvider.Tracer] for details.
func NewTracer(name string, options ...trace.TracerOption) *Tracer {
	provider := otel.GetTracerProvider()
	return &Tracer{provider.Tracer(name, options...)}
}

// NewMeter creates a new meter from the global provider.
// See [metric.MeterProvider.Meter] for details.
func NewMeter(name string, options ...metric.MeterOption) *Meter {
	provider := otel.GetMeterProvider()
	return &Meter{
		Meter: provider.Meter(name, options...),
	}
}

// Logger returns the globally configured logger.
func Logger() *slog.Logger {
	return slog.Default()
}

func RequestFilter(ignoredPrefix ...string) otelhttp.Filter {
	return func(r *http.Request) bool {
		return !slices.ContainsFunc(
			ignoredPrefix,
			func(endpoint string) bool {
				return strings.HasPrefix(r.URL.RequestURI(), endpoint)
			},
		)
	}
}
