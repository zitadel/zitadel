package tracing

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct{ trace.Tracer }

func NewTracer(name string) Tracer {
	return Tracer{otel.Tracer(name)}
}

type DecorateOption func(*decorateOptions)

type decorateOptions struct {
	startOpts []trace.SpanStartOption
	endOpts   []trace.SpanEndOption

	spanName string
}

func WithSpanName(name string) DecorateOption {
	return func(o *decorateOptions) {
		o.spanName = name
	}
}

func WithSpanStartOptions(opts ...trace.SpanStartOption) DecorateOption {
	return func(o *decorateOptions) {
		o.startOpts = append(o.startOpts, opts...)
	}
}

func WithSpanEndOptions(opts ...trace.SpanEndOption) DecorateOption {
	return func(o *decorateOptions) {
		o.endOpts = append(o.endOpts, opts...)
	}
}

func (t Tracer) Decorate(ctx context.Context, fn func(ctx context.Context) error, opts ...DecorateOption) {
	o := new(decorateOptions)
	for _, opt := range opts {
		opt(o)
	}

	if o.spanName == "" {
		o.spanName = functionName()
	}

	_, span := t.Tracer.Start(ctx, o.spanName, o.startOpts...)
	defer span.End(o.endOpts...)

	if err := fn(ctx); err != nil {
		span.RecordError(err)
	}
}

func functionName() string {
	counter, _, _, success := runtime.Caller(2)

	if !success {
		return "zitadel"
	}

	return runtime.FuncForPC(counter).Name()
}
