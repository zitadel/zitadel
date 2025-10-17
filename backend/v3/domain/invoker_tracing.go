package domain

import (
	"context"
	"fmt"
)

// traceInvoker decorates each command with tracing.
type traceInvoker struct {
	next Invoker
}

func newTraceInvoker(next Invoker) *traceInvoker {
	return &traceInvoker{next: next}
}

func (i *traceInvoker) Invoke(ctx context.Context, executor Executor, opts *InvokeOpts) (err error) {
	ctx, span := tracer.Start(ctx, fmt.Sprintf("%T", executor))
	defer func() {
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}()

	if i.next != nil {
		return i.next.Invoke(ctx, executor, opts)
	}
	return executor.Execute(ctx, opts)
}
