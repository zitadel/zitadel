package config

import (
	"github.com/zitadel/zitadel/v2/internal/errors"
	"github.com/zitadel/zitadel/v2/internal/telemetry/metrics/otel"
)

type Config struct {
	Type   string
	Config map[string]interface{} `mapstructure:",remain"`
}

var meter = map[string]func(map[string]interface{}) error{
	"otel": otel.NewTracerFromConfig,
	"none": NoMetrics,
	"":     NoMetrics,
}

func (c *Config) NewMeter() error {
	t, ok := meter[c.Type]
	if !ok {
		return errors.ThrowInternalf(nil, "METER-Dfqsx", "config type %s not supported", c.Type)
	}

	return t(c.Config)
}

func NoMetrics(_ map[string]interface{}) error {
	return nil
}
