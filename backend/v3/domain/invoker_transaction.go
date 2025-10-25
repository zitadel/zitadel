package domain

import (
	"context"
)

// transactionInvoker ensures that [InvokeOpts].DB is a [database.Transaction].
// if a new transaction is started, it will be committed or rolled back after the command execution.
type transactionInvoker struct {
	invoker
}

func NewTransactionInvoker(next Invoker) *transactionInvoker {
	return &transactionInvoker{
		invoker: invoker{next: next},
	}
}

type Transactional interface {
	RequiresTransaction() bool
}

// Invoke implements [Invoker].
func (i *transactionInvoker) Invoke(ctx context.Context, executor Executor, opts *InvokeOpts) (err error) {
	if transactional, ok := executor.(Transactional); !ok || !transactional.RequiresTransaction() {
		return i.execute(ctx, executor, opts)
	}

	close, err := i.ensureTx(ctx, opts)
	if err != nil {
		return err
	}
	defer func() { err = close(err) }()

	return i.execute(ctx, executor, opts)
}
