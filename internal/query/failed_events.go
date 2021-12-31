package query

import (
	"context"
	"database/sql"
	errs "errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

const (
	failedEventsColumnProjectionName = "projection_name"
	failedEventsColumnFailedSequence = "failed_sequence"
	failedEventsColumnFailureCount   = "failure_count"
	failedEventsColumnError          = "error"
)

var (
	failedEventsTable = table{
		name: projection.FailedEventsTable,
	}
	FailedEventsColumnProjectionName = Column{
		name:  failedEventsColumnProjectionName,
		table: failedEventsTable,
	}
	FailedEventsColumnFailedSequence = Column{
		name:  failedEventsColumnFailedSequence,
		table: failedEventsTable,
	}
	FailedEventsColumnFailureCount = Column{
		name:  failedEventsColumnFailureCount,
		table: failedEventsTable,
	}
	FailedEventsColumnError = Column{
		name:  failedEventsColumnError,
		table: failedEventsTable,
	}
)

type FailedEvents struct {
	SearchResponse
	FailedEvents []*FailedEvent
}

type FailedEvent struct {
	ProjectionName string
	FailedSequence uint64
	FailureCount   uint64
	Error          string
}

type FailedEventSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *Queries) SearchFailedEvents(ctx context.Context, queries *FailedEventSearchQueries) (failedEvents *FailedEvents, err error) {
	query, scan := prepareFailedEventsQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-n8rjJ", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-3j99J", "Errors.Internal")
	}
	return scan(rows)
}

func (q *Queries) RemoveFailedEvent(ctx context.Context, projectionName string, sequence uint64) (err error) {
	stmt, args, err := sq.Delete(projection.FailedEventsTable).
		Where(sq.Eq{
			failedEventsColumnProjectionName: projectionName,
			failedEventsColumnFailedSequence: sequence,
		}).ToSql()
	if err != nil {
		return errors.ThrowInternal(err, "QUERY-DGgh3", "Errors.RemoveFailed")
	}
	_, err = q.client.Exec(stmt, args...)
	if err != nil {
		return errors.ThrowInternal(err, "QUERY-0kbFF", "Errors.RemoveFailed")
	}
	return nil
}

func NewFailedEventProjectionNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(FailedEventsColumnProjectionName, value, method)
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

func prepareFailedEventQuery() (sq.SelectBuilder, func(*sql.Row) (*FailedEvent, error)) {
	return sq.Select(
			FailedEventsColumnProjectionName.identifier(),
			FailedEventsColumnFailedSequence.identifier(),
			FailedEventsColumnFailureCount.identifier(),
			FailedEventsColumnError.identifier()).
			From(failedEventsTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*FailedEvent, error) {
			p := new(FailedEvent)
			err := row.Scan(
				&p.ProjectionName,
				&p.FailedSequence,
				&p.FailureCount,
				&p.Error,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-5N00f", "Errors.FailedEvents.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-0oJf3", "Errors.Internal")
			}
			return p, nil
		}
}

func prepareFailedEventsQuery() (sq.SelectBuilder, func(*sql.Rows) (*FailedEvents, error)) {
	return sq.Select(
			FailedEventsColumnProjectionName.identifier(),
			FailedEventsColumnFailedSequence.identifier(),
			FailedEventsColumnFailureCount.identifier(),
			FailedEventsColumnError.identifier(),
			countColumn.identifier()).
			From(failedEventsTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*FailedEvents, error) {
			failedEvents := make([]*FailedEvent, 0)
			var count uint64
			for rows.Next() {
				failedEvent := new(FailedEvent)
				err := rows.Scan(
					&failedEvent.ProjectionName,
					&failedEvent.FailedSequence,
					&failedEvent.FailureCount,
					&failedEvent.Error,
					&count,
				)
				if err != nil {
					return nil, err
				}
				failedEvents = append(failedEvents, failedEvent)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-En99f", "Errors.Query.CloseRows")
			}

			return &FailedEvents{
				FailedEvents: failedEvents,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
