package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/cockroachdb/cockroach-go/v2/crdb"

	//sql import for cockroach
	_ "github.com/lib/pq"
)

const (
	crdbInsert = "WITH input_event ( " +
		"    event_type, " +
		"    aggregate_type, " +
		"    aggregate_id, " +
		"    aggregate_version, " +
		"    creation_date, " +
		"    event_data, " +
		"    editor_user, " +
		"    editor_service, " +
		"    resource_owner, " +
		"    previous_sequence, " +
		"    check_previous, " +
		// variables below are calculated
		"    max_event_seq, " +
		"    event_count " +
		") " +
		"	AS( " +
		"		SELECT " +
		"			$1::VARCHAR," +
		"			$2::VARCHAR," +
		"			$3::VARCHAR," +
		"			$4::VARCHAR," +
		"			COALESCE($5::TIMESTAMPTZ, NOW()), " +
		"			$6::JSONB, " +
		"			$7::VARCHAR, " +
		"			$8::VARCHAR, " +
		"			$9::VARCHAR, " +
		"			$10::BIGINT, " +
		"			$11::BOOLEAN," +
		"			MAX(event_sequence) AS max_event_seq, " +
		"			COUNT(*) AS event_count " +
		"	FROM eventstore.events " +
		"	WHERE " +
		"		aggregate_type = $2::VARCHAR " +
		"		AND aggregate_id = $3::VARCHAR " +
		") " +
		"INSERT INTO eventstore.events " +
		"	( " +
		"		event_type, " +
		"		aggregate_type," +
		"		aggregate_id, " +
		"		aggregate_version, " +
		"		creation_date, " +
		"		event_data, " +
		"		editor_user, " +
		"		editor_service, " +
		"		resource_owner, " +
		"		previous_sequence " +
		"	) " +
		"	( " +
		"		SELECT " +
		"			event_type, " +
		"			aggregate_type," +
		"			aggregate_id, " +
		"			aggregate_version, " +
		"			COALESCE(creation_date, NOW()), " +
		"			event_data, " +
		"			editor_user, " +
		"			editor_service, " +
		"			resource_owner, " +
		"			( " +
		"			    SELECT " +
		"			        CASE " +
		"			            WHEN NOT check_previous THEN " +
		"			                max_event_seq " +
		"			            ELSE " +
		"			                previous_sequence " +
		"			        END" +
		"			) " +
		"		FROM input_event " +
		"		WHERE EXISTS ( " +
		"		    SELECT " +
		"		        CASE " +
		"		            WHEN NOT check_previous THEN 1 " +
		"		            ELSE ( " +
		"		                SELECT 1 FROM input_event " +
		"		                    WHERE max_event_seq = previous_sequence OR (previous_sequence IS NULL AND event_count = 0) " +
		"		            ) " +
		"		        END " +
		"		) " +
		"	) " +
		"RETURNING event_sequence, creation_date "
)

type CRDB struct {
	db *sql.DB
}

func (db *CRDB) Health(ctx context.Context) error { return db.db.Ping() }

// Push adds all events to the eventstreams of the aggregates.
// This call is transaction save. The transaction will be rolled back if one event fails
func (db *CRDB) Push(ctx context.Context, events ...*Event) error {
	err := crdb.ExecuteTx(ctx, db.db, nil, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, crdbInsert)
		if err != nil {
			tx.Rollback()
			logging.Log("SQL-3to5p").WithError(err).Warn("prepare failed")
			return caos_errs.ThrowInternal(err, "SQL-OdXRE", "prepare failed")
		}
		for _, event := range events {
			previousSequence := event.PreviousSequence
			if event.PreviousEvent != nil {
				previousSequence = event.PreviousSequence
			}
			err = stmt.QueryRowContext(ctx,
				event.Type,
				event.AggregateType,
				event.AggregateID,
				event.Version,
				event.CreationDate,
				event.Data,
				event.EditorUser,
				event.EditorService,
				event.ResourceOwner,
				previousSequence,
				event.CheckPreviousSequence,
			).Scan(&event.Sequence, &event.CreationDate)

			if err != nil {
				tx.Rollback()

				logging.LogWithFields("SQL-IP3js",
					"aggregate", event.AggregateType,
					"aggregateId", event.AggregateID,
					"aggregateType", event.AggregateType,
					"eventType", event.Type).WithError(err).Info("query failed")
				return caos_errs.ThrowInternal(err, "SQL-SBP37", "unable to create event")
			}
		}
		return nil
	})
	if err != nil && !errors.Is(err, &caos_errs.CaosError{}) {
		err = caos_errs.ThrowInternal(err, "SQL-DjgtG", "unable to store events")
	}

	return err
}

// Filter returns all events matching the given search query
// func (db *CRDB) Filter(ctx context.Context, searchQuery *SearchQuery) (events []*Event, err error) {

// 	return events, nil
// }

//LatestSequence returns the latests sequence found by the the search query
func (db *CRDB) LatestSequence(ctx context.Context, queryFactory *SearchQuery) (uint64, error) {
	return 0, nil
}

func (db *CRDB) prepareQuery(columns Columns) (string, func(s scanner, dest interface{}) error) {
	switch columns {
	case Columns_Max_Sequence:
		return "SELECT MAX(event_sequence) FROM eventstore.events", func(scan scanner, dest interface{}) (err error) {
			sequence, ok := dest.(*Sequence)
			if !ok {
				return caos_errs.ThrowInvalidArgument(nil, "SQL-NBjA9", "type must be sequence")
			}
			err = scan(sequence)
			if err == nil || errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return caos_errs.ThrowInternal(err, "SQL-bN5xg", "something went wrong")
		}
	case Columns_Event:
		return selectStmt, func(row scanner, dest interface{}) (err error) {
			event, ok := dest.(*Event)
			if !ok {
				return caos_errs.ThrowInvalidArgument(nil, "SQL-4GP6F", "type must be event")
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
				&event.Version,
			)

			if err != nil {
				logging.Log("SQL-kn1Sw").WithError(err).Warn("unable to scan row")
				return caos_errs.ThrowInternal(err, "SQL-J0hFS", "unable to scan row")
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

func (db *CRDB) prepareFilter(filters []*Filter) string {
	filter := ""
	// for _, f := range filters{
	// 	f.
	// }
	return filter
}
