package view

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

func ChangesQuery(orgID string, latestSequence, limit uint64, sortAscending bool) *es_models.SearchQuery {
	query := es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgAggregate)

	if !sortAscending {
		query.OrderDesc()
	}

	query.LatestSequenceFilter(latestSequence).
		AggregateIDFilter(orgID).
		SetLimit(limit)
	return query
}
