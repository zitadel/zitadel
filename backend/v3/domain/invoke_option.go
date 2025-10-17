package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type InvokeOpt func(*InvokeOpts)

func WithOrganizationRepo(repo OrganizationRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.organizationRepo = repo
	}
}

func WithOrganizationDomainRepo(repo OrganizationDomainRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.organizationDomainRepo = repo
	}
}

func WithProjectRepo(repo ProjectRepository) InvokeOpt {
	return func(opts *InvokeOpts) {
		opts.projectRepo = repo
	}
}

// InvokeOpts are passed to each command
// they provide common fields used by commands like the database client.
type InvokeOpts struct {
	DB                     database.QueryExecutor
	Invoker                Invoker
	Permissions            PermissionChecker
	organizationRepo       OrganizationRepository
	organizationDomainRepo OrganizationDomainRepository
	projectRepo            ProjectRepository
}

type ensureTxOpts struct {
	*database.TransactionOptions
}

type EnsureTransactionOpt func(*ensureTxOpts)

// EnsureTx ensures that the DB is a transaction. If it is not, it will start a new transaction.
// The returned close function will end the transaction. If the DB is already a transaction, the close function
// will do nothing because another [Commander] is already responsible for ending the transaction.
func (o *InvokeOpts) EnsureTx(ctx context.Context, opts ...EnsureTransactionOpt) (close func(context.Context, error) error, err error) {
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
func (o *InvokeOpts) EnsureClient(ctx context.Context) (close func(_ context.Context) error, err error) {
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

func (o *InvokeOpts) Invoke(ctx context.Context, executor Executor) error {
	if o.Invoker == nil {
		return executor.Execute(ctx, o)
	}
	return o.Invoker.Invoke(ctx, executor, o)
}

func DefaultOpts(invoker Invoker) *InvokeOpts {
	if invoker == nil {
		invoker = &noopInvoker{}
	}
	return &InvokeOpts{
		DB:          pool,
		Invoker:     invoker,
		Permissions: &noopPermissionChecker{}, // prevent panics for now
	}
}
