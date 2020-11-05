package log

import (
	"github.com/caos/zitadel/internal/tracing"
	"github.com/caos/zitadel/internal/tracing/otel"
	"go.opentelemetry.io/otel/exporters/stdout"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	Fraction     float64
	MetricPrefix string
}

type Tracer struct {
	otel.Tracer
}

func (c *Config) NewTracer() error {
	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(c.Fraction))
	exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		return err
	}

	tracing.T = &Tracer{Tracer: *(otel.NewTracer(c.MetricPrefix, sampler, exporter))}
	return nil
}
