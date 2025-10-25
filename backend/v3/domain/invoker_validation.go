package domain

import (
	"context"
)

type validatorInvoker struct {
	invoker
}

func NewValidatorInvoker(next Invoker) *validatorInvoker {
	return &validatorInvoker{
		invoker: invoker{next: next},
	}
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

	return i.execute(ctx, executor, opts)
}
