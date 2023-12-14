package query

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	failedEventsColumnProjectionName = "projection_name"
	failedEventsColumnFailedSequence = "failed_sequence"
	failedEventsColumnAggregateType  = "aggregate_type"
	failedEventsColumnAggregateID    = "aggregate_id"
	failedEventsColumnFailureCount   = "failure_count"
	failedEventsColumnLastFailed     = "last_failed"
	failedEventsColumnError          = "error"
	failedEventsColumnInstanceID     = "instance_id"
)

var (
	failedEventsTable = table{
		name:          projection.FailedEventsTable,
		instanceIDCol: failedEventsColumnInstanceID,
	}
	FailedEventsColumnProjectionName = Column{
		name:  failedEventsColumnProjectionName,
		table: failedEventsTable,
	}
	FailedEventsColumnFailedSequence = Column{
		name:  failedEventsColumnFailedSequence,
		table: failedEventsTable,
	}
	FailedeventsColumnAggregateType = Column{
		name:  failedEventsColumnAggregateType,
		table: failedEventsTable,
	}
	FailedeventsColumnAggregateID = Column{
		name:  failedEventsColumnAggregateID,
		table: failedEventsTable,
	}
	FailedEventsColumnFailureCount = Column{
		name:  failedEventsColumnFailureCount,
		table: failedEventsTable,
	}
	FailedEventsColumnLastFailed = Column{
		name:  failedEventsColumnLastFailed,
		table: failedEventsTable,
	}
	FailedEventsColumnError = Column{
		name:  failedEventsColumnError,
		table: failedEventsTable,
	}
	FailedEventsColumnInstanceID = Column{
		name:  failedEventsColumnInstanceID,
		table: failedEventsTable,
	}
)

type FailedEvents struct {
	SearchResponse
	FailedEvents []*FailedEvent
}

type FailedEvent struct {
	ProjectionName string
	AggregateType  string
	AggregateID    string
	FailedSequence uint64
	FailureCount   uint64
	Error          string
	LastFailed     time.Time
}

type FailedEventSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *Queries) SearchFailedEvents(ctx context.Context, queries *FailedEventSearchQueries) (failedEvents *FailedEvents, err error) {
	query, scan := prepareFailedEventsQuery(ctx, q.client)
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-n8rjJ", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		failedEvents, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-3j99J", "Errors.Internal")
	}
	return failedEvents, nil
}

func (q *Queries) RemoveFailedEvent(ctx context.Context, projectionName, instanceID string, sequence uint64) (err error) {
	stmt, args, err := sq.Delete(projection.FailedEventsTable).
		Where(sq.Eq{
			failedEventsColumnProjectionName: projectionName,
			failedEventsColumnFailedSequence: sequence,
			failedEventsColumnInstanceID:     instanceID,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "QUERY-DGgh3", "Errors.RemoveFailed")
	}
	_, err = q.client.ExecContext(ctx, stmt, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "QUERY-0kbFF", "Errors.RemoveFailed")
	}
	return nil
}

func NewFailedEventInstanceIDSearchQuery(instanceID string) (SearchQuery, error) {
	return NewTextQuery(FailedEventsColumnInstanceID, instanceID, TextEquals)
}

func (r *ProjectSearchQueries) AppendProjectionNameQuery(projectionName string) error {
	query, err := NewProjectResourceOwnerSearchQuery(projectionName)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (q *FailedEventSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func prepareFailedEventsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*FailedEvents, error)) {
	return sq.Select(
			FailedEventsColumnProjectionName.identifier(),
			FailedEventsColumnFailedSequence.identifier(),
			FailedeventsColumnAggregateType.identifier(),
			FailedeventsColumnAggregateID.identifier(),
			FailedEventsColumnFailureCount.identifier(),
			FailedEventsColumnLastFailed.identifier(),
			FailedEventsColumnError.identifier(),
			countColumn.identifier()).
			From(failedEventsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*FailedEvents, error) {
			failedEvents := make([]*FailedEvent, 0)
			var count uint64
			for rows.Next() {
				failedEvent := new(FailedEvent)
				var lastFailed sql.NullTime
				err := rows.Scan(
					&failedEvent.ProjectionName,
					&failedEvent.FailedSequence,
					&failedEvent.AggregateType,
					&failedEvent.AggregateID,
					&failedEvent.FailureCount,
					&lastFailed,
					&failedEvent.Error,
					&count,
				)
				if err != nil {
					return nil, err
				}
				failedEvent.LastFailed = lastFailed.Time
				failedEvents = append(failedEvents, failedEvent)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-En99f", "Errors.Query.CloseRows")
			}

			return &FailedEvents{
				FailedEvents: failedEvents,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
