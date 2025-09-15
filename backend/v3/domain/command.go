package domain

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// Commander is the all it needs to implement the command pattern.
// It is the interface all manipulations need to implement.
// If possible it should also be used for queries. We will find out if this is possible in the future.
type Commander interface {
	Execute(ctx context.Context, opts *CommandOpts) (err error)
	Validate() (err error)
	fmt.Stringer
}

// Invoker is part of the command pattern.
// It is the interface that is used to execute commands.
type Invoker interface {
	Invoke(ctx context.Context, command Commander, opts *CommandOpts) error
}

// CommandOpts are passed to each command
// they provide common fields used by commands like the database client.
type CommandOpts struct {
	DB       database.QueryExecutor
	Invoker  Invoker
	_orgRepo func(db database.QueryExecutor) OrganizationRepository
}

func (opts *CommandOpts) SetOrgRepo(repo func(db database.QueryExecutor) OrganizationRepository) {
	opts._orgRepo = repo
}

func (opts *CommandOpts) orgRepo() OrganizationRepository {
	if opts._orgRepo != nil {
		return opts._orgRepo(opts.DB)
	}
	return orgRepo(opts.DB)
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
		DB:       pool,
		Invoker:  invoker,
		_orgRepo: orgRepo,
	}
}

// commandBatch is a batch of commands.
// It uses the [Invoker] provided by the opts to execute each command.
type commandBatch struct {
	Commands []Commander
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

func (b *commandBatch) Validate() (err error) {
	for _, cmd := range b.Commands {
		if err = cmd.Validate(); err != nil {
			return err
		}
	}
	return nil
}

var _ Commander = (*commandBatch)(nil)
