package domain

import "context"

type noopInvoker struct {
	invoker
}

func (i *noopInvoker) Invoke(ctx context.Context, executor Executor, opts *InvokeOpts) error {
	return i.execute(ctx, executor, opts)
}
