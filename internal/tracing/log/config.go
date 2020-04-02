package log

import (
	"go.opencensus.io/trace"

	"github.com/caos/zitadel/internal/tracing"
)

type Config struct {
	Fraction float64
}

func (c *Config) NewTracer() error {
	if c.Fraction < 1 {
		c.Fraction = 1
	}

	tracing.T = &Tracer{trace.ProbabilitySampler(c.Fraction)}

	return tracing.T.Start()
}
