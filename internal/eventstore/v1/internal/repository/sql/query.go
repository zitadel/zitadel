package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database/dialect"
	z_errors "github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

const (
	selectStmt = "SELECT" +
		" creation_date" +
		", event_type" +
		", event_sequence" +
		", previous_aggregate_sequence" +
		", event_data" +
		", editor_service" +
		", editor_user" +
		", resource_owner" +
		", instance_id" +
		", aggregate_type" +
		", aggregate_id" +
		", aggregate_version" +
		" FROM eventstore.events"
)

func buildQuery(ctx context.Context, db dialect.Database, queryFactory *es_models.SearchQueryFactory) (query string, limit uint64, values []interface{}, rowScanner func(s scan, dest interface{}) error) {
	searchQuery, err := queryFactory.Build()
	if err != nil {
		logging.New().WithError(err).Warn("search query factory invalid")
		return "", 0, nil, nil
	}
	query, rowScanner = prepareColumns(searchQuery.Columns)
	where, values := prepareCondition(searchQuery.Filters)
	if where == "" || query == "" {
		return "", 0, nil, nil
	}

	if travel := db.Timetravel(call.Took(ctx)); travel != "" {
		query += travel
	}
	query += where

	if searchQuery.Columns == es_models.Columns_Event {
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

func prepareCondition(filters [][]*es_models.Filter) (clause string, values []interface{}) {
	values = make([]interface{}, 0, len(filters))
	clauses := make([]string, len(filters))

	if len(filters) == 0 {
		return clause, values
	}
	for i, filter := range filters {
		subClauses := make([]string, 0, len(filter))
		for _, f := range filter {
			value := f.GetValue()

			subClauses = append(subClauses, getCondition(f))
			if subClauses[len(subClauses)-1] == "" {
				return "", nil
			}
			values = append(values, value)
		}
		clauses[i] = "( " + strings.Join(subClauses, " AND ") + " )"
	}
	return " WHERE " + strings.Join(clauses, " OR "), values
}

type scan func(dest ...interface{}) error

func prepareColumns(columns es_models.Columns) (string, func(s scan, dest interface{}) error) {
	switch columns {
	case es_models.Columns_Max_Sequence:
		return "SELECT MAX(event_sequence) FROM eventstore.events", func(row scan, dest interface{}) (err error) {
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
	case es_models.Columns_InstanceIDs:
		return "SELECT DISTINCT instance_id FROM eventstore.events", func(row scan, dest interface{}) (err error) {
			instanceID, ok := dest.(*string)
			if !ok {
				return z_errors.ThrowInvalidArgument(nil, "SQL-Fef5h", "type must be *string]")
			}
			err = row(instanceID)
			if err != nil {
				logging.New().WithError(err).Warn("unable to scan row")
				return z_errors.ThrowInternal(err, "SQL-SFef3", "unable to scan row")
			}
			return nil
		}
	case es_models.Columns_Event:
		return selectStmt, func(row scan, dest interface{}) (err error) {
			event, ok := dest.(*es_models.Event)
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
				&event.InstanceID,
				&event.AggregateType,
				&event.AggregateID,
				&event.AggregateVersion,
			)

			if err != nil {
				logging.New().WithError(err).Warn("unable to scan row")
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
	switch operation {
	case es_models.Operation_In:
		return "%s %s ANY(?)"
	case es_models.Operation_NotIn:
		return "%s %s ALL(?)"
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
	case es_models.Field_InstanceID:
		return "instance_id"
	case es_models.Field_EditorService:
		return "editor_service"
	case es_models.Field_EditorUser:
		return "editor_user"
	case es_models.Field_EventType:
		return "event_type"
	case es_models.Field_CreationDate:
		return "creation_date"
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
	case es_models.Operation_NotIn:
		return "<>"
	}
	return ""
}
