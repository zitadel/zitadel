package orchestrate

import (
	"github.com/zitadel/zitadel/backend/telemetry/logging"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

type options struct {
	tracer *tracing.Tracer
	logger *logging.Logger
}

type Option func(*options)

func WithTracer(tracer *tracing.Tracer) Option {
	return func(o *options) {
		o.tracer = tracer
	}
}

func WithLogger(logger *logging.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}
