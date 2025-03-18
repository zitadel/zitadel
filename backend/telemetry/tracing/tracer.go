package tracing

import (
	"context"
	"log"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/zitadel/zitadel/backend/handler"
)

type Tracer struct{ trace.Tracer }

func NewTracer(name string) *Tracer {
	return &Tracer{otel.Tracer(name)}
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

// Wrap decorates the given handle function with tracing.
// The function is safe to call with nil tracer.
func Wrap[Req, Res any](tracer *Tracer, name string, handle handler.Handle[Req, Res]) handler.Handle[Req, Res] {
	if tracer == nil {
		return handle
	}
	return func(ctx context.Context, r Req) (_ Res, err error) {
		ctx, span := tracer.Start(
			ctx,
			name,
		)
		log.Println("trace.wrap", name)
		defer func() {
			if err != nil {
				span.RecordError(err)
			}
			span.End()
		}()
		return handle(ctx, r)
	}
}

// Decorate decorates the given handle function with
// The function is safe to call with nil tracer.
func Decorate[Req, Res any](tracer *Tracer, opts ...DecorateOption) handler.Middleware[Req, Res] {
	return func(ctx context.Context, r Req, handle handler.Handle[Req, Res]) (_ Res, err error) {
		if tracer == nil {
			return handle(ctx, r)
		}
		o := new(DecorateOptions)
		for _, opt := range opts {
			opt(o)
		}
		log.Println("traced.decorate")

		ctx, end := o.Start(ctx, tracer)
		defer end(err)
		return handle(ctx, r)
	}
}

func (o *DecorateOptions) Start(ctx context.Context, tracer *Tracer) (context.Context, func(error)) {
	if o.spanName == "" {
		o.spanName = functionName()
	}
	ctx, o.span = tracer.Tracer.Start(ctx, o.spanName, o.startOpts...)
	return ctx, o.end
}

func (o *DecorateOptions) end(err error) {
	o.span.RecordError(err)
	o.span.End(o.endOpts...)
}

func functionName() string {
	counter, _, _, success := runtime.Caller(2)

	if !success {
		return "zitadel"
	}

	return runtime.FuncForPC(counter).Name()
}
