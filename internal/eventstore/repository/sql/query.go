package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
	z_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

type querier interface {
	columnName(repository.Field) string
	operation(repository.Operation) string
	conditionFormat(repository.Operation) string
	placeholder(query string) string
	eventQuery() string
	maxSequenceQuery() string
	instanceIDsQuery() string
	db() *database.DB
	orderByEventSequence(desc bool) string
	dialect.Database
}

type scan func(dest ...interface{}) error

type tx struct {
	*sql.Tx
}

func (t *tx) QueryContext(ctx context.Context, scan func(rows *sql.Rows) error, query string, args ...any) error {
	rows, err := t.Tx.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := rows.Close()
		logging.OnError(closeErr).Info("rows.Close failed")
	}()

	if err = scan(rows); err != nil {
		return err
	}
	return rows.Err()
}

func query(ctx context.Context, criteria querier, searchQuery *eventstore.SearchQueryBuilder, dest interface{}) error {
	q, err := repository.QueryFromBuilder(searchQuery)
	if err != nil {
		return err
	}
	query, rowScanner := prepareColumns(criteria, q.Columns)
	where, values := prepareCondition(criteria, q.Filters)
	if where == "" || query == "" {
		return z_errors.ThrowInvalidArgument(nil, "SQL-rWeBw", "invalid query factory")
	}
	if q.Tx == nil {
		if travel := prepareTimeTravel(ctx, criteria, q.AllowTimeTravel); travel != "" {
			query += travel
		}
	}
	query += where

	if q.Columns == eventstore.ColumnsEvent {
		query += criteria.orderByEventSequence(q.Desc)
	}

	if q.Limit > 0 {
		values = append(values, q.Limit)
		query += " LIMIT ?"
	}

	query = criteria.placeholder(query)

	var contextQuerier interface {
		QueryContext(context.Context, func(rows *sql.Rows) error, string, ...interface{}) error
	}
	contextQuerier = criteria.db()
	if q.Tx != nil {
		contextQuerier = &tx{Tx: q.Tx}
	}

	err = contextQuerier.QueryContext(ctx,
		func(rows *sql.Rows) error {
			for rows.Next() {
				err := rowScanner(rows.Scan, dest)
				if err != nil {
					return err
				}
			}
			return nil
		}, query, values...)
	if err != nil {
		logging.New().WithError(err).Info("query failed")
		return z_errors.ThrowInternal(err, "SQL-KyeAx", "unable to filter events")
	}

	return nil
}

func prepareColumns(criteria querier, columns eventstore.Columns) (string, func(s scan, dest interface{}) error) {
	switch columns {
	case eventstore.ColumnsMaxSequence:
		return criteria.maxSequenceQuery(), maxSequenceScanner
	case eventstore.ColumnsInstanceIDs:
		return criteria.instanceIDsQuery(), instanceIDsScanner
	case eventstore.ColumnsEvent:
		return criteria.eventQuery(), eventsScanner
	default:
		return "", nil
	}
}

func prepareTimeTravel(ctx context.Context, criteria querier, allow bool) string {
	if !allow {
		return ""
	}
	took := call.Took(ctx)
	return criteria.Timetravel(took)
}

func maxSequenceScanner(row scan, dest interface{}) (err error) {
	sequence, ok := dest.(*sql.NullTime)
	if !ok {
		return z_errors.ThrowInvalidArgumentf(nil, "SQL-NBjA9", "type must be sql.NullTime got: %T", dest)
	}
	err = row(sequence)
	if err == nil || errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return z_errors.ThrowInternal(err, "SQL-bN5xg", "something went wrong")
}

func instanceIDsScanner(scanner scan, dest interface{}) (err error) {
	ids, ok := dest.(*[]string)
	if !ok {
		return z_errors.ThrowInvalidArgument(nil, "SQL-Begh2", "type must be an array of string")
	}
	var id string
	err = scanner(&id)
	if err != nil {
		logging.WithError(err).Warn("unable to scan row")
		return z_errors.ThrowInternal(err, "SQL-DEFGe", "unable to scan row")
	}
	*ids = append(*ids, id)

	return nil
}

func eventsScanner(scanner scan, dest interface{}) (err error) {
	events, ok := dest.(*[]eventstore.Event)
	if !ok {
		return z_errors.ThrowInvalidArgument(nil, "SQL-4GP6F", "type must be event")
	}
	var (
		previousAggregateSequence     Sequence
		previousAggregateTypeSequence Sequence
	)
	data := make(Data, 0)
	event := new(repository.Event)

	var editor sql.NullString

	err = scanner(
		&event.CreationDate,
		&event.Typ,
		&event.Seq,
		&previousAggregateSequence,
		&previousAggregateTypeSequence,
		&data,
		&editor,
		&event.EditorUser,
		&event.ResourceOwner,
		&event.InstanceID,
		&event.AggregateType,
		&event.AggregateID,
		&event.Version,
	)

	if err != nil {
		logging.New().WithError(err).Warn("unable to scan row")
		return z_errors.ThrowInternal(err, "SQL-M0dsf", "unable to scan row")
	}

	event.EditorService = editor.String
	event.PreviousAggregateSequence = uint64(previousAggregateSequence)
	event.PreviousAggregateTypeSequence = uint64(previousAggregateTypeSequence)
	event.Data = make([]byte, len(data))
	copy(event.Data, data)

	*events = append(*events, event)

	return nil
}

func prepareCondition(criteria querier, filters [][]*repository.Filter) (clause string, values []interface{}) {
	values = make([]interface{}, 0, len(filters))

	if len(filters) == 0 {
		return clause, values
	}

	clauses := make([]string, len(filters))
	for idx, filter := range filters {
		subClauses := make([]string, 0, len(filter))
		for _, f := range filter {
			value := f.Value
			switch value.(type) {
			case map[string]interface{}:
				var err error
				value, err = json.Marshal(value)
				if err != nil {
					logging.WithError(err).Warn("unable to marshal search value")
					continue
				}
			}

			subClauses = append(subClauses, getCondition(criteria, f))
			if subClauses[len(subClauses)-1] == "" {
				return "", nil
			}
			values = append(values, value)
		}
		clauses[idx] = "( " + strings.Join(subClauses, " AND ") + " )"
	}
	// created_at <= now() must be added because clock_timestamp() could be in the future
	// this could lead to skipped events which are not visible as of system time but have a lower
	// created_at timestamp
	return " WHERE (" + strings.Join(clauses, " OR ") + ") AND hlc_to_timestamp(crdb_internal_mvcc_timestamp)::TIMESTAMPTZ <= (SELECT COALESCE(MIN(start)::TIMESTAMPTZ, NOW()) FROM crdb_internal.cluster_transactions where application_name = 'zitadel')", values
}

func getCondition(cond querier, filter *repository.Filter) (condition string) {
	field := cond.columnName(filter.Field)
	operation := cond.operation(filter.Operation)
	if field == "" || operation == "" {
		return ""
	}
	format := cond.conditionFormat(filter.Operation)

	return fmt.Sprintf(format, field, operation)
}
