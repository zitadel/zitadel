package otel

import (
	"github.com/caos/zitadel/internal/telemetry/metrics"
)

type Config struct {
	MeterName string
}

func (c *Config) NewMetrics() error {
	metrics.M = NewMeter(c.MeterName)
	return nil
}
