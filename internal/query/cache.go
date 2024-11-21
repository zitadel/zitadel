package query

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/cache/connector"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type Caches struct {
	instance cache.Cache[instanceIndex, string, *authzInstance]
	org      cache.Cache[orgIndex, string, *Org]
}

func startCaches(background context.Context, connectors connector.Connectors) (_ *Caches, err error) {
	caches := new(Caches)
	caches.instance, err = connector.StartCache[instanceIndex, string, *authzInstance](background, instanceIndexValues(), cache.PurposeAuthzInstance, connectors.Config.Instance, connectors)
	if err != nil {
		return nil, err
	}
	caches.org, err = connector.StartCache[orgIndex, string, *Org](background, orgIndexValues(), cache.PurposeOrganization, connectors.Config.Organization, connectors)
	if err != nil {
		return nil, err
	}

	caches.registerInstanceInvalidation()
	caches.registerOrgInvalidation()
	return caches, nil
}

type invalidator[I comparable] interface {
	Invalidate(ctx context.Context, index I, key ...string) error
}

func cacheInvalidationFunc[I comparable](cache invalidator[I], index I, getID func(*eventstore.Aggregate) string) func(context.Context, []*eventstore.Aggregate) {
	return func(ctx context.Context, aggregates []*eventstore.Aggregate) {
		ids := make([]string, len(aggregates))
		for i, aggregate := range aggregates {
			ids[i] = getID(aggregate)
		}
		err := cache.Invalidate(ctx, index, ids...)
		logging.OnError(err).Warn("cache invalidation failed")
	}
}

func getAggregateID(aggregate *eventstore.Aggregate) string {
	return aggregate.ID
}

func getResourceOwner(aggregate *eventstore.Aggregate) string {
	return aggregate.ResourceOwner
}
