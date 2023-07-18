package view

import (
	"github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func UserByIDQuery(id, instanceID string, latestSequence uint64, eventTypes ...es_models.EventType) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8isw", "Errors.User.UserIDMissing")
	}
	query := es_models.NewSearchQuery().
		AddQuery().
		AggregateTypeFilter(user.AggregateType).
		AggregateIDFilter(id).
		LatestSequenceFilter(latestSequence).
		InstanceIDFilter(instanceID).
		SearchQuery()
	if len(eventTypes) > 0 {
		query.AddQuery().EventTypesFilter(eventTypes...)
	}
	return query, nil
}
