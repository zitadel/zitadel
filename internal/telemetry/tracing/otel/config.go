package otel

import (
	"context"
	"strconv"

	otlpgrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	api_trace "go.opentelemetry.io/otel/trace"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Config struct {
	Fraction    float64
	Endpoint    string
	ServiceName string
}

func NewTracerFromConfig(rawConfig map[string]interface{}) (err error) {
	c := new(Config)
	c.Endpoint, _ = rawConfig["endpoint"].(string)
	c.ServiceName, _ = rawConfig["servicename"].(string)
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
	sampler := NewSampler(sdk_trace.TraceIDRatioBased(c.Fraction))
	exporter, err := otlpgrpc.New(context.Background(), otlpgrpc.WithEndpoint(c.Endpoint), otlpgrpc.WithInsecure())
	if err != nil {
		return err
	}

	tracing.T, err = NewTracer(sampler, exporter, c.ServiceName)
	return err
}

// NewSampler returns a sampler decorator which behaves differently,
// based on the parent of the span. If the span has no parent and is of kind server,
// the decorated sampler is used to make sampling decision.
// If the span has a parent, depending on whether the parent is remote and whether it
// is sampled, one of the following samplers will apply:
//   - remote parent sampled -> always sample
//   - remote parent not sampled -> sample based on the decorated sampler (fraction based)
//   - local parent sampled -> always sample
//   - local parent not sampled -> never sample
func NewSampler(sampler sdk_trace.Sampler) sdk_trace.Sampler {
	return sdk_trace.ParentBased(
		tracing.SpanKindBased(sampler, api_trace.SpanKindServer),
		sdk_trace.WithRemoteParentNotSampled(sampler),
	)
}
