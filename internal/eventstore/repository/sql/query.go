package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zitadel/logging"

	z_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

type querier interface {
	columnName(repository.Field) string
	operation(repository.Operation) string
	conditionFormat(repository.Operation) string
	placeholder(query string) string
	eventQuery() string
	maxCreationDateQuery() string
	instanceIDsQuery() string
	db() *sql.DB
	orderByCreationDate(desc bool) string
}

type scan func(dest ...interface{}) error

func query(ctx context.Context, criteria querier, searchQuery *repository.SearchQuery, dest interface{}) error {
	query, rowScanner := prepareColumns(criteria, searchQuery.Columns)
	values := make([]interface{}, 0, len(searchQuery.Filters)+2)
	where, conditionValues := prepareCondition(criteria, searchQuery)
	if where == "" || query == "" {
		return z_errors.ThrowInvalidArgument(nil, "SQL-rWeBw", "invalid query factory")
	}
	query += where
	values = append(values, conditionValues...)

	if searchQuery.Columns == repository.ColumnsEvent {
		query += criteria.orderByCreationDate(searchQuery.Desc)
	}

	if searchQuery.Limit > 0 {
		values = append(values, searchQuery.Limit)
		query += " LIMIT ?"
	}

	query = criteria.placeholder(query)

	var contextQuerier interface {
		QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	}
	contextQuerier = criteria.db()
	if searchQuery.Tx != nil {
		contextQuerier = searchQuery.Tx
	}

	rows, err := contextQuerier.QueryContext(ctx, query, values...)
	if err != nil {
		logging.New().WithError(err).Info("query failed")
		return z_errors.ThrowInternal(err, "SQL-KyeAx", "unable to filter events")
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
	case repository.ColumnsMaxCreationDate:
		return criteria.maxCreationDateQuery(), maxCreationDateScanner
	case repository.ColumnsInstanceIDs:
		return criteria.instanceIDsQuery(), instanceIDsScanner
	case repository.ColumnsEvent:
		return criteria.eventQuery(), eventsScanner
	default:
		return "", nil
	}
}

func maxCreationDateScanner(row scan, dest interface{}) (err error) {
	sequence, ok := dest.(*time.Time)
	if !ok {
		return z_errors.ThrowInvalidArgument(nil, "SQL-NBjA9", "type must be time.Time")
	}
	var creationDate sql.NullTime
	err = row(&creationDate)
	*sequence = creationDate.Time
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
	events, ok := dest.(*[]*repository.Event)
	if !ok {
		return z_errors.ThrowInvalidArgument(nil, "SQL-4GP6F", "type must be event")
	}
	data := make(Data, 0)
	event := new(repository.Event)

	err = scanner(
		&event.CreationDate,
		&event.Type,
		&data,
		&event.EditorService,
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

	event.Data = make([]byte, len(data))
	copy(event.Data, data)

	*events = append(*events, event)

	return nil
}

const sqlTimeLayout = "2006-01-02 15:04:05.999999-07:00"

func prepareCondition(criteria querier, searchQuery *repository.SearchQuery) (clause string, values []interface{}) {
	values = make([]interface{}, 0, len(searchQuery.Filters))

	if len(searchQuery.Filters) == 0 {
		return clause, values
	}

	clauses := make([]string, len(searchQuery.Filters))
	for idx, filter := range searchQuery.Filters {
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
	clause = " WHERE (" + strings.Join(clauses, " OR ") + ")"
	if !searchQuery.SystemTime.IsZero() {
		if searchQuery.Tx == nil {
			clause = " AS OF SYSTEM TIME '" + searchQuery.SystemTime.Format(sqlTimeLayout) + "'" + clause
		} else {
			clause += " AND " + criteria.columnName(repository.FieldCreationDate) + " = ?"
			values = append(values, searchQuery.SystemTime)
		}
	}
	return clause, values
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
