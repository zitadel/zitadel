package handlers

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type alreadyHandled struct {
	event      eventstore.Event
	eventTypes []eventstore.EventType
	data       map[string]interface{}

	handled bool
}

func (a *alreadyHandled) Reduce() error {
	return nil
}

func (a *alreadyHandled) AppendEvents(event ...eventstore.Event) {
	if len(event) > 0 {
		a.handled = true
	}
}

func (a *alreadyHandled) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(a.event.Aggregate().InstanceID).
		SequenceGreater(a.event.Sequence()).
		AddQuery().
		AggregateTypes(a.event.Aggregate().Type).
		AggregateIDs(a.event.Aggregate().ID).
		EventTypes(a.eventTypes...).
		EventData(a.data)
	rejectedType, ok := notificationRejectedTypes[a.event.Aggregate().Type]
	if !ok {
		return query.Builder()
	}
	// a rejected notification is handled as well, it must not be sent on a later attempt
	return query.Or().
		AggregateTypes(a.event.Aggregate().Type).
		AggregateIDs(a.event.Aggregate().ID).
		EventTypes(rejectedType).
		EventData(channels.RejectionsOf(a.event.Type())).
		Builder()
}

// notificationRejectedTypes are the events stating a rejected notification by the type of the aggregate that triggered it.
var notificationRejectedTypes = map[eventstore.AggregateType]eventstore.EventType{
	user.AggregateType:    user.NotificationRejectedType,
	session.AggregateType: session.NotificationRejectedType,
}

func (n *NotificationQueries) IsAlreadyHandled(ctx context.Context, event eventstore.Event, data map[string]interface{}, eventTypes ...eventstore.EventType) (bool, error) {
	already := &alreadyHandled{
		event:      event,
		eventTypes: eventTypes,
		data:       data,
	}
	err := n.es.FilterToQueryReducer(ctx, already)
	if err != nil {
		return false, err
	}
	return already.handled, nil
}
