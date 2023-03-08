package query

import (
	"context"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type Changes struct {
	Changes []*Change
}

type Change struct {
	ChangeDate            time.Time
	EventType             string
	Sequence              uint64
	ResourceOwner         string
	ModifierId            string
	ModifierName          string
	ModifierLoginName     string
	ModifierResourceOwner string
	ModifierAvatarKey     string
}

func (q *Queries) OrgChanges(ctx context.Context, orgID string, lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (_ *Changes, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query := func(query *eventstore.SearchQuery) {
		query.AggregateTypes(org.AggregateType).
			AggregateIDs(orgID)
	}
	return q.changes(ctx, query, lastSequence, limit, sortAscending, auditLogRetention)

}

func (q *Queries) ProjectChanges(ctx context.Context, projectID string, lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (_ *Changes, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query := func(query *eventstore.SearchQuery) {
		query.AggregateTypes(project.AggregateType).
			AggregateIDs(projectID)
	}
	return q.changes(ctx, query, lastSequence, limit, sortAscending, auditLogRetention)
}

func (q *Queries) ProjectGrantChanges(ctx context.Context, projectID, grantID string, lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (_ *Changes, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query := func(query *eventstore.SearchQuery) {
		query.AggregateTypes(project.AggregateType).
			AggregateIDs(projectID).
			EventData(map[string]interface{}{
				"grantId": grantID,
			})
	}
	return q.changes(ctx, query, lastSequence, limit, sortAscending, auditLogRetention)
}

func (q *Queries) ApplicationChanges(ctx context.Context, projectID, appID string, lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (_ *Changes, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query := func(query *eventstore.SearchQuery) {
		query.AggregateTypes(project.AggregateType).
			AggregateIDs(projectID).
			EventData(map[string]interface{}{
				"appId": appID,
			})
	}
	return q.changes(ctx, query, lastSequence, limit, sortAscending, auditLogRetention)
}

func (q *Queries) UserChanges(ctx context.Context, userID string, lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (_ *Changes, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query := func(query *eventstore.SearchQuery) {
		query.AggregateTypes(user.AggregateType).
			AggregateIDs(userID)
	}
	return q.changes(ctx, query, lastSequence, limit, sortAscending, auditLogRetention)
}

func (q *Queries) changes(ctx context.Context, query func(query *eventstore.SearchQuery), lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (*Changes, error) {
	builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).Limit(limit).AllowTimeTravel()
	if !sortAscending {
		builder.OrderDesc()
	}
	search := builder.AddQuery().SequenceGreater(lastSequence) //always use greater (less is done automatically by sorting desc)
	query(search)

	events, err := q.eventstore.Filter(ctx, builder)
	if err != nil {
		logging.Log("QUERY-ZRffs").WithError(err).Warn("eventstore unavailable")
		return nil, errors.ThrowInternal(err, "QUERY-328b1", "Errors.Internal")
	}
	if len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "QUERY-FpQqK", "Errors.Changes.NotFound")
	}
	changes := make([]*Change, 0, len(events))
	for _, event := range events {
		if auditLogRetention != 0 && event.CreationDate().Before(time.Now().Add(-auditLogRetention)) {
			continue
		}
		change := &Change{
			ChangeDate:        event.CreationDate(),
			EventType:         string(event.Type()),
			Sequence:          event.Sequence(),
			ResourceOwner:     event.Aggregate().ResourceOwner,
			ModifierId:        event.EditorUser(),
			ModifierName:      event.EditorUser(),
			ModifierLoginName: event.EditorUser(),
		}
		editor, _ := q.GetUserByID(ctx, false, change.ModifierId, false)
		if editor != nil {
			change.ModifierLoginName = editor.PreferredLoginName
			change.ModifierResourceOwner = editor.ResourceOwner
			if editor.Human != nil {
				change.ModifierName = editor.Human.DisplayName
				change.ModifierAvatarKey = editor.Human.AvatarKey
			}
			if editor.Machine != nil {
				change.ModifierName = editor.Machine.Name
			}
		}
		changes = append(changes, change)
	}
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "QUERY-DEGS2", "Errors.Changes.AuditRetention")
	}
	return &Changes{
		Changes: changes,
	}, nil
}
