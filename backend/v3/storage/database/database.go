package database

import (
	"context"
)

// Pool is a connection pool. e.g. pgxpool
type Pool interface {
	Beginner
	QueryExecutor
	Migrator

	Acquire(ctx context.Context) (Client, error)
	Close(ctx context.Context) error
}

type PoolTest interface {
	Pool
	MigrateTest(ctx context.Context) error
}

// Client is a single database connection which can be released back to the pool.
type Client interface {
	Beginner
	QueryExecutor
	Migrator

	Release(ctx context.Context) error
}

// Querier is a database client that can execute queries and return rows.
type Querier interface {
	Query(ctx context.Context, stmt string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, stmt string, args ...any) Row
}

// Executor is a database client that can execute statements.
// It returns the number of rows affected or an error
type Executor interface {
	Exec(ctx context.Context, stmt string, args ...any) (int64, error)
}

// QueryExecutor is a database client that can execute queries and statements.
type QueryExecutor interface {
	Querier
	Executor
}

// Scanner scans a single row of data into the destination.
type Scanner interface {
	Scan(dest ...any) error
}

// Row is an abstraction of sql.Row.
type Row interface {
	Scanner
}

// Rows is an abstraction of sql.Rows.
type Rows interface {
	Row
	Next() bool
	Close() error
	Err() error
}
