package log

import (
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/telemetry/tracing/otel"
)

type Config struct {
	Fraction    float64
	ServiceName string
}

func NewTracer(rawConfig map[string]interface{}) (err error) {
	c := new(Config)
	c.Fraction, err = otel.FractionFromConfig(rawConfig["fraction"])
	c.ServiceName, _ = rawConfig["servicename"].(string)
	if err != nil {
		return err
	}
	return c.NewTracer()
}

type Tracer struct {
	otel.Tracer
}

func (c *Config) NewTracer() error {
	sampler := otel.NewSampler(sdk_trace.TraceIDRatioBased(c.Fraction))
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		return err
	}

	tracing.T, err = otel.NewTracer(sampler, exporter, c.ServiceName)
	return err
}
