package domain

import (
	"context"
	"fmt"
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
		Invoker: newEventStoreInvoker(
			newLoggingInvoker(
				newTraceInvoker(
					newValidatorInvoker(nil),
				),
			),
		),
		DB:          pool,
		Permissions: &noopPermissionChecker{},
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
