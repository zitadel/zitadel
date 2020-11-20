package otel

import (
	"github.com/caos/zitadel/internal/metrics"
	"go.opentelemetry.io/otel/exporters/otlp"
)

type Config struct {
}

func (c *Config) NewMetrics() error {
	metrics.M = NewMeter()
	return nil
}
