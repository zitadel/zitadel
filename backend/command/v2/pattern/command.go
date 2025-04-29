package pattern

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/v2/storage/database"
)

// Command implements the command pattern.
// It is used to encapsulate a request as an object, thereby allowing for parameterization of clients with queues, requests, and operations.
// The command pattern allows for the decoupling of the sender and receiver of a request.
// It is often used in conjunction with the invoker pattern, which is responsible for executing the command.
// The command pattern is a behavioral design pattern that turns a request into a stand-alone object.
// This object contains all the information about the request.
// The command pattern is useful for implementing undo/redo functionality, logging, and queuing requests.
// It is also useful for implementing the macro command pattern, which allows for the execution of a series of commands as a single command.
// The command pattern is also used in event-driven architectures, where events are encapsulated as commands.
type Command interface {
	Execute(ctx context.Context) error
	Name() string
}

type Query[T any] interface {
	Command
	Result() T
}

type Invoker struct{}

// func bla() {
// 	sync.Pool{
// 		New: func() any {
// 			return new(Invoker)
// 		},
// 	}
// }

type Transaction struct {
	beginner database.Beginner
	cmd      Command
	opts     *database.TransactionOptions
}

func (t *Transaction) Execute(ctx context.Context) error {
	tx, err := t.beginner.Begin(ctx, t.opts)
	if err != nil {
		return err
	}
	defer func() { err = tx.End(ctx, err) }()
	return t.cmd.Execute(ctx)
}

func (t *Transaction) Name() string {
	return t.cmd.Name()
}

type batch struct {
	Commands []Command
}

func Batch(cmds ...Command) *batch {
	return &batch{
		Commands: cmds,
	}
}

func (b *batch) Execute(ctx context.Context) error {
	for _, cmd := range b.Commands {
		if err := cmd.Execute(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (b *batch) Name() string {
	return "batch"
}

func (b *batch) Append(cmds ...Command) {
	b.Commands = append(b.Commands, cmds...)
}

type NoopCommand struct{}

func (c *NoopCommand) Execute(_ context.Context) error {
	return nil
}
func (c *NoopCommand) Name() string {
	return "noop"
}

type NoopQuery[T any] struct {
	NoopCommand
}

func (q *NoopQuery[T]) Result() T {
	var zero T
	return zero
}
