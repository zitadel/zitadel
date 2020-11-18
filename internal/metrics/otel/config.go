package otel

import (
	"github.com/caos/zitadel/internal/metrics"
	"go.opentelemetry.io/otel/exporters/otlp"
)

type Config struct {
	Endpoint string
}

func (c *Config) NewMetrics() error {
	exporter, err := otlp.NewExporter(otlp.WithAddress(c.Endpoint), otlp.WithInsecure())
	if err != nil {
		return err
	}

	metrics.M = NewMeter(exporter)
	return nil
}
