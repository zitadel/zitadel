package domain

import "context"

type noopInvoker struct {
	next Invoker
}

func (i *noopInvoker) Invoke(ctx context.Context, command Executor, opts *InvokeOpts) error {
	if i.next != nil {
		return i.next.Invoke(ctx, command, opts)
	}
	return command.Execute(ctx, opts)
}
