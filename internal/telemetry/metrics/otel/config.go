package otel

import (
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

type Config struct {
	MeterName string
}

func NewTracerFromConfig(rawConfig map[string]interface{}) (err error) {
	metrics.M, err = NewMetrics(rawConfig)
	return err
}
