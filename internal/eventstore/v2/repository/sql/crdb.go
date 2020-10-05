package sql

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strconv"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
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
	client *sql.DB
}

func (db *CRDB) Health(ctx context.Context) error { return db.client.Ping() }

// Push adds all events to the eventstreams of the aggregates.
// This call is transaction save. The transaction will be rolled back if one event fails
func (db *CRDB) Push(ctx context.Context, events ...*repository.Event) error {
	err := crdb.ExecuteTx(ctx, db.client, nil, func(tx *sql.Tx) error {
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
func (db *CRDB) Filter(ctx context.Context, searchQuery *repository.SearchQuery) (events []*repository.Event, err error) {
	rows, rowScanner, err := db.query(searchQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		event := new(repository.Event)
		err := rowScanner(rows.Scan, event)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

//LatestSequence returns the latests sequence found by the the search query
func (db *CRDB) LatestSequence(ctx context.Context, searchQuery *repository.SearchQuery) (uint64, error) {
	rows, rowScanner, err := db.query(searchQuery)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, caos_errs.ThrowNotFound(nil, "SQL-cAEzS", "latest sequence not found")
	}

	var seq Sequence
	err = rowScanner(rows.Scan, &seq)
	if err != nil {
		return 0, err
	}

	return uint64(seq), nil
}

func (db *CRDB) query(searchQuery *repository.SearchQuery) (*sql.Rows, rowScan, error) {
	query, values, rowScanner := buildQuery(db, searchQuery)
	if query == "" {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "SQL-rWeBw", "invalid query factory")
	}

	rows, err := db.client.Query(query, values...)
	if err != nil {
		logging.Log("SQL-HP3Uk").WithError(err).Info("query failed")
		return nil, nil, caos_errs.ThrowInternal(err, "SQL-IJuyR", "unable to filter events")
	}
	return rows, rowScanner, nil
}

func (db *CRDB) eventQuery() string {
	return "SELECT" +
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
}
func (db *CRDB) maxSequenceQuery() string {
	return "SELECT MAX(event_sequence) FROM eventstore.events"
}

func (db *CRDB) columnName(col repository.Field) string {
	switch col {
	case repository.Field_AggregateID:
		return "aggregate_id"
	case repository.Field_AggregateType:
		return "aggregate_type"
	case repository.Field_LatestSequence:
		return "event_sequence"
	case repository.Field_ResourceOwner:
		return "resource_owner"
	case repository.Field_EditorService:
		return "editor_service"
	case repository.Field_EditorUser:
		return "editor_user"
	case repository.Field_EventType:
		return "event_type"
	default:
		return ""
	}
}

func (db *CRDB) conditionFormat(operation repository.Operation) string {
	if operation == repository.Operation_In {
		return "%s %s ANY(?)"
	}
	return "%s %s ?"
}

func (db *CRDB) operation(operation repository.Operation) string {
	switch operation {
	case repository.Operation_Equals, repository.Operation_In:
		return "="
	case repository.Operation_Greater:
		return ">"
	case repository.Operation_Less:
		return "<"
	}
	return ""
}

var (
	placeholder = regexp.MustCompile(`\?`)
)

//placeholder replaces all "?" with postgres placeholders ($<NUMBER>)
func (db *CRDB) placeholder(query string) string {
	occurances := placeholder.FindAllStringIndex(query, -1)
	if len(occurances) == 0 {
		return query
	}
	replaced := query[:occurances[0][0]]

	for i, l := range occurances {
		nextIDX := len(query)
		if i < len(occurances)-1 {
			nextIDX = occurances[i+1][0]
		}
		replaced = replaced + "$" + strconv.Itoa(i+1) + query[l[1]:nextIDX]
	}
	return replaced
}
