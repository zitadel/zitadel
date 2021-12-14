package query

import (
	"context"
	"database/sql"
	errs "errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

type LatestSequence struct {
	Sequence  uint64
	Timestamp time.Time
}

type CurrentSequences struct {
	SearchResponse
	CurrentSequences []*CurrentSequence
}

type CurrentSequence struct {
	ProjectionName  string
	CurrentSequence uint64
	Timestamp       time.Time
}

type CurrentSequencesSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *CurrentSequencesSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) SearchCurrentSequences(ctx context.Context, queries *CurrentSequencesSearchQueries) (failedEvents *CurrentSequences, err error) {
	query, scan := prepareCurrentSequencesQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-MmFef", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-22H8f", "Errors.Internal")
	}
	return scan(rows)
}

func (q *Queries) latestSequence(ctx context.Context, projection table) (*LatestSequence, error) {
	query, scan := prepareLatestSequence()
	stmt, args, err := query.Where(sq.Eq{
		CurrentSequenceColProjectionName.identifier(): projection.name,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-5CfX9", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func (q *Queries) ClearCurrentSequence(ctx context.Context, projectionName string) (err error) {
	tx, err := q.client.Begin()
	if err != nil {
		return errors.ThrowInternal(err, "QUERY-9iOpr", "Errors.RemoveFailed")
	}
	projectionQuery, args, err := sq.Select("count(*)").
		From("[show tables from zitadel.projections]").
		Where(
			sq.And{
				sq.NotEq{"table_name": []string{"locks", "current_sequences", "failed_events"}},
				sq.Eq{"concat('zitadel.projections.', table_name)": projectionName},
			}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	row := tx.QueryRowContext(ctx, projectionQuery, args...)
	var count int
	if err := row.Scan(&count); err != nil || count == 0 {
		return errors.ThrowInternal(err, "QUERY-ej8fn", "Errors.ProjectionName.Invalid")
	}
	tablesQuery, args, err := sq.Select("concat('zitadel.projections.', table_name)").
		From("[show tables from zitadel.projections]").
		Where(
			sq.And{
				sq.Eq{"type": "table"},
				sq.NotEq{"table_name": []string{"locks", "current_sequences", "failed_events"}},
				sq.Like{"concat('zitadel.projections.', table_name)": projectionName + "%"},
			}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var tables []string
	rows, err := tx.QueryContext(ctx, tablesQuery, args...)
	if err != nil {
		return errors.ThrowInternal(err, "QUERY-Dgfw", "Errors.ProjectionName.Invalid")
	}
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return errors.ThrowInternal(err, "QUERY-ej8fn", "Errors.ProjectionName.Invalid")
		}
		tables = append(tables, tableName)
	}
	for _, tableName := range tables {
		_, err = tx.Exec(fmt.Sprintf("TRUNCATE %s cascade", tableName))
		if err != nil {
			return errors.ThrowInternal(err, "QUERY-3n92f", "Errors.RemoveFailed")
		}
	}
	update, args, err := sq.Update(currentSequencesTable.identifier()).
		Set(CurrentSequenceColCurrentSequence.name, 0).
		Where(sq.Eq{CurrentSequenceColProjectionName.name: projectionName}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return errors.ThrowInternal(err, "QUERY-Ff3tw", "Errors.RemoveFailed")
	}
	_, err = tx.Exec(update, args...)
	if err != nil {
		return errors.ThrowInternal(err, "QUERY-NFiws", "Errors.RemoveFailed")
	}
	return tx.Commit()
}

func prepareLatestSequence() (sq.SelectBuilder, func(*sql.Row) (*LatestSequence, error)) {
	return sq.Select(
			CurrentSequenceColCurrentSequence.identifier(),
			CurrentSequenceColTimestamp.identifier()).
			From(currentSequencesTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*LatestSequence, error) {
			seq := new(LatestSequence)
			err := row.Scan(
				&seq.Sequence,
				&seq.Timestamp,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-gmd9o", "Errors.CurrentSequence.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-aAZ1D", "Errors.Internal")
			}
			return seq, nil
		}
}

func prepareCurrentSequencesQuery() (sq.SelectBuilder, func(*sql.Rows) (*CurrentSequences, error)) {
	return sq.Select(
			"max("+CurrentSequenceColCurrentSequence.identifier()+") as "+CurrentSequenceColCurrentSequence.name,
			"max("+CurrentSequenceColTimestamp.identifier()+") as "+CurrentSequenceColTimestamp.name,
			CurrentSequenceColProjectionName.identifier(),
			countColumn.identifier()).
			From(currentSequencesTable.identifier()).
			GroupBy(CurrentSequenceColProjectionName.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*CurrentSequences, error) {
			currentSequences := make([]*CurrentSequence, 0)
			var count uint64
			for rows.Next() {
				currentSequence := new(CurrentSequence)
				err := rows.Scan(
					&currentSequence.CurrentSequence,
					&currentSequence.Timestamp,
					&currentSequence.ProjectionName,
					&count,
				)
				if err != nil {
					return nil, err
				}
				currentSequences = append(currentSequences, currentSequence)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-jbJ77", "Errors.Query.CloseRows")
			}

			return &CurrentSequences{
				CurrentSequences: currentSequences,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

var (
	currentSequencesTable = table{
		name: projection.CurrentSeqTable,
	}
	CurrentSequenceColAggregateType = Column{
		name:  "aggregate_type",
		table: currentSequencesTable,
	}
	CurrentSequenceColCurrentSequence = Column{
		name:  "current_sequence",
		table: currentSequencesTable,
	}
	CurrentSequenceColTimestamp = Column{
		name:  "timestamp",
		table: currentSequencesTable,
	}
	CurrentSequenceColProjectionName = Column{
		name:  "projection_name",
		table: currentSequencesTable,
	}
)
