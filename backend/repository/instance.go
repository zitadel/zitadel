package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/handler"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/telemetry/logging"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

type Instance struct {
	ID   string
	Name string
}

type InstanceOptions struct {
	cache *InstanceCache
}

type instance struct {
	options[InstanceOptions]
	*InstanceOptions
}

func NewInstance(opts ...Option[InstanceOptions]) *instance {
	i := new(instance)
	i.InstanceOptions = &i.options.custom

	for _, opt := range opts {
		opt.apply(&i.options)
	}
	return i
}

func WithInstanceCache(c *InstanceCache) Option[InstanceOptions] {
	return func(opts *options[InstanceOptions]) {
		opts.custom.cache = c
	}
}

func (i *instance) Create(ctx context.Context, tx database.Transaction, instance *Instance) (*Instance, error) {
	return tracing.Wrap(i.tracer, "instance.SetUp",
		handler.Chains(
			handler.Decorates(
				execute(tx).CreateInstance,
				tracing.Decorate[*Instance, *Instance](i.tracer, tracing.WithSpanName("instance.sql.SetUp")),
				logging.Decorate[*Instance, *Instance](i.logger, "instance.sql.SetUp"),
			),
			handler.Decorates(
				events(tx).CreateInstance,
				tracing.Decorate[*Instance, *Instance](i.tracer, tracing.WithSpanName("instance.event.SetUp")),
				logging.Decorate[*Instance, *Instance](i.logger, "instance.event.SetUp"),
			),
			handler.SkipReturnPreviousHandler(i.cache,
				handler.Decorates(
					handler.NoReturnToHandle(i.cache.Set),
					tracing.Decorate[*Instance, *Instance](i.tracer, tracing.WithSpanName("instance.cache.SetUp")),
					logging.Decorate[*Instance, *Instance](i.logger, "instance.cache.SetUp"),
				),
			),
		),
	)(ctx, instance)
}

func (i *instance) ByID(ctx context.Context, querier database.Querier, id string) (*Instance, error) {
	return tracing.Wrap(i.tracer, "instance.byID",
		handler.SkipNext(
			handler.SkipNilHandler(i.cache,
				handler.ResFuncToHandle(i.cache.ByID),
			),
			handler.Chain(
				handler.Decorates(
					query(querier).InstanceByID,
					tracing.Decorate[string, *Instance](i.tracer, tracing.WithSpanName("instance.sql.ByID")),
					logging.Decorate[string, *Instance](i.logger, "instance.sql.ByID"),
				),
				handler.SkipNilHandler(i.cache, handler.NoReturnToHandle(i.cache.Set)),
			),
		),
	)(ctx, id)
}

func (i *instance) ByDomain(ctx context.Context, querier database.Querier, domain string) (*Instance, error) {
	return tracing.Wrap(i.tracer, "instance.byDomain",
		handler.SkipNext(
			handler.SkipNilHandler(i.cache,
				handler.ResFuncToHandle(i.cache.ByDomain),
			),
			handler.Chain(
				handler.Decorate(
					query(querier).InstanceByDomain,
					tracing.Decorate[string, *Instance](i.tracer, tracing.WithSpanName("instance.sql.ByDomain")),
				),
				handler.SkipNilHandler(i.cache, handler.NoReturnToHandle(i.cache.Set)),
			),
		),
	)(ctx, domain)
}

type ListRequest struct {
	Limit uint16
}

func (i *instance) List(ctx context.Context, querier database.Querier, request *ListRequest) ([]*Instance, error) {
	return tracing.Wrap(i.tracer, "instance.list",
		handler.Chains(
			handler.Decorates(
				query(querier).ListInstances,
				tracing.Decorate[*ListRequest, []*Instance](i.tracer, tracing.WithSpanName("instance.sql.List")),
				logging.Decorate[*ListRequest, []*Instance](i.logger, "instance.sql.List"),
			),
		),
	)(ctx, request)
}
