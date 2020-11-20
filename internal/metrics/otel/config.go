package otel

import (
	"github.com/caos/zitadel/internal/metrics"
)

type Config struct {
	meterName string
}

func (c *Config) NewMetrics() error {
	metrics.M = NewMeter(c.meterName)
	return nil
}
