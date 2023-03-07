package query

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/call"
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
	ID                string
	DisplayName       string
	Service           string
	PreferedLoginName string
	AvatarKey         string
}

func (q *Queries) SearchEvents(ctx context.Context, query *eventstore.SearchQueryBuilder, auditLogRetention time.Duration) (_ []*Event, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	events, err := q.eventstore.Filter(ctx, query.AllowTimeTravel())
	if err != nil {
		return nil, err
	}

	callTime := call.FromContext(ctx)
	if callTime.IsZero() {
		callTime = time.Now()
	}
	for i, event := range events {
		if event.CreationDate().Before(callTime.Add(-auditLogRetention)) {
			events = events[:i]
			break
		}
	}

	return q.convertEvents(ctx, events, auditLogRetention), nil
}

func (q *Queries) SearchEventTypes(ctx context.Context) []string {
	return q.eventstore.EventTypes()
}

func (q *Queries) SearchAggregateTypes(ctx context.Context) []string {
	return q.eventstore.AggregateTypes()
}

func (q *Queries) convertEvents(ctx context.Context, events []eventstore.Event, auditLogRetention time.Duration) []*Event {
	result := make([]*Event, len(events))
	users := make(map[string]*EventEditor)
	for i, event := range events {
		result[i] = q.convertEvent(ctx, event, users)
	}
	return result
}

func (q *Queries) convertEvent(ctx context.Context, event eventstore.Event, users map[string]*EventEditor) *Event {
	ctx, span := tracing.NewSpan(ctx)
	var err error
	defer func() { span.EndWithError(err) }()

	editor, ok := users[event.EditorUser()]
	if !ok {
		editor = q.editorUserByID(ctx, event.EditorUser())
		users[event.EditorUser()] = editor
	}

	return &Event{
		Editor: &EventEditor{
			ID:                event.EditorUser(),
			Service:           event.EditorService(),
			DisplayName:       editor.DisplayName,
			PreferedLoginName: editor.PreferedLoginName,
			AvatarKey:         editor.AvatarKey,
		},
		Aggregate:    event.Aggregate(),
		Sequence:     event.Sequence(),
		CreationDate: event.CreationDate(),
		Type:         string(event.Type()),
		Payload:      event.DataAsBytes(),
	}
}

func (q *Queries) editorUserByID(ctx context.Context, userID string) *EventEditor {
	user, err := q.GetUserByID(ctx, false, userID, false)
	if err != nil {
		return &EventEditor{ID: userID}
	}

	if user.Human != nil {
		return &EventEditor{
			ID:                user.ID,
			DisplayName:       user.Human.DisplayName,
			PreferedLoginName: user.PreferredLoginName,
			AvatarKey:         user.Human.AvatarKey,
		}
	} else if user.Machine != nil {
		return &EventEditor{
			ID:                user.ID,
			DisplayName:       user.Machine.Name,
			PreferedLoginName: user.PreferredLoginName,
		}
	}
	return &EventEditor{ID: userID}
}
