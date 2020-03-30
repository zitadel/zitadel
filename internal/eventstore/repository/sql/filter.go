package sql

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	caos_errs "github.com/caos/utils/errors"
	"github.com/caos/utils/logging"
	"github.com/caos/utils/tracing"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

func (db *SQL) Filter(ctx context.Context, searchQuery *es_models.SearchQuery) (events []*es_models.Event, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer span.EndWithError(err)

	query := "SELECT" +
		" id" +
		", creation_date" +
		", event_type" +
		", event_sequence" +
		", previous_sequence" +
		", event_data" +
		", modifier_service" +
		", modifier_tenant" +
		", modifier_user" +
		", resource_owner" +
		", aggregate_type" +
		", aggregate_id" +
		", aggregate_version" +
		" FROM eventstore.events"

	where, values := prepareWhere(searchQuery)
	query += where

	query += " ORDER BY event_sequence"
	if searchQuery.OrderDesc() {
		query += " DESC"
	}

	if searchQuery.Limit() > 0 {
		values = append(values, searchQuery.Limit())
		query += " LIMIT ?"
	}

	query = numberPlaceholder(query, "?", "$")

	rows, err := db.client.Query(query, values...)
	if err != nil {
		logging.Log("SQL-HP3Uk").WithError(err).Info("query failed")
		return nil, caos_errs.ThrowInternal(err, "SQL-IJuyR", "unable to filter events")
	}

	for rows.Next() {
		event := new(es_models.Event)
		rows.Scan(
			&event.ID,
			&event.CreationDate,
			&event.Typ,
			&event.Sequence,
			&event.PreviousSequence,
			&event.Data,
			&event.ModifierService,
			&event.ModifierTenant,
			&event.ModifierUser,
			&event.ResourceOwner,
			&event.AggregateType,
			&event.AggregateID,
			&event.AggregateVersion,
		)
		events = append(events, event)
	}

	return events, err
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
	values = make([]interface{}, len(searchQuery.Filters()))
	clauses := make([]string, len(searchQuery.Filters()))

	if len(values) == 0 {
		return clause, values
	}

	for i, filter := range searchQuery.Filters() {
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

func generateFilters(query *gorm.DB, searchQuery es_models.SearchQuery) *gorm.DB {
	for _, f := range searchQuery.Filters() {
		value := f.GetValue()
		switch value.(type) {
		case []bool, []float64, []int64, []string, *[]bool, *[]float64, *[]int64, *[]string:
			value = pq.Array(value)
		}
		query = query.Where(getCondition(f), value)
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
	if operation == es_models.In {
		return "%s %s (?)"
	}
	return "%s %s ?"
}

func getField(field es_models.Field) string {
	switch field {
	case es_models.AggregateID:
		return "aggregate_id"
	case es_models.AggregateType:
		return "aggregate_type"
	case es_models.LatestSequence:
		return "event_sequence"
	case es_models.ResourceOwner:
		return "resource_owner"
	case es_models.ModifierService:
		return "modifier_service"
	case es_models.ModifierUser:
		return "modifier_user"
	case es_models.ModifierTenant:
		return "modifier_tenant"
	}
	return ""
}

func getOperation(operation es_models.Operation) string {
	switch operation {
	case es_models.Equals:
		return "="
	case es_models.Greater:
		return ">"
	case es_models.Less:
		return "<"
	case es_models.In:
		return "IN"
	}
	return ""
}
