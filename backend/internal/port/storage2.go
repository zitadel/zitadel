package port

import (
	"context"
	"slices"
)

type Getter[T any] interface {
	Get(ctx context.Context, filters []Filter) (T, error)
}

func Get[T Object](ctx context.Context, get func(ctx context.Context, filters []Filter) (T, error), filters []Filter) (T, error) {
	return get(ctx, filters)
}

type Lister[T any] interface {
	List(ctx context.Context, filters []Filter) ([]T, error)
}

func List[T Object](ctx context.Context, lister Lister[T], filters []Filter) ([]T, error) {
	return lister.List(ctx, filters)
}

type tx struct{}

type instance struct{ id string }

func (i instance) Columns() []*Column {
	return []*Column{
		{Name: "id", Value: i.id},
	}
}

type instanceRepo struct {
	instances []*instance
}

type instanceTx struct {
	tx
	instances []*instance
}

func (ir *instanceRepo) ForTx(t tx) *instanceTx {
	return &instanceTx{
		tx:        t,
		instances: slices.Clone(ir.instances),
	}
}

var _ Querier[*instance] = (*instanceTx)(nil)

func (it *instanceTx) Get(ctx context.Context, filters []Filter) (*instance, error) {
	return nil, nil
}

func (it *instanceTx) List(ctx context.Context, filters []Filter) ([]*instance, error) {
	return it.instances, nil
}

type instanceSQLRepo struct{}

func (ir *instanceSQLRepo) ForTx(t tx) *instanceSQLTx {
	return &instanceSQLTx{tx: t}
}

type instanceSQLTx struct {
	tx
}

func bla() {
	var ir instanceRepo
	it := ir.ForTx(tx{})

	_, _ = Get(context.Background(), it.Get, nil)

}

type Getter2[T Object, C Client3] interface {
	Execute(ctx context.Context, client C) error
	Result() T
}

type Executor2[C Client3] interface {
	Execute(ctx context.Context, client C) error
}

type Client3 interface {
	Exec(Executor2[Client3]) error
}
