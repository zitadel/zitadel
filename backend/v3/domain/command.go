package domain

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	legacy_es "github.com/zitadel/zitadel/internal/eventstore"
)

// Commander is the all it needs to implement the command pattern.
// It is the interface all manipulations need to implement.
// If possible it should also be used for queries. We will find out if this is possible in the future.
type Commander interface {
	Execute(ctx context.Context, opts *CommandOpts) (err error)
	// Events returns the events that should be pushed to the event store after the command is executed.
	// If the command does not produce events, it should return nil or an empty slice.
	Events() []legacy_es.Command
	fmt.Stringer
}

// Invoker is part of the command pattern.
// It is the interface that is used to execute commands.
type Invoker interface {
	Invoke(ctx context.Context, command Commander, opts *CommandOpts) error
}

// CommandOpts are passed to each command
// the provide common fields used by commands like the database client.
type CommandOpts struct {
	DB      database.QueryExecutor
	Invoker Invoker
}

type ensureTxOpts struct {
	*database.TransactionOptions
}

type EnsureTransactionOpt func(*ensureTxOpts)

// EnsureTx ensures that the DB is a transaction. If it is not, it will start a new transaction.
// The returned close function will end the transaction. If the DB is already a transaction, the close function
// will do nothing because another [Commander] is already responsible for ending the transaction.
func (o *CommandOpts) EnsureTx(ctx context.Context, opts ...EnsureTransactionOpt) (close func(context.Context, error) error, err error) {
	beginner, ok := o.DB.(database.Beginner)
	if !ok {
		// db is already a transaction
		return func(_ context.Context, err error) error {
			return err
		}, nil
	}

	txOpts := &ensureTxOpts{
		TransactionOptions: new(database.TransactionOptions),
	}
	for _, opt := range opts {
		opt(txOpts)
	}

	tx, err := beginner.Begin(ctx, txOpts.TransactionOptions)
	if err != nil {
		return nil, err
	}
	o.DB = tx

	return func(ctx context.Context, err error) error {
		return tx.End(ctx, err)
	}, nil
}

// EnsureClient ensures that the o.DB is a client. If it is not, it will get a new client from the [database.Pool].
// The returned close function will release the client. If the o.DB is already a client or transaction, the close function
// will do nothing because another [Commander] is already responsible for releasing the client.
func (o *CommandOpts) EnsureClient(ctx context.Context) (close func(_ context.Context) error, err error) {
	pool, ok := o.DB.(database.Pool)
	if !ok {
		// o.DB is already a client
		return func(_ context.Context) error {
			return nil
		}, nil
	}
	client, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	o.DB = client
	return func(ctx context.Context) error {
		return client.Release(ctx)
	}, nil
}

func (o *CommandOpts) Invoke(ctx context.Context, command Commander) error {
	if o.Invoker == nil {
		return command.Execute(ctx, o)
	}
	return o.Invoker.Invoke(ctx, command, o)
}

func DefaultOpts(invoker Invoker) *CommandOpts {
	if invoker == nil {
		invoker = &noopInvoker{}
	}
	return &CommandOpts{
		DB:      pool,
		Invoker: invoker,
	}
}

// commandBatch is a batch of commands.
// It uses the [Invoker] provided by the opts to execute each command.
type commandBatch struct {
	Commands []Commander
}

// Events implements Commander.
func (cmd *commandBatch) Events() []legacy_es.Command {
	commands := make([]legacy_es.Command, 0, len(cmd.Commands))
	for _, c := range cmd.Commands {
		if e := c.Events(); len(e) > 0 {
			commands = append(commands, e...)
		}
	}
	return commands
}

func BatchCommands(cmds ...Commander) *commandBatch {
	return &commandBatch{
		Commands: cmds,
	}
}

// String implements [Commander].
func (cmd *commandBatch) String() string {
	return "commandBatch"
}

func (b *commandBatch) Execute(ctx context.Context, opts *CommandOpts) (err error) {
	for _, cmd := range b.Commands {
		if err = opts.Invoke(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

var _ Commander = (*commandBatch)(nil)
