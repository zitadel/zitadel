package sql

import (
	"context"
	"database/sql"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type Querier interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func (db *SQL) Filter(ctx context.Context, searchQuery *es_models.SearchQueryFactory) (events []*es_models.Event, err error) {
	return filter(db.client, searchQuery)
}

func filter(querier Querier, searchQuery *es_models.SearchQueryFactory) (events []*es_models.Event, err error) {
	query, limit, values, rowScanner := buildQuery(searchQuery)
	if query == "" {
		return nil, errors.ThrowInvalidArgument(nil, "SQL-rWeBw", "invalid query factory")
	}

	rows, err := querier.Query(query, values...)
	if err != nil {
		logging.Log("SQL-HP3Uk").WithError(err).Info("query failed")
		return nil, errors.ThrowInternal(err, "SQL-IJuyR", "unable to filter events")
	}
	defer rows.Close()

	events = make([]*es_models.Event, 0, limit)

	for rows.Next() {
		event := new(es_models.Event)
		err := rowScanner(rows.Scan, event)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (db *SQL) LatestSequence(ctx context.Context, queryFactory *es_models.SearchQueryFactory) (uint64, error) {
	query, _, values, rowScanner := buildQuery(queryFactory)
	if query == "" {
		return 0, errors.ThrowInvalidArgument(nil, "SQL-rWeBw", "invalid query factory")
	}
	row := db.client.QueryRow(query, values...)
	sequence := new(Sequence)
	err := rowScanner(row.Scan, sequence)
	if err != nil {
		logging.Log("SQL-WsxTg").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Info("query failed")
		return 0, errors.ThrowInternal(err, "SQL-Yczyx", "unable to filter latest sequence")
	}
	return uint64(*sequence), nil
}
