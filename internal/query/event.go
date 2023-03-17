package query

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
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

func (q *Queries) SearchEvents(ctx context.Context, query *eventstore.SearchQueryBuilder) (_ []*Event, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	events, err := q.eventstore.Filter(ctx, query.AllowTimeTravel())
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
	users := make(map[string]string)
	for i, event := range events {
		result[i] = q.convertEvent(ctx, event, users)
	}
	return result
}

func (q *Queries) convertEvent(ctx context.Context, event eventstore.Event, users map[string]string) *Event {
	ctx, span := tracing.NewSpan(ctx)
	var err error
	defer func() { span.EndWithError(err) }()

	displayName, ok := users[event.EditorUser()]
	if !ok {
		displayName = q.editorUserByID(ctx, event.EditorUser())
		users[event.EditorUser()] = displayName
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

func (q *Queries) editorUserByID(ctx context.Context, userID string) string {
	user, err := q.GetUserByID(ctx, false, userID, false)
	if err != nil {
		return userID
	}
	if user.Human != nil {
		return user.Human.DisplayName
	} else if user.Machine != nil {
		return user.Machine.Name
	}
	return userID
}
