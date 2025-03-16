package orchestrate

import (
	"context"
	"fmt"

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

func Instance(opts ...InstanceConfig) *instance {
	i := new(instance)
	for _, opt := range opts {
		opt.applyInstance(i)
	}
	return i
}

func WithInstanceCache(cache *cache.Instance) instanceOption {
	return func(i *instance) {
		i.cache = cache
	}
}

type InstanceConfig interface {
	applyInstance(*instance)
}

// instanceOption applies an option to the instance.
type instanceOption func(*instance)

func (io instanceOption) applyInstance(i *instance) {
	io(i)
}

func (o Option) applyInstance(i *instance) {
	o(&i.options)
}

func (i *instance) Create(ctx context.Context, tx database.Transaction, instance *repository.Instance) (*repository.Instance, error) {
	fmt.Println("----------------")
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
					i.cache.Set,
					traced.Decorate[*repository.Instance, *repository.Instance](i.tracer, tracing.WithSpanName("instance.cache.SetUp")),
					logged.Decorate[*repository.Instance, *repository.Instance](i.logger, "instance.cache.SetUp"),
				),
			),
		),
	)(ctx, instance)
}

func (i *instance) ByID(ctx context.Context, querier database.Querier, id string) (*repository.Instance, error) {
	return handler.SkipNext(
		i.cache.ByID,
		handler.Chain(
			handler.Decorate(
				sql.Query(querier).InstanceByID,
				traced.Decorate[string, *repository.Instance](i.tracer, tracing.WithSpanName("instance.sql.ByID")),
			),
			handler.SkipNilHandler(i.cache, i.cache.Set),
		),
	)(ctx, id)
}

func (i *instance) ByDomain(ctx context.Context, querier database.Querier, domain string) (*repository.Instance, error) {
	return handler.SkipNext(
		i.cache.ByDomain,
		handler.Chain(
			handler.Decorate(
				sql.Query(querier).InstanceByDomain,
				traced.Decorate[string, *repository.Instance](i.tracer, tracing.WithSpanName("instance.sql.ByDomain")),
			),
			handler.SkipNilHandler(i.cache, i.cache.Set),
		),
	)(ctx, domain)
}
