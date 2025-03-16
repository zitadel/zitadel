package orchestrate

import (
	"github.com/zitadel/zitadel/backend/telemetry/logging"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

// options are the default options for orchestrators.
type options[T any] struct {
	custom *T
	tracer *tracing.Tracer
	logger *logging.Logger
}

func newOptions[T any]() options[T] {
	return options[T]{
		custom: new(T),
	}
}

type applier interface {
	apply()
}

type Option[T any] func(*options[T])

func WithTracer[T any](tracer *tracing.Tracer) Option[T] {
	return func(o *options[T]) {
		o.tracer = tracer
	}
}

func WithLogger[T any](logger *logging.Logger) Option[T] {
	return func(o *options[T]) {
		o.logger = logger
	}
}

func (o Option[T]) apply(opts *options[T]) {
	o(opts)
}
