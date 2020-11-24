package otel

import (
	"github.com/caos/zitadel/internal/telemetry/metrics"
)

type Config struct {
}

func (c *Config) NewMetrics() (err error) {
	metrics.M, err = NewMetrics()
	return err
}
