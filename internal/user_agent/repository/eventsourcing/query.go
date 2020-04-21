package eventsourcing

import (
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user_agent/repository/eventsourcing/model"
)

func UserAgentByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-2sv2A", "id should be filled")
	}
	return UserAgentQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func UserAgentQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.UserAgentAggregate).
		LatestSequenceFilter(latestSequence)
}
