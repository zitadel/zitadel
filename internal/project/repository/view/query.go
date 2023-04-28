package view

import (
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func ProjectByIDQuery(id, instanceID string, latestSequence uint64) (*eventstore.SearchQueryBuilder, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dke74", "Errors.Project.ProjectIDMissing")
	}
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(instanceID).
		AddQuery().
		AggregateTypes(project.AggregateType).
		AggregateIDs(id).
		SequenceGreater(latestSequence).
		Builder(), nil
}
