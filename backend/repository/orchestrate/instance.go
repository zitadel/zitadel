package orchestrate

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/repository/cache"
	"github.com/zitadel/zitadel/backend/repository/event"
	"github.com/zitadel/zitadel/backend/repository/orchestrate/handler"
	"github.com/zitadel/zitadel/backend/repository/sql"
	"github.com/zitadel/zitadel/backend/repository/telemetry/logged"
	"github.com/zitadel/zitadel/backend/repository/telemetry/traced"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

type instance struct {
	options

	cache *cache.Instance
}

func Instance(opts ...Option) *instance {
	i := new(instance)
	for _, opt := range opts {
		opt(&i.options)
	}
	return i
}

func (i *instance) apply(o Option) {
	o(&i.options)
}

func (i *instance) SetUp(ctx context.Context, tx database.Transaction, instance *repository.Instance) (*repository.Instance, error) {
	return handler.NewChained(
		handler.NewDecorated(
			traced.DecorateHandle[*repository.Instance, *repository.Instance](i.tracer, tracing.WithSpanName("instance.sql.SetUp")),
			sql.SetUpInstance(tx),
		),
		handler.NewChained(
			handler.NewDecorated(
				traced.DecorateHandle[*repository.Instance, *repository.Instance](i.tracer, tracing.WithSpanName("instance.event.SetUp")),
				event.SetUpInstanceWithout(tx),
			),
			handler.NewDecorated(
				traced.DecorateHandle[*repository.Instance, *repository.Instance](i.tracer, tracing.WithSpanName("instance.cache.SetUp")),
				cache.SetUpInstanceWithout(i.cache),
			),
		),
	)(ctx, instance)
}

func (i *instance) ByID(ctx context.Context, querier database.Querier, id string) (*repository.Instance, error) {
	return traced.Wrap(i.tracer, "instance.byID",
		logged.Wrap(i.logger, "instance.byID",
			cache.ForInstanceByID(i.cache,
				sql.InstanceByID(querier),
			),
		),
	)(ctx, id)
}

func (i *instance) ByDomain(ctx context.Context, querier database.Querier, domain string) (*repository.Instance, error) {
	return traced.Wrap(i.tracer, "instance.byDomain",
		logged.Wrap(i.logger, "instance.byDomain",
			cache.ForInstanceByDomain(i.cache,
				sql.InstanceByDomain(querier),
			),
		),
	)(ctx, domain)
}
