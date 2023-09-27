package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
	columnName(field repository.Field, useV1 bool) string
	operation(repository.Operation) string
	conditionFormat(repository.Operation) string
	placeholder(query string) string
	eventQuery(useV1 bool) string
	maxSequenceQuery(useV1 bool) string
	instanceIDsQuery(useV1 bool) string
	db() *database.DB
	orderByEventSequence(desc, useV1 bool) string
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

func query(ctx context.Context, criteria querier, searchQuery *eventstore.SearchQueryBuilder, dest interface{}, useV1 bool) error {
	q, err := repository.QueryFromBuilder(searchQuery)
	if err != nil {
		return err
	}
	query, rowScanner := prepareColumns(criteria, q.Columns, useV1)
	where, values := prepareCondition(criteria, q, useV1)
	if where == "" || query == "" {
		return z_errors.ThrowInvalidArgument(nil, "SQL-rWeBw", "invalid query factory")
	}
	if q.Tx == nil {
		if travel := prepareTimeTravel(ctx, criteria, q.AllowTimeTravel); travel != "" {
			query += travel
		}
	}
	query += where

	// instead of using the max function of the database (which doesn't work for postgres)
	// we select the most recent row
	if q.Columns == eventstore.ColumnsMaxSequence {
		q.Limit = 1
		q.Desc = true
	}

	switch q.Columns {
	case eventstore.ColumnsEvent,
		eventstore.ColumnsMaxSequence:
		query += criteria.orderByEventSequence(q.Desc, useV1)
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

func prepareColumns(criteria querier, columns eventstore.Columns, useV1 bool) (string, func(s scan, dest interface{}) error) {
	switch columns {
	case eventstore.ColumnsMaxSequence:
		return criteria.maxSequenceQuery(useV1), maxSequenceScanner
	case eventstore.ColumnsInstanceIDs:
		return criteria.instanceIDsQuery(useV1), instanceIDsScanner
	case eventstore.ColumnsEvent:
		return criteria.eventQuery(useV1), eventsScanner(useV1)
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
	position, ok := dest.(*sql.NullFloat64)
	if !ok {
		return z_errors.ThrowInvalidArgumentf(nil, "SQL-NBjA9", "type must be sql.NullInt64 got: %T", dest)
	}
	err = row(position)
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

func eventsScanner(useV1 bool) func(scanner scan, dest interface{}) (err error) {
	return func(scanner scan, dest interface{}) (err error) {
		events, ok := dest.(*[]eventstore.Event)
		if !ok {
			return z_errors.ThrowInvalidArgument(nil, "SQL-4GP6F", "type must be event")
		}
		data := make(Data, 0)
		event := new(repository.Event)

		position := new(sql.NullFloat64)

		if useV1 {
			err = scanner(
				&event.CreationDate,
				&event.Typ,
				&event.Seq,
				&data,
				&event.EditorUser,
				&event.ResourceOwner,
				&event.InstanceID,
				&event.AggregateType,
				&event.AggregateID,
				&event.Version,
			)
		} else {
			var revision uint8
			err = scanner(
				&event.CreationDate,
				&event.Typ,
				&event.Seq,
				position,
				&data,
				&event.EditorUser,
				&event.ResourceOwner,
				&event.InstanceID,
				&event.AggregateType,
				&event.AggregateID,
				&revision,
			)
			event.Version = eventstore.Version("v" + strconv.Itoa(int(revision)))
		}

		if err != nil {
			logging.New().WithError(err).Warn("unable to scan row")
			return z_errors.ThrowInternal(err, "SQL-M0dsf", "unable to scan row")
		}

		event.Data = make([]byte, len(data))
		copy(event.Data, data)
		event.Pos = position.Float64

		*events = append(*events, event)

		return nil
	}
}

func prepareCondition(criteria querier, query *repository.SearchQuery, useV1 bool) (clause string, values []interface{}) {
	values = make([]interface{}, 0, len(query.Filters))

	if len(query.Filters) == 0 {
		return clause, values
	}

	clauses := make([]string, len(query.Filters))
	for idx, filter := range query.Filters {
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

			subClauses = append(subClauses, getCondition(criteria, f, useV1))
			if subClauses[len(subClauses)-1] == "" {
				return "", nil
			}
			values = append(values, value)
		}
		clauses[idx] = "( " + strings.Join(subClauses, " AND ") + " )"
	}

	where := " WHERE (" + strings.Join(clauses, " OR ") + ") "
	if query.AwaitOpenTransactions {
		where += awaitOpenTransactions(useV1)
	}

	return where, values
}

func getCondition(cond querier, filter *repository.Filter, useV1 bool) (condition string) {
	field := cond.columnName(filter.Field, useV1)
	operation := cond.operation(filter.Operation)
	if field == "" || operation == "" {
		return ""
	}
	format := cond.conditionFormat(filter.Operation)

	return fmt.Sprintf(format, field, operation)
}
