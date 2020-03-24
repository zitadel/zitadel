package google

import (
	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/caos/zitadel/internal/tracing"
)

type Config struct {
	ProjectID    string
	MetricPrefix string
	Fraction     float64
}

func (c *Config) NewTracer() error {
	if !envIsSet() {
		return status.Error(codes.InvalidArgument, "env not properly set, GOOGLE_APPLICATION_CREDENTIALS is misconfigured or missing")
	}

	tracing.T = &Tracer{projectID: c.ProjectID, metricPrefix: c.MetricPrefix, sampler: trace.ProbabilitySampler(c.Fraction)}

	return tracing.T.Start()
}
