package database

import (
	"context"
	"io/fs"
)

type Row interface {
	Scan(dest ...any) error
}

type Rows interface {
	Row
	Next() bool
	Close() error
	Err() error
}

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	End(ctx context.Context, err error) error

	QueryExecutor
}

type Client interface {
	Beginner
	QueryExecutor

	Release(ctx context.Context) error
}

type Pool interface {
	Beginner
	QueryExecutor

	Acquire(ctx context.Context) (Client, error)
	Close(ctx context.Context) error
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

type Beginner interface {
	Begin(ctx context.Context, opts *TransactionOptions) (Transaction, error)
}

type QueryExecutor interface {
	Querier
	Executor
}

type Querier interface {
	Query(ctx context.Context, stmt string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, stmt string, args ...any) Row
}

func Query[Out any](q Querier, fn func(q Querier) ([]Out, error)) ([]Out, error) {
	return fn(q)
}

func QueryRow[Out any](q Querier, fn func(q Querier) (Out, error)) (Out, error) {
	return fn(q)
}

type Executor interface {
	Exec(ctx context.Context, stmt string, args ...any) error
}

// LoadStatements sets the sql statements strings
// TODO: implement
func LoadStatements(fs.FS) error {
	return nil
}
