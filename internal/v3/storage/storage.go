package storage

import "context"

type Client interface {
	Begin(ctx context.Context) (Transaction, error)
}

// type Command interface {
// }

// type Query[R any] interface {
// 	Result() R
// }

// type Executor interface {
// 	Execute(ctx context.Context, command Command) error
// }

// type Querier interface {
// 	Query[R](ctx context.Context, query Query[R]) error
// }

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	// End the transaction based on err. If err is nil the transaction is committed, otherwise it is rolled back.
	End(ctx context.Context, err error) error
	OnCommit(hook func(ctx context.Context) error)
	OnRollback(hook func(ctx context.Context) error)
}
