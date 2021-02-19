package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func IAMByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-0soe4", "Errors.IAM.IDMissing")
	}
	return IAMQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func IAMQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.IAMAggregate).
		LatestSequenceFilter(latestSequence)
}

func IAMAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, iam *model.IAM) (*es_models.Aggregate, error) {
	if iam == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-lo04e", "Errors.Internal")
	}
	return aggCreator.NewAggregate(ctx, iam.AggregateID, model.IAMAggregate, model.IAMVersion, iam.Sequence)
}
