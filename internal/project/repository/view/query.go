package view

import (
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/project"
)

func ProjectByIDQuery(id, instanceID string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dke74", "Errors.Project.ProjectIDMissing")
	}
	return ProjectQuery(latestSequence).
		AddQuery().
		AggregateIDFilter(id).
		InstanceIDFilter(instanceID).
		SearchQuery(), nil
}

func ProjectQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AddQuery().
		AggregateTypeFilter(project.AggregateType).
		LatestSequenceFilter(latestSequence).
		SearchQuery()
}
