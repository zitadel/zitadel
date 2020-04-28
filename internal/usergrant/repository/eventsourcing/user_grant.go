package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
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
		AggregateTypeFilter(grant_model.UserGrantAggregate).
		LatestSequenceFilter(latestSequence)
}

func UserGrantAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, grant *model.UserGrant) (*es_models.Aggregate, error) {
	if grant == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dis83", "existing user should not be nil")
	}
	return aggCreator.NewAggregate(ctx, grant.AggregateID, grant_model.UserGrantAggregate, model.UserGrantVersion, grant.Sequence)
}

func UserGrantAddedAggregate(aggCreator *es_models.AggregateCreator, grant *model.UserGrant) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if grant == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dlox9", "grant should not be nil")
		}
		agg, err := UserGrantAggregate(ctx, aggCreator, grant)
		if err != nil {
			return nil, err
		}
		agg.AppendEvent(grant_model.UserGrantAdded, grant)
		return agg, nil
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
		agg.AppendEvent(grant_model.UserGrantChanged, changes)

		return agg, nil
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
		agg.AppendEvent(grant_model.UserGrantDeactivated, nil)

		return agg, nil
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
		agg.AppendEvent(grant_model.UserGrantReactivated, nil)

		return agg, nil
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
		agg.AppendEvent(grant_model.UserGrantRemoved, nil)

		return agg, nil
	}
}
