package otel

import (
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

type Config struct {
	MeterName string
}

func NewTracerFromConfig(rawConfig map[string]interface{}) (err error) {
	c := new(Config)
	c.MeterName, _ = rawConfig["metername"].(string)
	return c.NewMetrics()
}

func (c *Config) NewMetrics() (err error) {
	metrics.M, err = NewMetrics(c.MeterName)
	return err
}
