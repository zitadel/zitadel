package otel

import (
	"context"
	"strconv"

	otlpgrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Config struct {
	Fraction float64
	Endpoint string
}

func NewTracerFromConfig(rawConfig map[string]interface{}) (err error) {
	c := new(Config)
	c.Endpoint, _ = rawConfig["endpoint"].(string)
	c.Fraction, err = FractionFromConfig(rawConfig["fraction"])
	if err != nil {
		return err
	}
	return c.NewTracer()
}

func FractionFromConfig(i interface{}) (float64, error) {
	if i == nil {
		return 0, nil
	}
	switch fraction := i.(type) {
	case float64:
		return fraction, nil
	case int:
		return float64(fraction), nil
	case string:
		f, err := strconv.ParseFloat(fraction, 64)
		if err != nil {
			return 0, zerrors.ThrowInternal(err, "OTEL-SAfe1", "could not map fraction")
		}
		return f, nil
	default:
		return 0, zerrors.ThrowInternal(nil, "OTEL-Dd2s", "could not map fraction, unknown type")
	}
}

func (c *Config) NewTracer() error {
	sampler := sdk_trace.ParentBased(sdk_trace.TraceIDRatioBased(c.Fraction))
	exporter, err := otlpgrpc.New(context.Background(), otlpgrpc.WithEndpoint(c.Endpoint), otlpgrpc.WithInsecure())
	if err != nil {
		return err
	}

	tracing.T, err = NewTracer(sampler, exporter)
	return err
}
