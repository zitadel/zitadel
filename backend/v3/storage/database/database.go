package database

import (
	"context"
)

var (
	db *database
)

type database struct {
	connector Connector
	pool      Pool
}

type Pool interface {
	Beginner
	QueryExecutor

	Acquire(ctx context.Context) (Client, error)
	Close(ctx context.Context) error
}

type Client interface {
	Beginner
	QueryExecutor

	Release(ctx context.Context) error
}

type Querier interface {
	Query(ctx context.Context, stmt string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, stmt string, args ...any) Row
}

type Executor interface {
	Exec(ctx context.Context, stmt string, args ...any) error
}

type QueryExecutor interface {
	Querier
	Executor
}

type Scanner interface {
	Scan(dest ...any) error
}

type Row interface {
	Scanner
}

type Rows interface {
	Row
	Next() bool
	Close() error
	Err() error
}

type Query[T any] func(querier Querier) (result T, err error)
