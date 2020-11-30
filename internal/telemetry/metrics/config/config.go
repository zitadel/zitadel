package config

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/metrics"
	"github.com/caos/zitadel/internal/telemetry/metrics/otel"
)

type MetricsConfig struct {
	Type   string
	Config metrics.Config
}

var meter = map[string]func() metrics.Config{
	"otel": func() metrics.Config { return &otel.Config{} },
	"none": func() metrics.Config { return &NoMetrics{} },
	"":     func() metrics.Config { return &NoMetrics{} },
}

func (c *MetricsConfig) UnmarshalJSON(data []byte) error {
	var rc struct {
		Type   string
		Config json.RawMessage
	}

	if err := json.Unmarshal(data, &rc); err != nil {
		return errors.ThrowInternal(err, "METER-4M9so", "error parsing config")
	}

	c.Type = rc.Type

	var err error
	c.Config, err = newMetricsConfig(c.Type, rc.Config)
	if err != nil {
		return err
	}

	return c.Config.NewMetrics()
}

func newMetricsConfig(tracerType string, configData []byte) (metrics.Config, error) {
	t, ok := meter[tracerType]
	if !ok {
		return nil, errors.ThrowInternalf(nil, "METER-3M0ps", "config type %s not supported", tracerType)
	}

	metricsConfig := t()
	if len(configData) == 0 {
		return metricsConfig, nil
	}

	if err := json.Unmarshal(configData, metricsConfig); err != nil {
		return nil, errors.ThrowInternal(err, "METER-4M9sf", "Could not read config: %v")
	}

	return metricsConfig, nil
}

type NoMetrics struct{}

func (_ *NoMetrics) NewMetrics() error {
	return nil
}
