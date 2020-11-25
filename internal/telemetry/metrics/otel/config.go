package otel

import (
	"github.com/caos/zitadel/internal/telemetry/metrics"
)

type Config struct {
	MeterName string
}

func (c *Config) NewMetrics() (err error) {
	metrics.M, err = NewMetrics(c.MeterName)
	return err
}
