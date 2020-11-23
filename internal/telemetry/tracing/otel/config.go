package otel

import (
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"go.opentelemetry.io/otel/exporters/otlp"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	Fraction     float64
	MetricPrefix string
	Endpoint     string
}

func (c *Config) NewTracer() error {
	sampler := sdk_trace.ParentBased(sdk_trace.TraceIDRatioBased(c.Fraction))
	exporter, err := otlp.NewExporter(otlp.WithAddress(c.Endpoint), otlp.WithInsecure())
	if err != nil {
		return err
	}

	tracing.T = NewTracer(c.MetricPrefix, sampler, exporter)
	return nil
}
