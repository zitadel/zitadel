package otel

import (
	"context"

	"github.com/caos/zitadel/internal/telemetry/tracing"
	otlpgrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	Fraction     float64
	MetricPrefix string
	Endpoint     string
}

func (c *Config) NewTracer() error {
	sampler := sdk_trace.ParentBased(sdk_trace.TraceIDRatioBased(c.Fraction))
	exporter, err := otlpgrpc.New(context.Background(), otlpgrpc.WithEndpoint(c.Endpoint), otlpgrpc.WithInsecure())
	if err != nil {
		return err
	}

	tracing.T = NewTracer(c.MetricPrefix, sampler, exporter)
	return nil
}
