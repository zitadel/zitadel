package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
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

func UserGrantUniqueQuery(resourceOwner, projectID, userID string) *es_models.SearchQuery {
	grantID := resourceOwner + projectID + userID
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.UserGrantUniqueAggregate).
		AggregateIDFilter(grantID).
		OrderDesc().
		SetLimit(1)
}

func UserGrantAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, grant *model.UserGrant) (*es_models.Aggregate, error) {
	if grant == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dis83", "existing grant should not be nil")
	}
	return aggCreator.NewAggregate(ctx, grant.AggregateID, model.UserGrantAggregate, model.UserGrantVersion, grant.Sequence)
}

func UserGrantAddedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, grant *model.UserGrant) ([]*es_models.Aggregate, error) {
	agg, err := UserGrantAggregate(ctx, aggCreator, grant)
	if err != nil {
		return nil, err
	}
	agg, err = agg.AppendEvent(model.UserGrantAdded, grant)
	if err != nil {
		return nil, err
	}
	uniqueAggregate, err := reservedUniqueUserGrantAggregate(ctx, aggCreator, grant)
	if err != nil {
		return nil, err
	}
	return []*es_models.Aggregate{
		agg,
		uniqueAggregate,
	}, nil
}

func reservedUniqueUserGrantAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, grant *model.UserGrant) (*es_models.Aggregate, error) {
	grantID := auth.GetCtxData(ctx).OrgID + grant.ProjectID + grant.UserID
	aggregate, err := aggCreator.NewAggregate(ctx, grantID, model.UserGrantUniqueAggregate, model.UserGrantVersion, 0)
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.UserGrantReserved, nil)
	if err != nil {
		return nil, err
	}

	return aggregate.SetPrecondition(UserGrantUniqueQuery(auth.GetCtxData(ctx).OrgID, grant.ProjectID, grant.UserID), isEventValidation(aggregate, model.UserGrantReserved)), nil
}

func releasedUniqueUserGrantAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, grant *model.UserGrant) (aggregate *es_models.Aggregate, err error) {
	grantID := grant.ResourceOwner + grant.ProjectID + grant.UserID
	aggregate, err = aggCreator.NewAggregate(ctx, grantID, model.UserGrantUniqueAggregate, model.UserGrantVersion, 0)
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(model.UserGrantReleased, nil)
	if err != nil {
		return nil, err
	}

	return aggregate.SetPrecondition(UserGrantUniqueQuery(grant.ResourceOwner, grant.ProjectID, grant.UserID), isEventValidation(aggregate, model.UserGrantReleased)), nil
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

func UserGrantRemovedAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *model.UserGrant, grant *model.UserGrant) ([]*es_models.Aggregate, error) {
	if grant == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-lo21s", "grant should not be nil")
	}
	agg, err := UserGrantAggregate(ctx, aggCreator, existing)
	if err != nil {
		return nil, err
	}
	agg, err = agg.AppendEvent(model.UserGrantRemoved, nil)
	if err != nil {
		return nil, err
	}
	uniqueAggregate, err := releasedUniqueUserGrantAggregate(ctx, aggCreator, existing)
	if err != nil {
		return nil, err
	}
	return []*es_models.Aggregate{
		agg,
		uniqueAggregate,
	}, nil
}

func isEventValidation(aggregate *es_models.Aggregate, eventType es_models.EventType) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		if len(events) == 0 {
			aggregate.PreviousSequence = 0
			return nil
		}
		if events[0].Type == eventType {
			return errors.ThrowPreconditionFailedf(nil, "EVENT-eJQqe", "user_grant is already %v", eventType)
		}
		aggregate.PreviousSequence = events[0].Sequence
		return nil
	}
}
