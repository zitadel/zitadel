package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/caos/logging"
	z_errors "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/lib/pq"
)

const (
	selectStmt = "SELECT" +
		" creation_date" +
		", event_type" +
		", event_sequence" +
		", previous_sequence" +
		", event_data" +
		", editor_service" +
		", editor_user" +
		", resource_owner" +
		", aggregate_type" +
		", aggregate_id" +
		", aggregate_version" +
		" FROM eventstore.events"
	selectWithSystemTimeStmt = "SELECT" +
		" creation_date" +
		", event_type" +
		", event_sequence" +
		", previous_sequence" +
		", event_data" +
		", editor_service" +
		", editor_user" +
		", resource_owner" +
		", aggregate_type" +
		", aggregate_id" +
		", aggregate_version" +
		" FROM eventstore.events" +
		" AS OF SYSTEM TIME '-50ms'"
)

func buildQuery(queryFactory *models.SearchQueryFactory) (query string, limit uint64, values []interface{}, rowScanner func(s scan, dest interface{}) error) {
	searchQuery, err := queryFactory.Build()
	if err != nil {
		logging.Log("SQL-cshKu").WithError(err).Warn("search query factory invalid")
		return "", 0, nil, nil
	}
	query, rowScanner = prepareColumns(searchQuery.Columns, searchQuery.IsPrecondition)
	where, values := prepareCondition(searchQuery.Filters)
	if where == "" || query == "" {
		return "", 0, nil, nil
	}
	query += where

	if searchQuery.Columns != models.Columns_Max_Sequence {
		query += " ORDER BY event_sequence"
		if searchQuery.Desc {
			query += " DESC"
		}
	}

	if searchQuery.Limit > 0 {
		values = append(values, searchQuery.Limit)
		query += " LIMIT ?"
	}

	query = numberPlaceholder(query, "?", "$")

	return query, searchQuery.Limit, values, rowScanner
}

func prepareCondition(filters []*models.Filter) (clause string, values []interface{}) {
	values = make([]interface{}, len(filters))
	clauses := make([]string, len(filters))

	if len(filters) == 0 {
		return clause, values
	}
	for i, filter := range filters {
		value := filter.GetValue()
		switch value.(type) {
		case []bool, []float64, []int64, []string, []models.AggregateType, []models.EventType, *[]bool, *[]float64, *[]int64, *[]string, *[]models.AggregateType, *[]models.EventType:
			value = pq.Array(value)
		}

		clauses[i] = getCondition(filter)
		if clauses[i] == "" {
			return "", nil
		}
		values[i] = value
	}
	return " WHERE " + strings.Join(clauses, " AND "), values
}

type scan func(dest ...interface{}) error

func prepareColumns(columns models.Columns, isPrecondition bool) (string, func(s scan, dest interface{}) error) {
	switch columns {
	case models.Columns_Max_Sequence:
		stmt := "SELECT MAX(event_sequence) FROM eventstore.events"
		if !isPrecondition {
			stmt = "SELECT MAX(event_sequence) FROM eventstore.events AS OF SYSTEM TIME '-50ms'"
		}
		return stmt, func(row scan, dest interface{}) (err error) {
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
	case models.Columns_Event:
		stmt := selectStmt
		if !isPrecondition {
			stmt = selectWithSystemTimeStmt
		}
		return stmt, func(row scan, dest interface{}) (err error) {
			event, ok := dest.(*models.Event)
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
				&event.AggregateVersion,
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
	default:
		return "", nil
	}
}

func numberPlaceholder(query, old, new string) string {
	for i, hasChanged := 1, true; hasChanged; i++ {
		newQuery := strings.Replace(query, old, new+strconv.Itoa(i), 1)
		hasChanged = query != newQuery
		query = newQuery
	}
	return query
}

func getCondition(filter *es_models.Filter) (condition string) {
	field := getField(filter.GetField())
	operation := getOperation(filter.GetOperation())
	if field == "" || operation == "" {
		return ""
	}
	format := getConditionFormat(filter.GetOperation())

	return fmt.Sprintf(format, field, operation)
}

func getConditionFormat(operation es_models.Operation) string {
	if operation == es_models.Operation_In {
		return "%s %s ANY(?)"
	}
	return "%s %s ?"
}

func getField(field es_models.Field) string {
	switch field {
	case es_models.Field_AggregateID:
		return "aggregate_id"
	case es_models.Field_AggregateType:
		return "aggregate_type"
	case es_models.Field_LatestSequence:
		return "event_sequence"
	case es_models.Field_ResourceOwner:
		return "resource_owner"
	case es_models.Field_EditorService:
		return "editor_service"
	case es_models.Field_EditorUser:
		return "editor_user"
	case es_models.Field_EventType:
		return "event_type"
	}
	return ""
}

func getOperation(operation es_models.Operation) string {
	switch operation {
	case es_models.Operation_Equals, es_models.Operation_In:
		return "="
	case es_models.Operation_Greater:
		return ">"
	case es_models.Operation_Less:
		return "<"
	}
	return ""
}
