package otel

import (
	"github.com/caos/zitadel/internal/tracing"
	"go.opentelemetry.io/otel/exporters/otlp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// import "github.com/caos/zitadel/internal/tracing"

// type Config struct {
// 	MetricPrefix string
// 	Fraction     float64
// }

// func (c *Config) NewTracer() error {

// 	// tracing.T = &Tracer{projectID: c.ProjectID, metricPrefix: c.MetricPrefix, sampler: trace.ProbabilitySampler(c.Fraction)}

// 	return tracing.T.Start()
// }

type Config struct {
	Fraction     float64
	MetrixPrefix string
	Endpoint     string
}

func (c *Config) NewTracer() error {
	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(c.Fraction))
	exporter, err := otlp.NewExporter(otlp.WithAddress(c.Endpoint), otlp.WithInsecure())
	// exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		return err
	}

	tracing.T = NewTracer(c.MetrixPrefix, sampler, exporter)
	return nil
}
