package google

import (
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/telemetry/tracing/otel"
)

type Config struct {
	ProjectID   string
	Fraction    float64
	ServiceName string
}

func NewTracer(rawConfig map[string]interface{}) (err error) {
	c := new(Config)
	c.ProjectID, _ = rawConfig["projectid"].(string)
	c.ServiceName, _ = rawConfig["servicename"].(string)
	c.Fraction, err = otel.FractionFromConfig(rawConfig["fraction"])
	if err != nil {
		return err
	}
	return c.NewTracer()
}

type Tracer struct {
	otel.Tracer
}

func (c *Config) NewTracer() error {
	sampler := otel.NewSampler(sdk_trace.TraceIDRatioBased(c.Fraction))
	exporter, err := texporter.New(texporter.WithProjectID(c.ProjectID))
	if err != nil {
		return err
	}

	tracing.T, err = otel.NewTracer(sampler, exporter, c.ServiceName)
	return err
}
