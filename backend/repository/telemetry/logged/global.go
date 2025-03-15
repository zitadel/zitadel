package logged

import (
	"context"
	"log/slog"

	"github.com/zitadel/zitadel/backend/repository/orchestrate/handler"
	"github.com/zitadel/zitadel/backend/telemetry/logging"
)

// Wrap decorates the given handle function with logging.
// The function is safe to call with nil logger.
func Wrap[Req, Res any](logger *logging.Logger, name string, handle handler.Handle[Req, Res]) handler.Handle[Req, Res] {
	if logger == nil {
		return handle
	}
	return func(ctx context.Context, r Req) (_ Res, err error) {
		logger.Debug("execute", slog.String("handler", name))
		defer logger.Debug("done", slog.String("handler", name))
		return handle(ctx, r)
	}
}

func WrapInside(logger *logging.Logger, name string) func(ctx context.Context, fn func(context.Context) error) {
	logger = logger.With(slog.String("handler", name))
	return func(ctx context.Context, fn func(context.Context) error) {
		logger.Debug("execute")
		var err error
		defer func() {
			if err != nil {
				logger.Error("failed", slog.String("cause", err.Error()))
			}
			logger.Debug("done")
		}()
		err = fn(ctx)
	}
}

func DecorateHandle[Req, Res any](logger *logging.Logger, handle func(context.Context, Req) (Res, error)) func(ctx context.Context, r Req) (_ Res, err error) {
	return func(ctx context.Context, r Req) (_ Res, err error) {
		logger.DebugContext(ctx, "execute")
		defer func() {
			if err != nil {
				logger.ErrorContext(ctx, "failed", slog.String("cause", err.Error()))
			}
			logger.DebugContext(ctx, "done")
		}()
		return handle(ctx, r)
	}
}

// // Handler wraps the given handle function with logging.
// // The function is safe to call with nil logger.
// func Handler[Req, Res any, H handler.Handle[Req, Res]](logger *logging.Logger, name string, handle H) *handler.Handler[Req, Res, H] {
// 	return &handler.Handler[Req, Res, H]{
// 		Handle: Wrap(logger, name, handle),
// 	}
// }

// // Chained wraps the given handle function with logging.
// // The function is safe to call with nil logger.
// // The next handler is called after the handle function.
// func Chained[Req, Res any, H, N handler.Handle[Req, Res]](logger *logging.Logger, name string, handle H, next N) *handler.Chained[Req, Res, H, N] {
// 	return handler.NewChained(
// 		Wrap(logger, name, handle),
// 		next,
// 	)
// }
