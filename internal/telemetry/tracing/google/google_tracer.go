package google

import (
	"strconv"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/telemetry/tracing/otel"
)

type Config struct {
	ProjectID    string
	MetricPrefix string
	Fraction     float64
}

func NewTracer(rawConfig map[string]interface{}) (err error) {
	c := new(Config)
	c.ProjectID, _ = rawConfig["projectid"].(string)
	c.MetricPrefix, _ = rawConfig["metricprefix"].(string)
	fraction, ok := rawConfig["fraction"].(string)
	if ok {
		c.Fraction, err = strconv.ParseFloat(fraction, 32)
		if err != nil {
			return errors.ThrowInternal(err, "GOOGLE-Dsag3", "could not map fraction")
		}
	}

	return c.NewTracer()
}

type Tracer struct {
	otel.Tracer
}

func (c *Config) NewTracer() error {
	sampler := sdk_trace.ParentBased(sdk_trace.TraceIDRatioBased(c.Fraction))
	exporter, err := texporter.New(texporter.WithProjectID(c.ProjectID))
	if err != nil {
		return err
	}

	tracing.T = &Tracer{Tracer: *(otel.NewTracer(c.MetricPrefix, sampler, exporter))}

	return nil
}
