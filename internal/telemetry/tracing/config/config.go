package config

import (
	"github.com/zitadel/zitadel/v2/internal/errors"
	"github.com/zitadel/zitadel/v2/internal/telemetry/tracing/google"
	"github.com/zitadel/zitadel/v2/internal/telemetry/tracing/log"
	"github.com/zitadel/zitadel/v2/internal/telemetry/tracing/otel"
)

type Config struct {
	Type   string
	Config map[string]interface{} `mapstructure:",remain"`
}

func (c *Config) NewTracer() error {
	t, ok := tracer[c.Type]
	if !ok {
		return errors.ThrowInternalf(nil, "TRACE-dsbjh", "config type %s not supported", c.Type)
	}

	return t(c.Config)
}

var tracer = map[string]func(map[string]interface{}) error{
	"otel":   otel.NewTracerFromConfig,
	"google": google.NewTracer,
	"log":    log.NewTracer,
	"none":   NoTracer,
	"":       NoTracer,
}

func NoTracer(_ map[string]interface{}) error {
	return nil
}
