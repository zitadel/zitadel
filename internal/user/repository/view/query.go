package view

import (
	"github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func UserByIDQuery(id, instanceID string, latestSequence uint64, eventTypes []es_models.EventType) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8isw", "Errors.User.UserIDMissing")
	}
	return es_models.NewSearchQuery().
		AddQuery().
		AggregateTypeFilter(user.AggregateType).
		AggregateIDFilter(id).
		EventTypesFilter(eventTypes...).
		LatestSequenceFilter(latestSequence).
		InstanceIDFilter(instanceID).
		SearchQuery(), nil
}
