package query

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type Event struct {
	Editor       *EventEditor
	Aggregate    eventstore.Aggregate
	Sequence     uint64
	CreationDate time.Time
	Type         string
	Payload      []byte
}

type EventEditor struct {
	ID          string
	DisplayName string
	Service     string
}

func (q *Queries) SearchEvents(ctx context.Context, query *eventstore.SearchQueryBuilder) ([]*Event, error) {
	events, err := q.eventstore.Filter(ctx, query)
	if err != nil {
		return nil, err
	}

	return q.convertEvents(ctx, events), nil
}

func (q *Queries) SearchEventTypes(ctx context.Context) []string {
	return q.eventstore.EventTypes()
}

func (q *Queries) SearchAggregateTypes(ctx context.Context) []string {
	return q.eventstore.AggregateTypes()
}

func (q *Queries) convertEvents(ctx context.Context, events []eventstore.Event) []*Event {
	result := make([]*Event, len(events))
	for i, event := range events {
		result[i] = q.convertEvent(ctx, event)
	}
	return result
}

func (q *Queries) convertEvent(ctx context.Context, event eventstore.Event) *Event {
	displayName := event.EditorUser()
	user, err := q.GetUserByID(ctx, false, event.EditorUser(), false)
	if err == nil {
		if user.Human != nil {
			displayName = user.Human.DisplayName
		} else if user.Machine != nil {
			displayName = user.Machine.Name
		}
	}

	return &Event{
		Editor: &EventEditor{
			ID:          event.EditorUser(),
			Service:     event.EditorService(),
			DisplayName: displayName,
		},
		Aggregate:    event.Aggregate(),
		Sequence:     event.Sequence(),
		CreationDate: event.CreationDate(),
		Type:         string(event.Type()),
		Payload:      event.DataAsBytes(),
	}
}