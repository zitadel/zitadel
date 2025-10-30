package domain

import (
	"context"
	"fmt"
)

// traceInvoker decorates each command with tracing.
type traceInvoker struct {
	invoker
}

func NewTraceInvoker(next Invoker) *traceInvoker {
	return &traceInvoker{
		invoker: invoker{next: next},
	}
}

func (i *traceInvoker) Invoke(ctx context.Context, executor Executor, opts *InvokeOpts) (err error) {
	ctx, span := tracer.Start(ctx, fmt.Sprintf("%T", executor))
	defer func() {
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}()

	return i.execute(ctx, executor, opts)
}
