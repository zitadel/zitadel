package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/caos/logging"
	z_errors "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/lib/pq"
)

type querier interface {
	columnName(repository.Field) string
	operation(repository.Operation) string
	conditionFormat(repository.Operation) string
	placeholder(query string) string
	eventQuery() string
	maxSequenceQuery() string
	db() *sql.DB
	orderByEventSequence(desc bool) string
}

type scan func(dest ...interface{}) error

func query(ctx context.Context, criteria querier, searchQuery *repository.SearchQuery, dest interface{}) error {
	query, rowScanner := prepareColumns(criteria, searchQuery.Columns)
	where, values := prepareCondition(criteria, searchQuery.Filters)
	if where == "" || query == "" {
		return z_errors.ThrowInvalidArgument(nil, "SQL-rWeBw", "invalid query factory")
	}
	query += where

	if searchQuery.Columns != repository.ColumnsMaxSequence {
		query += criteria.orderByEventSequence(searchQuery.Desc)
	}

	if searchQuery.Limit > 0 {
		values = append(values, searchQuery.Limit)
		query += " LIMIT ?"
	}

	query = criteria.placeholder(query)

	rows, err := criteria.db().QueryContext(ctx, query, values...)
	if err != nil {
		logging.Log("SQL-HP3Uk").WithError(err).Info("query failed")
		return z_errors.ThrowInternal(err, "SQL-IJuyR", "unable to filter events")
	}
	defer rows.Close()

	for rows.Next() {
		err = rowScanner(rows.Scan, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func prepareColumns(criteria querier, columns repository.Columns) (string, func(s scan, dest interface{}) error) {
	switch columns {
	case repository.ColumnsMaxSequence:
		return criteria.maxSequenceQuery(), maxSequenceScanner
	case repository.ColumnsEvent:
		return criteria.eventQuery(), eventsScanner
	default:
		return "", nil
	}
}

func maxSequenceScanner(row scan, dest interface{}) (err error) {
	sequence, ok := dest.(*Sequence)
	if !ok {
		return z_errors.ThrowInvalidArgument(nil, "SQL-NBjA9", "type must be sequence")
	}
	err = row(sequence)
	if err == nil || errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return z_errors.ThrowInternal(err, "SQL-bN5xg", "something went wrong")
}

func eventsScanner(scanner scan, dest interface{}) (err error) {
	events, ok := dest.(*[]*repository.Event)
	if !ok {
		return z_errors.ThrowInvalidArgument(nil, "SQL-4GP6F", "type must be event")
	}
	var previousSequence Sequence
	data := make(Data, 0)
	event := new(repository.Event)

	err = scanner(
		&event.CreationDate,
		&event.Type,
		&event.Sequence,
		&previousSequence,
		&data,
		&event.EditorService,
		&event.EditorUser,
		&event.ResourceOwner,
		&event.AggregateType,
		&event.AggregateID,
		&event.Version,
	)

	if err != nil {
		logging.Log("SQL-3mofs").WithError(err).Warn("unable to scan row")
		return z_errors.ThrowInternal(err, "SQL-M0dsf", "unable to scan row")
	}

	event.PreviousSequence = uint64(previousSequence)
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
		subClauses := make([]string, len(filter))
		for subIdx, f := range filter {
			value := f.Value
			switch value.(type) {
			case []bool, []float64, []int64, []string, []repository.AggregateType, []repository.EventType, *[]bool, *[]float64, *[]int64, *[]string, *[]repository.AggregateType, *[]repository.EventType:
				value = pq.Array(value)
			case map[string]interface{}:
				var err error
				value, err = json.Marshal(value)
				logging.Log("SQL-BSsNy").OnError(err).Warn("unable to marshal search value")
			}

			subClauses[subIdx] = getCondition(criteria, f)
			if subClauses[subIdx] == "" {
				return "", nil
			}
			values = append(values, value)
		}
		clauses[idx] = "( " + strings.Join(subClauses, " AND ") + " )"
	}
	return " WHERE " + strings.Join(clauses, " OR "), values
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
