package sql

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/lib/pq"
)

const (
	selectStmt = "SELECT" +
		" id" +
		", creation_date" +
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
)

func (db *SQL) Filter(ctx context.Context, searchQuery *es_models.SearchQuery) (events []*models.Event, err error) {
	where, values := prepareWhere(searchQuery)
	query := selectStmt + where

	query += " ORDER BY event_sequence"
	if searchQuery.Desc {
		query += " DESC"
	}

	if searchQuery.Limit > 0 {
		values = append(values, searchQuery.Limit)
		query += " LIMIT ?"
	}

	query = numberPlaceholder(query, "?", "$")

	rows, err := db.client.Query(query, values...)
	if err != nil {
		logging.Log("SQL-HP3Uk").WithError(err).Info("query failed")
		return nil, errors.ThrowInternal(err, "SQL-IJuyR", "unable to filter events")
	}
	defer rows.Close()

	events = make([]*es_models.Event, 0, searchQuery.Limit)

	for rows.Next() {
		event := new(models.Event)
		var previousSequence Sequence

		err = rows.Scan(
			&event.ID,
			&event.CreationDate,
			&event.Type,
			&event.Sequence,
			&previousSequence,
			&event.Data,
			&event.EditorService,
			&event.EditorUser,
			&event.ResourceOwner,
			&event.AggregateType,
			&event.AggregateID,
			&event.AggregateVersion,
		)

		if err != nil {
			logging.Log("SQL-wHNPo").WithError(err).Warn("unable to scan row")
			return nil, errors.ThrowInternal(err, "SQL-BfZwF", "unable to scan row")
		}

		event.PreviousSequence = uint64(previousSequence)
		events = append(events, event)
	}

	return events, nil
}

func numberPlaceholder(query, old, new string) string {
	for i, hasChanged := 1, true; hasChanged; i++ {
		newQuery := strings.Replace(query, old, new+strconv.Itoa(i), 1)
		hasChanged = query != newQuery
		query = newQuery
	}
	return query
}

func prepareWhere(searchQuery *es_models.SearchQuery) (clause string, values []interface{}) {
	values = make([]interface{}, len(searchQuery.Filters))
	clauses := make([]string, len(searchQuery.Filters))

	if len(values) == 0 {
		return clause, values
	}

	for i, filter := range searchQuery.Filters {
		value := filter.GetValue()
		switch value.(type) {
		case []bool, []float64, []int64, []string, *[]bool, *[]float64, *[]int64, *[]string:
			value = pq.Array(value)
		}

		clauses[i] = getCondition(filter)
		values[i] = value
	}
	return " WHERE " + strings.Join(clauses, " AND "), values
}

func getCondition(filter *es_models.Filter) string {
	field := getField(filter.GetField())
	operation := getOperation(filter.GetOperation())
	format := prepareConditionFormat(filter.GetOperation())

	return fmt.Sprintf(format, field, operation)
}

func prepareConditionFormat(operation es_models.Operation) string {
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
