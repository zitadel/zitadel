package otel

import (
	"context"

	"github.com/caos/zitadel/internal/tracing"
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
	exporter, err := otlp.NewExporter(context.Background(), otlp.WithAddress(c.Endpoint), otlp.WithInsecure())
	if err != nil {
		return err
	}

	tracing.T = NewTracer(c.MetricPrefix, sampler, exporter)
	return nil
}
