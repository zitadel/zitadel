package domain

import (
	"context"
)

// transactionInvoker ensures that [InvokeOpts].DB() returns a [database.Transaction].
// if a new transaction is started, it will be committed or rolled back after the [Executor] execution.
type transactionInvoker struct {
	invoker
}

// NewTransactionInvoker creates a new [transactionInvoker].
func NewTransactionInvoker(next Invoker) *transactionInvoker {
	return &transactionInvoker{
		invoker: invoker{next: next},
	}
}

//go:generate mockgen -destination=mock/transactional.mock.go -package=domainmock . Transactional
type Transactional interface {
	RequiresTransaction()
}

// Invoke implements [Invoker].
func (i *transactionInvoker) Invoke(ctx context.Context, executor Executor, opts *InvokeOpts) (err error) {
	if _, ok := executor.(Transactional); !ok {
		return i.execute(ctx, executor, opts)
	}

	close, err := i.ensureTx(ctx, opts)
	if err != nil {
		return err
	}
	defer func() { err = close(err) }()

	return i.execute(ctx, executor, opts)
}
