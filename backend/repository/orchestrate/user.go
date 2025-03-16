package orchestrate

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/repository/cache"
	"github.com/zitadel/zitadel/backend/repository/event"
	"github.com/zitadel/zitadel/backend/repository/orchestrate/handler"
	"github.com/zitadel/zitadel/backend/repository/sql"
	"github.com/zitadel/zitadel/backend/repository/telemetry/traced"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

type UserOptions struct {
	cache *cache.User
}

type user struct {
	options[UserOptions]
	*UserOptions
}

func User(opts ...Option[UserOptions]) *user {
	i := user{
		options: newOptions[UserOptions](),
	}
	i.UserOptions = i.options.custom
	for _, opt := range opts {
		opt(&i.options)
	}
	return &i
}

func WithUserCache(cache *cache.User) Option[UserOptions] {
	return func(i *options[UserOptions]) {
		i.custom.cache = cache
	}
}

func (i *user) Create(ctx context.Context, tx database.Transaction, user *repository.User) (*repository.User, error) {
	return traced.Wrap(i.tracer, "user.Create",
		handler.Chain(
			handler.Decorate(
				sql.Execute(tx).CreateUser,
				traced.Decorate[*repository.User, *repository.User](i.tracer, tracing.WithSpanName("user.sql.Create")),
			),
			handler.Decorate(
				event.Store(tx).CreateUser,
				traced.Decorate[*repository.User, *repository.User](i.tracer, tracing.WithSpanName("user.event.Create")),
			),
		),
	)(ctx, user)
}

func (i *user) ByID(ctx context.Context, querier database.Querier, id string) (*repository.User, error) {
	return handler.SkipNext(
		i.custom.cache.ByID,
		handler.Chain(
			handler.Decorate(
				sql.Query(querier).UserByID,
				traced.Decorate[string, *repository.User](i.tracer, tracing.WithSpanName("user.sql.ByID")),
			),
			handler.SkipNilHandler(i.custom.cache, i.custom.cache.Set),
		),
	)(ctx, id)
}
