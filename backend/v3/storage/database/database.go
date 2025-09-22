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

	Ping(ctx context.Context) error
}

type PoolTest interface {
	Pool
	// MigrateTest is the same as [Migrator] but executes the migrations multiple times instead of only once.
	MigrateTest(ctx context.Context) error
}

// Client is a single database connection which can be released back to the pool.
type Client interface {
	Beginner
	QueryExecutor
	Migrator

	Release(ctx context.Context) error

	Ping(ctx context.Context) error
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
	Scanner
	Next() bool
	Close() error
	Err() error
}

type CollectableRows interface {
	// Collect collects all rows and scans them into dest.
	// dest must be a pointer to a slice of pointer to structs
	// e.g. *[]*MyStruct
	// Rows are closed after this call.
	Collect(dest any) error
	// CollectFirst collects the first row and scans it into dest.
	// dest must be a pointer to a struct
	// e.g. *MyStruct{}
	// Rows are closed after this call.
	CollectFirst(dest any) error
	// CollectExactlyOneRow collects exactly one row and scans it into dest.
	// e.g. *MyStruct{}
	// Rows are closed after this call.
	CollectExactlyOneRow(dest any) error
}
