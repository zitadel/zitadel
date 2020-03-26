package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/caos/eventstore-lib/pkg/models"
	caos_errs "github.com/caos/utils/errors"
	"github.com/caos/utils/logging"
	"github.com/caos/utils/tracing"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

func (db *SQL) Filter(ctx context.Context, events models.Events, searchQuery models.SearchQuery) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer span.EndWithError(err)

	query := "SELECT" +
		" id," +
		" creation_date," +
		" event_type," +
		" event_sequence," +
		" previous_sequence," +
		" event_data," +
		" modifier_service," +
		" modifier_tenant," +
		" modifier_user," +
		" resource_owner," +
		" aggregate_type," +
		" aggregate_id" +
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
		return caos_errs.ThrowInternal(err, "SQL-IJuyR", "unable to filter events")
	}

	for rows.Next() {
		event := new(Event)
		rows.Scan(
			&event.ID,
			&event.CreationDate,
			&event.Typ,
			&event.Sequence,
			&event.PreviousSequence,
			&event.Data,
			&event.ModifierService,
			&event.ModifierTenant,
			&event.ModiferUser,
			&event.ResourceOwner,
			&event.AggregateType,
			&event.AggregateID,
		)
		events.Append(eventToApp(event))
	}

	return nil
}

func numberPlaceholder(query, old, new string) string {
	for i, hasChanged := 1, true; hasChanged; i++ {
		newQuery := strings.Replace(query, old, new+strconv.Itoa(i), 1)
		hasChanged = query != newQuery
		query = newQuery
	}
	return query
}

func prepareWhere(searchQuery models.SearchQuery) (clause string, values []interface{}) {
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

func generateFilters(query *gorm.DB, searchQuery models.SearchQuery) *gorm.DB {
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

func getCondition(filter models.Filter) string {
	field := getField(filter.GetField())
	operation := getOperation(filter.GetOperation())
	format := prepareConditionFormat(filter.GetOperation())

	return fmt.Sprintf(format, field, operation)
}

func prepareConditionFormat(operation models.Operation) string {
	if operation == es_models.In {
		return "%s %s (?)"
	}
	return "%s %s ?"
}

func getField(field models.Field) string {
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

func getOperation(operation models.Operation) string {
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
