package handlers

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func (n *NotificationQueries) IsAlreadyHandled(ctx context.Context, event eventstore.Event, data map[string]interface{}, eventTypes ...eventstore.EventType) (bool, error) {
	events, err := n.es.Filter(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			InstanceID(event.Aggregate().InstanceID).
			AddQuery().
			AggregateTypes(user.AggregateType).
			AggregateIDs(event.Aggregate().ID).
			SequenceGreater(event.Sequence()).
			EventTypes(eventTypes...).
			EventData(data).
			Builder(),
	)
	if err != nil {
		return false, err
	}
	return len(events) > 0, nil
}
