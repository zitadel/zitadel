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

type DecorateOption func(*DecorateOptions)

type DecorateOptions struct {
	startOpts []trace.SpanStartOption
	endOpts   []trace.SpanEndOption

	spanName string

	span trace.Span
}

func WithSpanName(name string) DecorateOption {
	return func(o *DecorateOptions) {
		o.spanName = name
	}
}

func WithSpanStartOptions(opts ...trace.SpanStartOption) DecorateOption {
	return func(o *DecorateOptions) {
		o.startOpts = append(o.startOpts, opts...)
	}
}

func WithSpanEndOptions(opts ...trace.SpanEndOption) DecorateOption {
	return func(o *DecorateOptions) {
		o.endOpts = append(o.endOpts, opts...)
	}
}

func (o *DecorateOptions) Start(ctx context.Context, tracer *Tracer) context.Context {
	if o.spanName == "" {
		o.spanName = functionName()
	}
	ctx, o.span = tracer.Tracer.Start(ctx, o.spanName, o.startOpts...)
	return ctx
}

func (o *DecorateOptions) End(err error) {
	o.span.RecordError(err)
	o.span.End(o.endOpts...)
}

// func (t Tracer) Decorate(ctx context.Context, fn func(ctx context.Context) error, opts ...DecorateOption) {
// 	o := new(DecorateOptions)
// 	for _, opt := range opts {
// 		opt(o)
// 	}

// 	if o.spanName == "" {
// 		o.spanName = functionName()
// 	}

// 	ctx, span := t.Tracer.Start(ctx, o.spanName, o.startOpts...)
// 	defer span.End(o.endOpts...)

// 	err := fn(ctx)
// 	span.RecordError(err)
// }

func functionName() string {
	counter, _, _, success := runtime.Caller(2)

	if !success {
		return "zitadel"
	}

	return runtime.FuncForPC(counter).Name()
}
