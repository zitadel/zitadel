package google

import (
	"go.opencensus.io/trace"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/tracing"
)

type Config struct {
	ProjectID    string
	MetricPrefix string
	Fraction     float64
}

func (c *Config) NewTracer() error {
	if !envIsSet() {
		return errors.ThrowInvalidArgument(nil, "GOOGL-sdh3a", "env not properly set, GOOGLE_APPLICATION_CREDENTIALS is misconfigured or missing")
	}

	tracing.T = &Tracer{projectID: c.ProjectID, metricPrefix: c.MetricPrefix, sampler: trace.ProbabilitySampler(c.Fraction)}

	return tracing.T.Start()
}
