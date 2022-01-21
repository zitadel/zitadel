package query

import (
	"context"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/repository/user"

	"github.com/caos/zitadel/internal/repository/project"

	"github.com/caos/zitadel/internal/repository/org"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
)

type Changes struct {
	LastSequence uint64
	Changes      []*Change
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

func (q *Queries) OrgChanges(ctx context.Context, orgID string, lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (*Changes, error) {
	query := func(query *eventstore.SearchQuery) {
		query.AggregateTypes(org.AggregateType).
			AggregateIDs(orgID)
	}
	return q.changes(ctx, query, lastSequence, limit, sortAscending, auditLogRetention)

}

func (q *Queries) ProjectChanges(ctx context.Context, projectID string, lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (*Changes, error) {
	query := func(query *eventstore.SearchQuery) {
		query.AggregateTypes(project.AggregateType).
			AggregateIDs(projectID)
	}
	return q.changes(ctx, query, lastSequence, limit, sortAscending, auditLogRetention)
}

func (q *Queries) ApplicationChanges(ctx context.Context, projectID, appID string, lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (*Changes, error) {
	query := func(query *eventstore.SearchQuery) {
		query.AggregateTypes(project.AggregateType).
			AggregateIDs(projectID).
			EventData(map[string]interface{}{
				"appId": appID,
			})
	}
	return q.changes(ctx, query, lastSequence, limit, sortAscending, auditLogRetention)
}

func (q *Queries) UserChanges(ctx context.Context, userID string, lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (*Changes, error) {
	query := func(query *eventstore.SearchQuery) {
		query.AggregateTypes(user.AggregateType).
			AggregateIDs(userID)
	}
	return q.changes(ctx, query, lastSequence, limit, sortAscending, auditLogRetention)
}

func (q *Queries) changes(ctx context.Context, query func(query *eventstore.SearchQuery), lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (*Changes, error) {
	builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).Limit(limit)
	search := builder.AddQuery()
	query(search)
	if sortAscending {
		builder.OrderAsc()
		search.SequenceGreater(lastSequence)
	} else {
		builder.OrderAsc()
		search.SequenceLess(lastSequence)
	}

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
		if event.CreationDate().After(time.Now().Add(-auditLogRetention)) {
			continue
		}
		lastSequence = event.Sequence()
		change := &Change{
			ChangeDate:        event.CreationDate(),
			EventType:         string(event.Type()),
			Sequence:          event.Sequence(),
			ResourceOwner:     event.Aggregate().ResourceOwner,
			ModifierId:        event.EditorUser(),
			ModifierName:      event.EditorUser(),
			ModifierLoginName: event.EditorUser(),
		}
		editor, _ := q.GetUserByID(ctx, change.ModifierId)
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
		LastSequence: lastSequence,
		Changes:      changes,
	}, nil
}
