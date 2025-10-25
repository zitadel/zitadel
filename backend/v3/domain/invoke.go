package domain

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate mockgen -typed -package domainmock -destination ./mock/executor.mock.go . Executor
type Executor interface {
	Execute(ctx context.Context, opts *InvokeOpts) (err error)
	fmt.Stringer
}

// Commander is all that is needed to implement the command pattern.
// It is the interface all manipulations need to implement.
// If possible it should also be used for queries. We will find out if this is possible in the future.
//
//go:generate mockgen -typed -package domainmock -destination ./mock/commander.mock.go . Commander
type Commander interface {
	Validator
	EventProducer
	Executor
}

// Querier used to query data.
//
//go:generate mockgen -typed -package domainmock -destination ./mock/querier.mock.go . Querier
type Querier[T any] interface {
	Validator
	Executor
	// Result returns the result of the query.
	// If `Execute` returns an error, the result is nil.
	Result() T
}

// Invoke provides a way to execute commands within the domain package.
// It uses a chain of responsibility pattern to handle the command execution.
// The default chain includes logging, tracing, and event publishing.
// If you want to invoke multiple commands in a single transaction, you can use the [batchExecutor].
func Invoke(ctx context.Context, executor Executor, opts ...InvokeOpt) error {
	invokeOpts := &InvokeOpts{
		Invoker: NewLoggingInvoker(
			NewTraceInvoker(
				NewEventStoreInvoker(
					NewTransactionInvoker(
						NewValidatorInvoker(nil),
					),
				),
			),
		),
	}
	for _, opt := range opts {
		opt(invokeOpts)
	}
	return invokeOpts.Invoke(ctx, executor)
}

// Invoker is part of the command pattern.
// It is the interface that is used to execute commands.
type Invoker interface {
	Invoke(ctx context.Context, command Executor, opts *InvokeOpts) error
}

type invoker struct {
	next Invoker
}

func (i invoker) execute(ctx context.Context, executor Executor, opts *InvokeOpts) error {
	if i.next != nil {
		return i.next.Invoke(ctx, executor, opts)
	}
	return executor.Execute(ctx, opts)
}

// ensureTx ensures that the InvokeOpts has a transaction.
// the close function ends the transaction and resets [InvokeOpts].db.
// If a new transaction is started, it returns a function to end the transaction.
// The caller is responsible to call the returned function to end the transaction.
// If no new transaction is started, the returned function is a no-op and still safe to call.
func (i invoker) ensureTx(ctx context.Context, opts *InvokeOpts) (close func(err error) error, err error) {
	beginner, ok := opts.DB().(database.Beginner)
	if !ok {
		return func(err error) error { return err }, nil
	}

	previousDB := opts.DB()
	tx, err := beginner.Begin(ctx, nil)
	if err != nil {
		return nil, err
	}
	opts.db = tx
	return func(err error) error {
		opts.db = previousDB
		return tx.End(ctx, err)
	}, nil
}
