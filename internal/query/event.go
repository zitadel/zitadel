package query

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type Event struct {
	Editor       *EventEditor
	Aggregate    *eventstore.Aggregate
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

type eventsReducer struct {
	ctx     context.Context
	q       *Queries
	events  []*Event
	editors map[string]*EventEditor
}

func (r *eventsReducer) AppendEvents(events ...eventstore.Event) {
	r.events = append(r.events, r.convertEvents(r.ctx, events)...)
}

func (r *eventsReducer) Reduce() error { return nil }

func (q *Queries) SearchEvents(ctx context.Context, query *eventstore.SearchQueryBuilder) (_ []*Event, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	auditLogRetention := q.defaultAuditLogRetention
	if instanceAuditLogRetention := authz.GetInstance(ctx).AuditLogRetention(); instanceAuditLogRetention != nil {
		auditLogRetention = *instanceAuditLogRetention
	}
	if auditLogRetention != 0 {
		query = filterAuditLogRetention(ctx, auditLogRetention, query)
	}
	reducer := &eventsReducer{ctx: ctx, q: q, editors: make(map[string]*EventEditor, query.GetLimit())}
	if err = q.eventstore.FilterToReducer(ctx, query, reducer); err != nil {
		return nil, err
	}
	return reducer.events, nil
}

func filterAuditLogRetention(ctx context.Context, auditLogRetention time.Duration, builder *eventstore.SearchQueryBuilder) *eventstore.SearchQueryBuilder {
	callTime := call.FromContext(ctx)
	if callTime.IsZero() {
		callTime = time.Now()
	}
	oldestAllowed := callTime.Add(-auditLogRetention)
	// The audit log retention time should overwrite the creation date after query only if it is older
	// For example API calls should still be able to restrict the creation date after to a more recent date
	if builder.GetCreationDateAfter().Before(oldestAllowed) {
		return builder.CreationDateAfter(oldestAllowed)
	}
	return builder
}

func (q *Queries) SearchEventTypes(ctx context.Context) []string {
	return q.eventstore.EventTypes()
}

func (q *Queries) SearchAggregateTypes(ctx context.Context) []string {
	return q.eventstore.AggregateTypes()
}

func (er *eventsReducer) convertEvents(ctx context.Context, events []eventstore.Event) []*Event {
	result := make([]*Event, len(events))
	for i, event := range events {
		result[i] = er.convertEvent(ctx, event)
	}
	return result
}

func (er *eventsReducer) convertEvent(ctx context.Context, event eventstore.Event) *Event {
	editor, ok := er.editors[event.Creator()]
	if !ok {
		editor = er.q.editorUserByID(ctx, event.Creator())
		er.editors[event.Creator()] = editor
	}

	return &Event{
		Editor: &EventEditor{
			ID:                event.Creator(),
			Service:           "zitadel",
			DisplayName:       editor.DisplayName,
			PreferedLoginName: editor.PreferedLoginName,
			AvatarKey:         editor.AvatarKey,
		},
		Aggregate:    event.Aggregate(),
		Sequence:     event.Sequence(),
		CreationDate: event.CreatedAt(),
		Type:         string(event.Type()),
		Payload:      event.DataAsBytes(),
	}
}

func (q *Queries) editorUserByID(ctx context.Context, userID string) *EventEditor {
	user, err := q.GetUserByID(ctx, false, userID)
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
