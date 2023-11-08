package handlers

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func (n *NotificationQueries) IsAlreadyHandled(ctx context.Context, event eventstore.Event, data map[string]interface{}, aggregateType eventstore.AggregateType, eventTypes ...eventstore.EventType) (bool, error) {
	events, err := n.es.Filter(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			InstanceID(event.Aggregate().InstanceID).
			SequenceGreater(event.Sequence()).
			AddQuery().
			AggregateTypes(aggregateType).
			AggregateIDs(event.Aggregate().ID).
			EventTypes(eventTypes...).
			EventData(data).
			Builder(),
	)
	if err != nil {
		return false, err
	}
	return len(events) > 0, nil
}
