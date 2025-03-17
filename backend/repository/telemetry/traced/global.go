package traced

import (
	"context"
	"log"

	"github.com/zitadel/zitadel/backend/repository/orchestrate/handler"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

// Wrap decorates the given handle function with tracing.
// The function is safe to call with nil tracer.
func Wrap[Req, Res any](tracer *tracing.Tracer, name string, handle handler.Handler[Req, Res]) handler.Handler[Req, Res] {
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

// Decorate decorates the given handle function with tracing.
// The function is safe to call with nil tracer.
func Decorate[Req, Res any](tracer *tracing.Tracer, opts ...tracing.DecorateOption) handler.Decorator[Req, Res] {
	return func(ctx context.Context, r Req, handle handler.Handler[Req, Res]) (_ Res, err error) {
		if tracer == nil {
			return handle(ctx, r)
		}
		o := new(tracing.DecorateOptions)
		for _, opt := range opts {
			opt(o)
		}
		log.Println("traced.decorate")

		ctx, end := o.Start(ctx, tracer)
		defer end(err)
		return handle(ctx, r)
	}
}
