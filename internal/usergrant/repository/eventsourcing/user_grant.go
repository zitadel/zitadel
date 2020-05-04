package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"
)

func UserGrantByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ols34", "id should be filled")
	}
	return UserGrantQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func UserGrantQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.UserGrantAggregate).
		LatestSequenceFilter(latestSequence)
}

func UserGrantAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, grant *model.UserGrant) (*es_models.Aggregate, error) {
	if grant == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dis83", "existing grant should not be nil")
	}
	return aggCreator.NewAggregate(ctx, grant.AggregateID, model.UserGrantAggregate, model.UserGrantVersion, grant.Sequence)
}

func UserGrantAddedAggregate(aggCreator *es_models.AggregateCreator, grant *model.UserGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := UserGrantAggregate(ctx, aggCreator, grant)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserGrantAdded, grant)
	}
}

func UserGrantChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.UserGrant, grant *model.UserGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-osl8x", "grant should not be nil")
		}
		agg, err := UserGrantAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := existing.Changes(grant)
		return agg.AppendEvent(model.UserGrantChanged, changes)
	}
}

func UserGrantDeactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *model.UserGrant, grant *model.UserGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-lo21s", "grant should not be nil")
		}
		agg, err := UserGrantAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserGrantDeactivated, nil)
	}
}

func UserGrantReactivatedAggregate(aggCreator *es_models.AggregateCreator, existing *model.UserGrant, grant *model.UserGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-mks34", "grant should not be nil")
		}
		agg, err := UserGrantAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserGrantReactivated, nil)
	}
}

func UserGrantRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.UserGrant, grant *model.UserGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-lo21s", "grant should not be nil")
		}
		agg, err := UserGrantAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.UserGrantRemoved, nil)
	}
}
