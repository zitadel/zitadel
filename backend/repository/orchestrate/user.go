package orchestrate

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/repository/cached"
	"github.com/zitadel/zitadel/backend/repository/event"
	"github.com/zitadel/zitadel/backend/repository/orchestrate/handler"
	"github.com/zitadel/zitadel/backend/repository/sql"
	"github.com/zitadel/zitadel/backend/repository/telemetry/traced"
	"github.com/zitadel/zitadel/backend/storage/cache"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

type UserOptions struct {
	cache *cached.User
}

type user struct {
	options[UserOptions]
	*UserOptions
}

func User(opts ...Option[UserOptions]) *user {
	i := new(user)
	i.UserOptions = &i.options.custom

	for _, opt := range opts {
		opt(&i.options)
	}
	return i
}

func WithUserCache(cache cache.Cache[repository.UserIndex, string, *repository.User]) Option[UserOptions] {
	return func(i *options[UserOptions]) {
		i.custom.cache = cached.NewUser(cache)
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
		handler.SkipNilHandler(i.cache,
			handler.ResFuncToHandle(i.cache.ByID),
		),
		handler.Chain(
			handler.Decorate(
				sql.Query(querier).UserByID,
				traced.Decorate[string, *repository.User](i.tracer, tracing.WithSpanName("user.sql.ByID")),
			),
			handler.SkipNilHandler(i.custom.cache, handler.NoReturnToHandle(i.cache.Set)),
		),
	)(ctx, id)
}
