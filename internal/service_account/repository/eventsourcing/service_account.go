package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/service_account/repository/eventsourcing/model"
)

func ServiceAccountByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-2gqWE", "Errors.ServiceAccount.ServiceAccountIDMissing")
	}
	return ServiceAccountQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func ServiceAccountQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.ServiceAccountAggregate).
		LatestSequenceFilter(latestSequence)
}

func ServiceAccountAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, serviceAccount *model.ServiceAccount) (*es_models.Aggregate, error) {
	return aggCreator.NewAggregate(ctx, serviceAccount.AggregateID, model.ServiceAccountAggregate, "v1", 0)
}

func ServiceAccountCreateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, serviceAccount *model.ServiceAccount) ([]*es_models.Aggregate, error) {
	aggregates := make([]*es_models.Aggregate, 0)
	aggregate, err := ServiceAccountAggregate(ctx, aggCreator, serviceAccount)
	if err != nil {
		return nil, err
	}
	accountAggregate, err := aggregate.AppendEvent(model.ServiceAccountAdded, serviceAccount)
	if err != nil {
		return nil, err
	}

	aggregates = append(aggregates, accountAggregate)
	return aggregates, nil
}
