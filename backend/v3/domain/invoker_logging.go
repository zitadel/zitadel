package domain

import (
	"context"
	"time"
)

// loggingInvoker decorates each command with logging.
// It is an example implementation and logs the command name at the beginning and success or failure after the command got executed.
type loggingInvoker struct {
	next Invoker
}

func newLoggingInvoker(next Invoker) *loggingInvoker {
	return &loggingInvoker{next: next}
}

func (i *loggingInvoker) Invoke(ctx context.Context, executor Executor, opts *InvokeOpts) (err error) {
	start := time.Now()
	logger.InfoContext(ctx, "invoking", "name", executor.String())

	if i.next != nil {
		err = i.next.Invoke(ctx, executor, opts)
	} else {
		err = executor.Execute(ctx, opts)
	}

	if err != nil {
		logger.ErrorContext(ctx, "invocation failed", "name", executor.String(), "error", err, "took", time.Since(start))
		return err
	}
	logger.InfoContext(ctx, "invocation succeeded", "name", executor.String(), "took", time.Since(start))
	return nil
}
