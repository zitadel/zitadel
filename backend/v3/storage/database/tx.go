package database

import "context"

// Transaction is an SQL transaction.
type Transaction interface {
	// Commit marks a transaction as successful and commits it to the database.
	Commit(ctx context.Context) error
	// Rollback undoes all changes made in the transaction.
	Rollback(ctx context.Context) error
	// End commits the transaction if err is nil, otherwise rollbacks.
	//
	// If err is nil the returned error is from Commit.
	// If err is not nil, and rollback is successful, the original err is returned.
	// If err is not nil, and rollback fails, the two errors are joined using [errors.Join].
	End(ctx context.Context, err error) error

	Begin(ctx context.Context) (Transaction, error)

	QueryExecutor
}

// Beginner can start a new transaction.
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
