package eventsourcing

import (
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/service_account/model"
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
