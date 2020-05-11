package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func IamByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-0soe4", "id should be filled")
	}
	return IamQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func IamQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.IamAggregate).
		LatestSequenceFilter(latestSequence)
}

func IamAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, iam *model.Iam) (*es_models.Aggregate, error) {
	if iam == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-lo04e", "existing iam should not be nil")
	}
	return aggCreator.NewAggregate(ctx, iam.AggregateID, model.IamAggregate, model.IamVersion, iam.Sequence)
}

func IamAggregateOverwriteContext(ctx context.Context, aggCreator *es_models.AggregateCreator, iam *model.Iam, resourceOwnerID string, userID string) (*es_models.Aggregate, error) {
	if iam == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dis83", "existing iam should not be nil")
	}

	return aggCreator.NewAggregate(ctx, iam.AggregateID, model.IamAggregate, model.IamVersion, iam.Sequence, es_models.OverwriteResourceOwner(resourceOwnerID), es_models.OverwriteEditorUser(userID))
}

func IamSetupStartedAggregate(aggCreator *es_models.AggregateCreator, iam *model.Iam) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if iam == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-lo0sw", "iam should not be nil")
		}

		agg, err := IamAggregate(ctx, aggCreator, iam)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.IamSetupStarted, iam)
	}
}
