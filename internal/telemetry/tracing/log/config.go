package log

import (
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/telemetry/tracing/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	Fraction     float64
	MetricPrefix string
}

type Tracer struct {
	otel.Tracer
}

func (c *Config) NewTracer() error {
	sampler := sdk_trace.ParentBased(sdk_trace.TraceIDRatioBased(c.Fraction))
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		return err
	}

	tracing.T = &Tracer{Tracer: *(otel.NewTracer(c.MetricPrefix, sampler, exporter))}
	return nil
}
