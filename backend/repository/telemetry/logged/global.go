package logged

import (
	"context"
	"log"
	"log/slog"

	"github.com/zitadel/zitadel/backend/repository/orchestrate/handler"
	"github.com/zitadel/zitadel/backend/telemetry/logging"
)

// Wrap decorates the given handle function with logging.
// The function is safe to call with nil logger.
func Wrap[Req, Res any](logger *logging.Logger, name string, handle handler.Handler[Req, Res]) handler.Handler[Req, Res] {
	if logger == nil {
		return handle
	}
	return func(ctx context.Context, r Req) (_ Res, err error) {
		logger.Debug("execute", slog.String("handler", name))
		defer logger.Debug("done", slog.String("handler", name))
		log.Println("log.wrap", name)
		return handle(ctx, r)
	}
}

// Decorate decorates the given handle function with logging.
// The function is safe to call with nil logger.
func Decorate[Req, Res any](logger *logging.Logger, name string) handler.Decorator[Req, Res] {
	return func(ctx context.Context, request Req, handle handler.Handler[Req, Res]) (res Res, err error) {
		if logger == nil {
			return handle(ctx, request)
		}
		logger = logger.With("handler", name)
		logger.DebugContext(ctx, "execute")
		log.Println("logged.decorate", name)
		defer func() {
			if err != nil {
				logger.ErrorContext(ctx, "failed", slog.String("cause", err.Error()))
			}
			logger.DebugContext(ctx, "done")
		}()
		return handle(ctx, request)
	}
}
