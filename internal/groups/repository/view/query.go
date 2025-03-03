package view

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func GroupByIDQuery(id, instanceID string, latestSequence uint64) (*eventstore.SearchQueryBuilder, error) {
	if id == "" {
		return nil, zerrors.ThrowPreconditionFailed(nil, "EVENT-dkf84", "Errors.Group.GroupIDMissing")
	}
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(instanceID).
		AwaitOpenTransactions().
		SequenceGreater(latestSequence).
		AddQuery().
		AggregateTypes(group.AggregateType).
		AggregateIDs(id).
		EventTypes(
			group.GroupAddedType,
			group.GroupChangedType,
			group.GroupDeactivatedType,
			group.GroupReactivatedType,
			group.GroupRemovedType,
		).
		Builder(), nil
}
