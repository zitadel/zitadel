package execution

import (
	"context"
	"database/sql"
	_ "embed"

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
	InstanceByID(ctx context.Context, id string) (instance authz.Instance, err error)
	GetNotifyUserByID(ctx context.Context, shouldTriggered bool, userID string) (*query.NotifyUser, error)
	TargetsByExecutionID(ctx context.Context, ids []string) (execution []*query.ExecutionTarget, err error)

	ActiveInstances() []string
}

type ExecutionsQueries struct {
	Queries
	client *database.DB
}

func NewExecutionsQueries(
	baseQueries Queries,
	client *database.DB,
) *ExecutionsQueries {
	return &ExecutionsQueries{
		Queries: baseQueries,
		client:  client,
	}
}

type EventExecutions struct {
	query.SearchResponse
	EventExecutions []*EventExecution
}

func (w *ExecutionsQueries) searchEventExecutions(ctx context.Context) (eventExecutions *EventExecutions, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = w.client.QueryContext(ctx,
		func(rows *sql.Rows) error {
			eventExecutions, err = scanEventExecution(rows)
			return err
		},
		searchEventExecutions,
		authz.GetInstance(ctx).InstanceID(),
	)
	return eventExecutions, err
}

func scanEventExecution(rows *sql.Rows) (*EventExecutions, error) {
	eventExecutions := make([]*EventExecution, 0)
	var count uint64
	for rows.Next() {
		e := new(EventExecution)
		agg := new(eventstore.Aggregate)
		err := rows.Scan(
			&agg.InstanceID,
			&agg.ResourceOwner,
			&agg.Type,
			&agg.Version,
			&agg.ID,
			&e.Sequence,
			&e.EventType,
			&e.CreatedAt,
			&e.Position,
			&e.UserID,
			&e.EventData,
			&e.TargetsData,
			&count,
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
			Count: count,
		},
	}, nil
}
