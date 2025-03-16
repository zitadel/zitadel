package orchestrate

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/repository/event"
	"github.com/zitadel/zitadel/backend/repository/orchestrate/handler"
	"github.com/zitadel/zitadel/backend/repository/sql"
	"github.com/zitadel/zitadel/backend/repository/telemetry/logged"
	"github.com/zitadel/zitadel/backend/repository/telemetry/traced"
	"github.com/zitadel/zitadel/backend/storage/cache"
	"github.com/zitadel/zitadel/backend/storage/cache/connector/noop"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

type InstanceOptions struct {
	cache cache.Cache[repository.InstanceIndex, string, *repository.Instance]
}

type instance struct {
	options[InstanceOptions]
	*InstanceOptions
}

func Instance(opts ...Option[InstanceOptions]) *instance {
	i := new(instance)
	i.InstanceOptions = &i.options.custom
	i.cache = noop.NewCache[repository.InstanceIndex, string, *repository.Instance]()

	for _, opt := range opts {
		opt.apply(&i.options)
	}
	return i
}

func WithInstanceCache(cache cache.Cache[repository.InstanceIndex, string, *repository.Instance]) Option[InstanceOptions] {
	return func(opts *options[InstanceOptions]) {
		opts.custom.cache = cache
	}
}

func (i *instance) Create(ctx context.Context, tx database.Transaction, instance *repository.Instance) (*repository.Instance, error) {
	return traced.Wrap(i.tracer, "instance.SetUp",
		handler.Chains(
			handler.Decorates(
				sql.Execute(tx).CreateInstance,
				traced.Decorate[*repository.Instance, *repository.Instance](i.tracer, tracing.WithSpanName("instance.sql.SetUp")),
				logged.Decorate[*repository.Instance, *repository.Instance](i.logger, "instance.sql.SetUp"),
			),
			handler.Decorates(
				event.Store(tx).CreateInstance,
				traced.Decorate[*repository.Instance, *repository.Instance](i.tracer, tracing.WithSpanName("instance.event.SetUp")),
				logged.Decorate[*repository.Instance, *repository.Instance](i.logger, "instance.event.SetUp"),
			),
			handler.SkipNilHandler(i.cache,
				handler.Decorates(
					handler.NoReturnToHandle(i.cache.Set),
					traced.Decorate[*repository.Instance, *repository.Instance](i.tracer, tracing.WithSpanName("instance.cache.SetUp")),
					logged.Decorate[*repository.Instance, *repository.Instance](i.logger, "instance.cache.SetUp"),
				),
			),
		),
	)(ctx, instance)
}

func (i *instance) ByID(ctx context.Context, querier database.Querier, id string) (*repository.Instance, error) {
	return handler.SkipNext(
		handler.CacheGetToHandle(i.cache.Get, repository.InstanceByID),
		handler.Chain(
			handler.Decorate(
				sql.Query(querier).InstanceByID,
				traced.Decorate[string, *repository.Instance](i.tracer, tracing.WithSpanName("instance.sql.ByID")),
			),
			handler.SkipNilHandler(i.cache, handler.NoReturnToHandle(i.cache.Set)),
		),
	)(ctx, id)
}

func (i *instance) ByDomain(ctx context.Context, querier database.Querier, domain string) (*repository.Instance, error) {
	return handler.SkipNext(
		handler.CacheGetToHandle(i.cache.Get, repository.InstanceByDomain),
		handler.Chain(
			handler.Decorate(
				sql.Query(querier).InstanceByDomain,
				traced.Decorate[string, *repository.Instance](i.tracer, tracing.WithSpanName("instance.sql.ByDomain")),
			),
			handler.SkipNilHandler(i.cache, handler.NoReturnToHandle(i.cache.Set)),
		),
	)(ctx, domain)
}
