package config

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/telemetry/tracing/google"
	"github.com/caos/zitadel/internal/telemetry/tracing/log"
	"github.com/caos/zitadel/internal/telemetry/tracing/otel"
)

type TracingConfig struct {
	Type   string
	Config tracing.Config
}

var tracer = map[string]func() tracing.Config{
	"otel":   func() tracing.Config { return &otel.Config{} },
	"google": func() tracing.Config { return &google.Config{} },
	"log":    func() tracing.Config { return &log.Config{} },
	"none":   func() tracing.Config { return &NoTracing{} },
	"":       func() tracing.Config { return &NoTracing{} },
}

func (c *TracingConfig) UnmarshalJSON(data []byte) error {
	var rc struct {
		Type   string
		Config json.RawMessage
	}

	if err := json.Unmarshal(data, &rc); err != nil {
		return errors.ThrowInternal(err, "TRACE-vmjS", "error parsing config")
	}

	c.Type = rc.Type

	var err error
	c.Config, err = newTracingConfig(c.Type, rc.Config)
	if err != nil {
		return err
	}

	return c.Config.NewTracer()
}

func newTracingConfig(tracerType string, configData []byte) (tracing.Config, error) {
	t, ok := tracer[tracerType]
	if !ok {
		return nil, errors.ThrowInternalf(nil, "TRACE-HMEJ", "config type %s not supported", tracerType)
	}

	tracingConfig := t()
	if len(configData) == 0 {
		return tracingConfig, nil
	}

	if err := json.Unmarshal(configData, tracingConfig); err != nil {
		return nil, errors.ThrowInternal(err, "TRACE-1tSS", "Could not read config: %v")
	}

	return tracingConfig, nil
}

type NoTracing struct{}

func (_ *NoTracing) NewTracer() error {
	return nil
}
