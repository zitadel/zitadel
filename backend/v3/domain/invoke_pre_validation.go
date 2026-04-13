package domain

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/tracing"
)

type preValidationInvoker struct {
	invoker
}

// NewPreValidationInvoker creates a new [preValidationInvoker].
// It is responsible for calling the [PreValidator.PreValidate] function before event collection and transaction handling.
func NewPreValidationInvoker(next Invoker) *preValidationInvoker {
	return &preValidationInvoker{
		invoker: invoker{next: next},
	}
}

type PreValidator interface {
	PreValidate(ctx context.Context, opts *InvokeOpts) (err error)
}

func (i *preValidationInvoker) Invoke(ctx context.Context, executor Executor, opts *InvokeOpts) error {
	if err := i.executePreValidation(ctx, executor, opts); err != nil {
		return err
	}
	return i.execute(ctx, executor, opts)
}

func (i *preValidationInvoker) executePreValidation(ctx context.Context, executor Executor, opts *InvokeOpts) (err error) {
	if validator, ok := executor.(PreValidator); ok {
		ctx, span := tracing.NewNamedSpan(ctx, fmt.Sprintf("%s.PreValidate", executor.String()))
		defer func() {
			span.EndWithError(err)
		}()
		return validator.PreValidate(ctx, opts)
	}
	if batch, ok := executor.(*batchExecutor); ok {
		for _, sub := range batch.executors {
			if err = i.executePreValidation(ctx, sub, opts); err != nil {
				return err
			}
		}
	}
	return nil
}
