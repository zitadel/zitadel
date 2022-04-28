package log

import (
	"strconv"

	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/telemetry/tracing/otel"
)

type Config struct {
	Fraction float64
}

func NewTracer(rawConfig map[string]interface{}) (err error) {
	c := new(Config)
	fraction, ok := rawConfig["fraction"].(string)
	if ok {
		c.Fraction, err = strconv.ParseFloat(fraction, 32)
		if err != nil {
			return errors.ThrowInternal(err, "LOG-Dsag3", "could not map fraction")
		}
	}

	return c.NewTracer()
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

	tracing.T, err = otel.NewTracer(sampler, exporter)
	return err
}
