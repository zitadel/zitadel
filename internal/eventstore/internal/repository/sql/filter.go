package sql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
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
)

type Querier interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func (db *SQL) Filter(ctx context.Context, searchQuery *es_models.SearchQueryFactory) (events []*models.Event, err error) {
	return filter(db.client, searchQuery)
}

func filter(querier Querier, searchQuery *es_models.SearchQueryFactory) (events []*es_models.Event, err error) {
	query, limit, values, rowScanner := buildQuery(searchQuery)

	rows, err := querier.Query(query, values...)
	if err != nil {
		logging.Log("SQL-HP3Uk").WithError(err).Info("query failed")
		return nil, errors.ThrowInternal(err, "SQL-IJuyR", "unable to filter events")
	}
	defer rows.Close()

	events = make([]*es_models.Event, 0, limit)

	for rows.Next() {
		event, err := rowScanner(rows)
		if err != nil {
			return nil, err
		}

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
