package sql

import (
	"context"
	"strings"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/lib/pq"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
)

func (db *SQL) LatestSequence(ctx context.Context, queryFactory *es_models.SearchQueryFactory) (uint64, error) {
	sequenceFactory := *queryFactory

	sequenceFactory = *(&sequenceFactory).Columns(es_models.Columns_Max_Sequence)
	sequenceFactory = *(&sequenceFactory).SequenceGreater(0)

	query, _, values, rowScanner := buildQuery(&sequenceFactory)
	row := db.client.QueryRow(query, values...)
	event, err := rowScanner(row)
	if err != nil {
		logging.Log("SQL-IXjUN").WithError(err).Info("query failed")
		return 0, errors.ThrowInternal(err, "SQL-WMIAq", "unable to filter events")
	}
	return event.Sequence, nil
}

func buildQuery(queryFactory *es_models.SearchQueryFactory) (query string, limit uint64, values []interface{}, rowScanner func(s scanner) (*es_models.Event, error)) {
	searchQuery := queryFactory.Build()
	query, rowScanner = prepareColumns(searchQuery.Columns)
	where, values := prepareCondition(searchQuery.Limit, searchQuery.Desc, searchQuery.Filters)
	query += where

	if searchQuery.Columns != es_models.Columns_Max_Sequence {
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

func prepareCondition(limit uint64, isDesc bool, filters []*es_models.Filter) (clause string, values []interface{}) {
	values = make([]interface{}, len(filters))
	clauses := make([]string, len(filters))

	if len(values) == 0 {
		return clause, values
	}

	for i, filter := range filters {
		value := filter.GetValue()
		switch value.(type) {
		case []bool, []float64, []int64, []string, []es_models.AggregateType, []es_models.EventType, *[]bool, *[]float64, *[]int64, *[]string, *[]es_models.AggregateType, *[]es_models.EventType:
			value = pq.Array(value)
		}

		clauses[i] = getCondition(filter)
		values[i] = value
	}
	return " WHERE " + strings.Join(clauses, " AND "), values
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func prepareColumns(columns es_models.Columns) (string, func(s scanner) (*es_models.Event, error)) {
	switch columns {
	case es_models.Columns_Max_Sequence:
		return "SELECT MAX(event_sequence) FROM eventstore.events", func(row scanner) (event *es_models.Event, err error) {
			var sequence Sequence
			err = row.Scan(&sequence)
			return &es_models.Event{Sequence: uint64(sequence)}, err
		}
	default:
		return selectStmt, func(row scanner) (event *es_models.Event, err error) {
			event = new(es_models.Event)
			var previousSequence Sequence
			data := make(Data, 0)

			err = row.Scan(
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
				return nil, errors.ThrowInternal(err, "SQL-J0hFS", "unable to scan row")
			}

			event.PreviousSequence = uint64(previousSequence)

			event.Data = make([]byte, len(data))
			copy(event.Data, data)

			return event, err
		}
	}
}
