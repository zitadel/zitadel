package view

import (
	"time"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func UserByIDQuery(id, instanceID string, lastCreationDate time.Time, eventTypes []eventstore.EventType) (*eventstore.SearchQueryBuilder, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8isw", "Errors.User.UserIDMissing")
	}
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		InstanceID(instanceID).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(id).
		EventTypes(eventTypes...).
		CreationDateAfter(lastCreationDate.Add(-1 * time.Microsecond)).
		Builder(), nil
}
