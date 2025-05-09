package database

import "context"

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	End(ctx context.Context, err error) error

	Begin(ctx context.Context) (Transaction, error)

	QueryExecutor
}

type Beginner interface {
	Begin(ctx context.Context, opts *TransactionOptions) (Transaction, error)
}

type TransactionOptions struct {
	IsolationLevel IsolationLevel
	AccessMode     AccessMode
}

type IsolationLevel uint8

const (
	IsolationLevelSerializable IsolationLevel = iota
	IsolationLevelReadCommitted
)

type AccessMode uint8

const (
	AccessModeReadWrite AccessMode = iota
	AccessModeReadOnly
)
