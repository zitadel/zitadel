package domain

import (
	"context"
)

type validatorInvoker struct {
	next Invoker
}

func newValidatorInvoker(next Invoker) *validatorInvoker {
	return &validatorInvoker{next: next}
}

type Validator interface {
	Validate(ctx context.Context, opts *InvokeOpts) (err error)
}

func (i *validatorInvoker) Invoke(ctx context.Context, executor Executor, opts *InvokeOpts) error {
	if validator, ok := executor.(Validator); ok {
		if err := validator.Validate(ctx, opts); err != nil {
			return err
		}
	}

	if i.next != nil {
		return i.next.Invoke(ctx, executor, opts)
	}

	return executor.Execute(ctx, opts)
}
