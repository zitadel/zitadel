package domain

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/tracing"
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
	ctx, span := tracing.NewNamedSpan(ctx, fmt.Sprintf("%T", executor))
	defer func() {
		span.EndWithError(err)
	}()

	return i.execute(ctx, executor, opts)
}
