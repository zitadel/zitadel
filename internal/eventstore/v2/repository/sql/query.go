package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/caos/logging"
	z_errors "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/lib/pq"
)

type criteriaer interface {
	columnName(repository.Field) string
	operation(repository.Operation) string
	conditionFormat(repository.Operation) string
	placeholder(query string) string
	eventQuery() string
	maxSequenceQuery() string
}

type rowScan func(scan, interface{}) error
type scan func(dest ...interface{}) error

func buildQuery(criteria criteriaer, searchQuery *repository.SearchQuery) (query string, values []interface{}, rowScanner rowScan) {
	query, rowScanner = prepareColumns(criteria, searchQuery.Columns)
	where, values := prepareCondition(criteria, searchQuery.Filters)
	if where == "" || query == "" {
		return "", nil, nil
	}
	query += where

	if searchQuery.Columns != repository.Columns_Max_Sequence {
		query += " ORDER BY event_sequence"
		if searchQuery.Desc {
			query += " DESC"
		}
	}

	if searchQuery.Limit > 0 {
		values = append(values, searchQuery.Limit)
		query += " LIMIT ?"
	}

	query = criteria.placeholder(query)

	return query, values, rowScanner
}

func prepareColumns(criteria criteriaer, columns repository.Columns) (string, func(s scan, dest interface{}) error) {
	switch columns {
	case repository.Columns_Max_Sequence:
		return criteria.maxSequenceQuery(), maxSequenceRowScanner
	case repository.Columns_Event:
		return criteria.eventQuery(), eventRowScanner
	default:
		return "", nil
	}
}

func maxSequenceRowScanner(row scan, dest interface{}) (err error) {
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

func eventRowScanner(row scan, dest interface{}) (err error) {
	event, ok := dest.(*repository.Event)
	if !ok {
		return z_errors.ThrowInvalidArgument(nil, "SQL-4GP6F", "type must be event")
	}
	var previousSequence Sequence
	data := make(Data, 0)

	err = row(
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
		logging.Log("SQL-kn1Sw").WithError(err).Warn("unable to scan row")
		return z_errors.ThrowInternal(err, "SQL-J0hFS", "unable to scan row")
	}

	event.PreviousSequence = uint64(previousSequence)

	event.Data = make([]byte, len(data))
	copy(event.Data, data)

	return nil
}

func prepareCondition(criteria criteriaer, filters []*repository.Filter) (clause string, values []interface{}) {
	values = make([]interface{}, len(filters))
	clauses := make([]string, len(filters))

	if len(filters) == 0 {
		return clause, values
	}
	for i, filter := range filters {
		value := filter.Value()
		switch value.(type) {
		case []bool, []float64, []int64, []string, []repository.AggregateType, []repository.EventType, *[]bool, *[]float64, *[]int64, *[]string, *[]repository.AggregateType, *[]repository.EventType:
			value = pq.Array(value)
		}

		clauses[i] = getCondition(criteria, filter)
		if clauses[i] == "" {
			return "", nil
		}
		values[i] = value
	}
	return " WHERE " + strings.Join(clauses, " AND "), values
}

func getCondition(cond criteriaer, filter *repository.Filter) (condition string) {
	field := cond.columnName(filter.Field())
	operation := cond.operation(filter.Operation())
	if field == "" || operation == "" {
		return ""
	}
	format := cond.conditionFormat(filter.Operation())

	return fmt.Sprintf(format, field, operation)
}
