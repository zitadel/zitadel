package handlers

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
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
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(a.event.Aggregate().InstanceID).
		SequenceGreater(a.event.Sequence()).
		AddQuery().
		AggregateTypes(a.event.Aggregate().Type).
		AggregateIDs(a.event.Aggregate().ID).
		EventTypes(a.eventTypes...).
		EventData(a.data).
		Builder()
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
