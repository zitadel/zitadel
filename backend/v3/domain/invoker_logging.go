package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

// loggingInvoker decorates each command with logging.
// It is an example implementation and logs the command name at the beginning and success or failure after the command got executed.
type loggingInvoker struct {
	invoker
}

func NewLoggingInvoker(next Invoker) *loggingInvoker {
	return &loggingInvoker{
		invoker: invoker{next: next},
	}
}

func (i *loggingInvoker) Invoke(ctx context.Context, executor Executor, opts *InvokeOpts) (err error) {
	start := time.Now()
	logger := logging.FromCtx(ctx)
	logger.InfoContext(ctx, "invoking", "name", executor.String())

	err = i.execute(ctx, executor, opts)
	if err != nil {
		logger.ErrorContext(ctx, "invocation failed", "name", executor.String(), "error", err, "took", time.Since(start))
		return err
	}
	logger.InfoContext(ctx, "invocation succeeded", "name", executor.String(), "took", time.Since(start))
	return nil
}
