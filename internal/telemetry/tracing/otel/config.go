package otel

import (
	"context"
	"strconv"

	otlpgrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type Config struct {
	Fraction     float64
	MetricPrefix string
	Endpoint     string
}

func NewTracerFromConfig(rawConfig map[string]interface{}) (err error) {
	c := new(Config)
	c.Endpoint, _ = rawConfig["endpoint"].(string)
	c.MetricPrefix, _ = rawConfig["metricprefix"].(string)
	fraction, ok := rawConfig["fraction"].(string)
	if ok {
		c.Fraction, err = strconv.ParseFloat(fraction, 32)
		if err != nil {
			return errors.ThrowInternal(err, "OTEL-Dd2s", "could not map fraction")
		}
	}

	return c.NewTracer()
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
