package view

import (
	"time"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

func UserByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8isw", "Errors.User.UserIDMissing")
	}
	return UserQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func UserQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.UserAggregate).
		LatestSequenceFilter(latestSequence)
}

func ChangesQuery(userID string, latestSequence, limit uint64, sortAscending bool, retention time.Duration) *es_models.SearchQuery {
	query := es_models.NewSearchQuery().
		AggregateTypeFilter(model.UserAggregate)
	if !sortAscending {
		query.OrderDesc()
	}
	if retention > 0 {
		query.CreationDateNewerFilter(time.Now().Add(-retention))
	}

	query.LatestSequenceFilter(latestSequence).
		AggregateIDFilter(userID).
		SetLimit(limit)
	return query
}
