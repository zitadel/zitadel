package execution

import (
	"context"
	"database/sql"
	_ "embed"
	"slices"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

//go:embed search_event_executions.sql
var searchEventExecutions string

type Queries interface {
	ActiveInstancesWithFeatureFlag(ctx context.Context, featureFlag bool) []string
	InstanceByID(ctx context.Context, id string) (instance authz.Instance, err error)
}

type baseQueries interface {
	ActiveInstances() []string
	InstanceByID(ctx context.Context, id string) (instance authz.Instance, err error)
	GetNotifyUserByID(ctx context.Context, shouldTriggered bool, userID string) (*query.NotifyUser, error)
	TargetsByExecutionID(ctx context.Context, ids []string) (execution []*query.ExecutionTarget, err error)
	GetInstanceFeatures(ctx context.Context, cascade bool) (_ *query.InstanceFeatures, err error)
}

type ExecutionsQueries struct {
	Queries
	queries baseQueries
	client  *database.DB
}

func NewExecutionsQueries(
	baseQueries baseQueries,
	client *database.DB,
) *ExecutionsQueries {
	return &ExecutionsQueries{
		queries: baseQueries,
		client:  client,
	}
}

func (q *ExecutionsQueries) ActiveInstancesWithFeatureFlag(ctx context.Context) []string {
	return slices.DeleteFunc(q.queries.ActiveInstances(), func(s string) bool {
		features, err := q.queries.GetInstanceFeatures(ctx, true)
		if err != nil {
			return true
		}
		if features == nil || !features.Actions.Value {
			return true
		}
		return false
	})
}

type EventExecutions struct {
	query.SearchResponse
	EventExecutions []*EventExecution
}

func (w *ExecutionsQueries) searchEventExecutions(ctx context.Context, limit uint16) (eventExecutions *EventExecutions, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = w.client.QueryContext(ctx,
		func(rows *sql.Rows) error {
			eventExecutions, err = scanEventExecution(rows)
			return err
		},
		searchEventExecutions,
		authz.GetInstance(ctx).InstanceID(),
		limit,
	)
	return eventExecutions, err
}

func scanEventExecution(rows *sql.Rows) (*EventExecutions, error) {
	eventExecutions := make([]*EventExecution, 0)
	for rows.Next() {
		e := new(EventExecution)
		e.Aggregate = new(eventstore.Aggregate)
		err := rows.Scan(
			&e.Aggregate.InstanceID,
			&e.Aggregate.ResourceOwner,
			&e.Aggregate.Type,
			&e.Aggregate.Version,
			&e.Aggregate.ID,
			&e.Sequence,
			&e.EventType,
			&e.CreatedAt,
			&e.UserID,
			&e.EventData,
			&e.TargetsData,
		)
		if err != nil {
			return nil, err
		}

		eventExecutions = append(eventExecutions, e)
	}

	if err := rows.Close(); err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-TODO", "Errors.Query.CloseRows")
	}

	return &EventExecutions{
		EventExecutions: eventExecutions,
		SearchResponse: query.SearchResponse{
			Count: uint64(len(eventExecutions)),
		},
	}, nil
}
