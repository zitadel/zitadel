package config

import (
	"github.com/zitadel/zitadel/internal/telemetry/profiler/google"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Config struct {
	Type   string
	Config map[string]any `mapstructure:",remain"`
}

var profiler = map[string]func(map[string]any) error{
	"google": google.NewProfiler,
	"none":   NoProfiler,
	"":       NoProfiler,
}

func (c *Config) NewProfiler() error {
	t, ok := profiler[c.Type]
	if !ok {
		return zerrors.ThrowInternalf(nil, "PROFI-Dfqsx", "config type %s not supported", c.Type)
	}

	return t(c.Config)
}

func NoProfiler(_ map[string]any) error {
	return nil
}
