package query

import "context"

type Query[T any] interface {
	Execute(ctx context.Context) (T, error)
	Name() string
}
