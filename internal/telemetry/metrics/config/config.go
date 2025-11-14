package config

import (
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
	"github.com/zitadel/zitadel/internal/telemetry/metrics/otel"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Config struct {
	Type   string
	Config map[string]interface{} `mapstructure:",remain"`
}

var meter = map[string]func(map[string]interface{}) error{
	"otel": otel.NewTracerFromConfig,
	"none": registerNoopMetrics,
	"":     registerNoopMetrics,
}

func (c *Config) NewMeter() error {
	// When using start-from-init or start-from-setup the metric provider
	// was already set in the setup phase and the start phase must not overwrite it.
	if metrics.M != nil {
		return nil
	}
	t, ok := meter[c.Type]
	if !ok {
		return zerrors.ThrowInternalf(nil, "METER-Dfqsx", "config type %s not supported", c.Type)
	}

	return t(c.Config)
}

func registerNoopMetrics(rawConfig map[string]interface{}) (err error) {
	metrics.M = &metrics.NoopMetrics{}
	return nil
}
