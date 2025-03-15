package traced

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository/orchestrate/handler"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

// Wrap decorates the given handle function with tracing.
// The function is safe to call with nil tracer.
func Wrap[Req, Res any](tracer *tracing.Tracer, name string, handle handler.Handle[Req, Res]) handler.Handle[Req, Res] {
	if tracer == nil {
		return handle
	}
	return func(ctx context.Context, r Req) (_ Res, err error) {
		ctx, span := tracer.Start(
			ctx,
			name,
		)
		defer func() {
			if err != nil {
				span.RecordError(err)
			}
			span.End()
		}()
		return handle(ctx, r)
	}
}

func WrapInside(tracer *tracing.Tracer, name string) func(ctx context.Context, fn func() error) {
	return func(ctx context.Context, fn func() error) {
		var err error
		_, span := tracer.Start(
			ctx,
			name,
		)
		defer func() {
			if err != nil {
				span.RecordError(err)
			}
			span.End()
		}()
		err = fn()
	}
}

func DecorateHandle[Req, Res any](tracer *tracing.Tracer, opts ...tracing.DecorateOption) handler.Decorate[Req, Res] {
	return func(ctx context.Context, r Req, handle handler.Handle[Req, Res]) (_ Res, err error) {
		o := new(tracing.DecorateOptions)
		for _, opt := range opts {
			opt(o)
		}

		ctx = o.Start(ctx, tracer)
		defer o.End(err)
		return handle(ctx, r)
	}
}

// // Handler wraps the given handle function with tracing.
// // The function is safe to call with nil logger.
// func Handler[Req, Res any, H handler.Handle[Req, Res]](tracer *tracing.Tracer, name string, handle H) *handler.Handler[Req, Res, H] {
// 	return &handler.Handler[Req, Res, H]{
// 		Handle: Wrap(tracer, name, handle),
// 	}
// }

// // Chained wraps the given handle function with tracing.
// // The function is safe to call with nil logger.
// // The next handler is called after the handle function.
// func Chained[Req, Res any, H, N handler.Handle[Req, Res]](tracer *tracing.Tracer, name string, handle H, next N) *handler.Chained[Req, Res, H, N] {
// 	return handler.NewChained(
// 		Wrap(tracer, name, handle),
// 		next,
// 	)
// }
