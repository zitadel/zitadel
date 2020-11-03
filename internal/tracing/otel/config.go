package otel

import (
	"github.com/caos/zitadel/internal/tracing"
	"go.opentelemetry.io/otel/exporters/otlp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	Fraction     float64
	MetrixPrefix string
	Endpoint     string
}

func (c *Config) NewTracer() error {
	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(c.Fraction))
	exporter, err := otlp.NewExporter(otlp.WithAddress(c.Endpoint), otlp.WithInsecure())
	if err != nil {
		return err
	}

	tracing.T = NewTracer(c.MetrixPrefix, sampler, exporter)
	return nil
}
