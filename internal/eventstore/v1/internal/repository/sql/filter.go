package sql

import (
	"context"
	"database/sql"
	"errors"
	"runtime/debug"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	errs "github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type Querier interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func (db *SQL) Filter(ctx context.Context, searchQuery *es_models.SearchQueryFactory) (events []*es_models.Event, err error) {
	if !searchQuery.InstanceFiltered {
		logging.WithFields("stack", string(debug.Stack())).Warn("instanceid not filtered")
	}
	return db.filter(ctx, db.client, searchQuery)
}

func (server *SQL) filter(ctx context.Context, db *database.DB, searchQuery *es_models.SearchQueryFactory) (events []*es_models.Event, err error) {
	query, limit, values, rowScanner := server.buildQuery(ctx, db, searchQuery)
	if query == "" {
		return nil, errs.ThrowInvalidArgument(nil, "SQL-rWeBw", "invalid query factory")
	}

	events = make([]*es_models.Event, 0, limit)
	err = db.QueryContext(ctx,
		func(rows *sql.Rows) error {
			for rows.Next() {
				event := new(es_models.Event)
				err := rowScanner(rows.Scan, event)
				if err != nil {
					return err
				}

				events = append(events, event)
			}
			return nil
		},
		query, values...,
	)
	if err != nil {
		logging.New().WithError(err).Info("query failed")
		return nil, errs.ThrowInternal(err, "SQL-IJuyR", "unable to filter events")
	}
	return events, nil
}

func (db *SQL) LatestSequence(ctx context.Context, queryFactory *es_models.SearchQueryFactory) (uint64, error) {
	query, _, values, rowScanner := db.buildQuery(ctx, db.client, queryFactory)
	if query == "" {
		return 0, errs.ThrowInvalidArgument(nil, "SQL-rWeBw", "invalid query factory")
	}
	sequence := new(Sequence)
	err := db.client.QueryRowContext(ctx, func(row *sql.Row) error {
		return rowScanner(row.Scan, sequence)
	}, query, values...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logging.New().WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Info("query failed")
		return 0, errs.ThrowInternal(err, "SQL-Yczyx", "unable to filter latest sequence")
	}
	return uint64(*sequence), nil
}

func (db *SQL) InstanceIDs(ctx context.Context, queryFactory *es_models.SearchQueryFactory) (ids []string, err error) {
	query, _, values, rowScanner := db.buildQuery(ctx, db.client, queryFactory)
	if query == "" {
		return nil, errs.ThrowInvalidArgument(nil, "SQL-Sfwg2", "invalid query factory")
	}

	err = db.client.QueryContext(ctx,
		func(rows *sql.Rows) error {
			for rows.Next() {
				var id string
				err := rowScanner(rows.Scan, &id)
				if err != nil {
					return err
				}

				ids = append(ids, id)
			}
			return nil
		},
		query, values...)
	if err != nil {
		logging.New().WithError(err).Info("query failed")
		return nil, errs.ThrowInternal(err, "SQL-Sfg3r", "unable to filter instance ids")
	}

	return ids, nil
}
