package otel

import (
	"context"
	"github.com/caos/zitadel/internal/metrics"
	"go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/otlp"
)

type Meter struct {
	meter    metric.Meter
	counters map[string]metric.Int64Counter
}

func NewMeter(name string, exporter *otlp.Exporter) metrics.Meter {
	return &Meter{
		meter: global.MeterProvider().Meter(
			name,
			metric.WithInstrumentationVersion(contrib.SemVersion()),
		),
	}
}

func (m *Meter) AddMeasure() context.Context {

}
