package port

import "context"

type Operation uint8

const (
	OperationEqual Operation = iota
)

type Object interface {
	Columns() []*Column
}

type Column struct {
	Name  string
	Value any
}

type Filter interface {
	Column() *Column
	Operation() Operation
}

var _ Filter = (*filter)(nil)

type filter struct {
	column *Column
	op     Operation
}

func newFilter(column *Column, op Operation) Filter {
	return &filter{column: column, op: op}
}

func (f *filter) Column() *Column {
	return f.column
}

func (f *filter) Operation() Operation {
	return f.op
}

func NewEqualFilter(column *Column) Filter {
	return newFilter(column, OperationEqual)
}

type Querier[T any] interface {
	Get(ctx context.Context, filters []Filter) (T, error)
	List(ctx context.Context, filters []Filter) ([]T, error)
}

type Executor[T Object] interface {
	Create(ctx context.Context, object T) error
	Update(ctx context.Context, columns []*Column, filters []Filter) error
	Delete(ctx context.Context, filters []Filter) error
}

type Pool[T Object] interface {
	Acquire(ctx context.Context) (Client[T], error)
	Begin(ctx context.Context) (Transaction[T], error)
}

type Client[T Object] interface {
	Querier[T]
	Executor[T]
	Begin(ctx context.Context) (Transaction[T], error)
	Release(ctx context.Context) error
}

type Transaction[T Object] interface {
	Executor[T]
	Querier[T]
	End(ctx context.Context, gotErr error) error
}
