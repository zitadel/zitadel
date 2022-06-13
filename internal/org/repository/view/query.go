package view

import (
	"github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func OrgByIDQuery(id, instanceID string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dke74", "id should be filled")
	}
	return es_models.NewSearchQuery().
		AddQuery().
		AggregateTypeFilter(org.AggregateType).
		LatestSequenceFilter(latestSequence).
		InstanceIDFilter(instanceID).
		AggregateIDFilter(id).
		SearchQuery(), nil
}
