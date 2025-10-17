package domain

import (
	"context"
	"fmt"
	"strings"
)

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
// If you want to invoke multiple commands in a single transaction, you can use the [commandBatch].
func Invoke(ctx context.Context, executor Executor, opts ...InvokeOpt) error {
	invokeOpts := &InvokeOpts{
		Invoker: newEventStoreInvoker(
			newLoggingInvoker(
				newTraceInvoker(
					newValidatorInvoker(nil),
				),
			),
		),
		DB: pool,
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

type noopInvoker struct {
	next Invoker
}

func (i *noopInvoker) Invoke(ctx context.Context, command Executor, opts *InvokeOpts) error {
	if i.next != nil {
		return i.next.Invoke(ctx, command, opts)
	}
	return command.Execute(ctx, opts)
}

// executorBatch is a batch of commands.
// It uses the [Invoker] provided by the opts to execute each command.
type executorBatch struct {
	Commands []Executor
}

func BatchExecutors(cmds ...Executor) *executorBatch {
	return &executorBatch{
		Commands: cmds,
	}
}

// String implements [Commander].
func (cmd *executorBatch) String() string {
	names := make([]string, len(cmd.Commands))
	for i, c := range cmd.Commands {
		names[i] = c.String()
	}
	return fmt.Sprintf("commandBatch[%s]", strings.Join(names, ", "))
}

func (b *executorBatch) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	for _, cmd := range b.Commands {
		if err = opts.Invoke(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

var _ Executor = (*executorBatch)(nil)
