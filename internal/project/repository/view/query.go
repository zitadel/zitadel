package view

import (
	"time"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
)

func ProjectByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dke74", "Errors.Project.ProjectIDMissing")
	}
	return ProjectQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func ProjectQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.ProjectAggregate).
		LatestSequenceFilter(latestSequence)
}

func ChangesQuery(projectID string, latestSequence, limit uint64, sortAscending bool, retention time.Duration) *es_models.SearchQuery {
	query := es_models.NewSearchQuery().
		AggregateTypeFilter(model.ProjectAggregate)
	if !sortAscending {
		query.OrderDesc()
	}
	if retention > 0 {
		query.CreationDateNewerFilter(time.Now().Add(-retention))
	}

	query.LatestSequenceFilter(latestSequence).
		AggregateIDFilter(projectID).
		SetLimit(limit)
	return query
}
