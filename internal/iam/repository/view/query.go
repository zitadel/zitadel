package view

import (
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func IAMByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-4ng8sd", "id should be filled")
	}
	return IAMQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func IAMQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(iam_es_model.IAMAggregate).
		LatestSequenceFilter(latestSequence)
}
