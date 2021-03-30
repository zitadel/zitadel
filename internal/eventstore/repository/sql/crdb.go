package sql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"regexp"
	"strconv"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/cockroachdb/cockroach-go/v2/crdb"

	//sql import for cockroach
	_ "github.com/lib/pq"
)

const (
	//as soon as stored procedures are possible in crdb
	// we could move the code to migrations and coll the procedure
	// traking issue: https://github.com/cockroachdb/cockroach/issues/17511
	crdbInsert = "WITH data ( " +
		"    event_type, " +
		"    aggregate_type, " +
		"    aggregate_id, " +
		"    aggregate_version, " +
		"    creation_date, " +
		"    event_data, " +
		"    editor_user, " +
		"    editor_service, " +
		"    resource_owner, " +
		// variables below are calculated
		"    previous_sequence" +
		") AS (" +
		//previous_data selects the needed data of the latest event of the aggregate
		// and buffers it (crdb inmemory)
		"    WITH previous_data AS (" +
		"        SELECT MAX(event_sequence) AS seq, resource_owner " +
		"        FROM eventstore.events " +
		//TODO: remove LIMIT 1 / order by as soon as data cleaned up (only 1 resource_owner per aggregate)
		"        WHERE aggregate_type = $2 AND aggregate_id = $3 GROUP BY resource_owner order by seq desc LIMIT 1" +
		"    )" +
		// defines the data to be inserted
		"    SELECT " +
		"        $1::VARCHAR AS event_type, " +
		"        $2::VARCHAR AS aggregate_type, " +
		"        $3::VARCHAR AS aggregate_id, " +
		"        $4::VARCHAR AS aggregate_version, " +
		"        NOW() AS creation_date, " +
		"        $5::JSONB AS event_data, " +
		"        $6::VARCHAR AS editor_user, " +
		"        $7::VARCHAR AS editor_service, " +
		"        CASE WHEN EXISTS (SELECT * FROM previous_data) " +
		"            THEN (SELECT resource_owner FROM previous_data) " +
		"            ELSE $8::VARCHAR " +
		"        end AS resource_owner, " +
		"        CASE WHEN EXISTS (SELECT * FROM previous_data) " +
		"            THEN (SELECT seq FROM previous_data) " +
		"            ELSE NULL " +
		"        end AS previous_sequence" +
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
		"			previous_sequence " +
		"		FROM data " +
		"	) " +
		"RETURNING id, event_sequence, previous_sequence, creation_date, resource_owner"
	uniqueInsert = `INSERT INTO eventstore.unique_constraints
					(
						unique_type,
						unique_field
					) 
					VALUES (  
						$1,
						$2
					)`

	uniqueDelete = `DELETE FROM eventstore.unique_constraints
					WHERE unique_type = $1 and unique_field = $2`
)

type CRDB struct {
	client *sql.DB
}

func NewCRDB(client *sql.DB) *CRDB {
	return &CRDB{client}
}

func (db *CRDB) Health(ctx context.Context) error { return db.client.Ping() }

// Push adds all events to the eventstreams of the aggregates.
// This call is transaction save. The transaction will be rolled back if one event fails
func (db *CRDB) Push(ctx context.Context, events []*repository.Event, uniqueConstraints ...*repository.UniqueConstraint) error {
	err := crdb.ExecuteTx(ctx, db.client, nil, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, crdbInsert)
		if err != nil {
			logging.Log("SQL-3to5p").WithError(err).Warn("prepare failed")
			return caos_errs.ThrowInternal(err, "SQL-OdXRE", "prepare failed")
		}

		var previousSequence Sequence
		for _, event := range events {
			err = stmt.QueryRowContext(ctx,
				event.Type,
				event.AggregateType,
				event.AggregateID,
				event.Version,
				Data(event.Data),
				event.EditorUser,
				event.EditorService,
				event.ResourceOwner,
			).Scan(&event.ID, &event.Sequence, &previousSequence, &event.CreationDate, &event.ResourceOwner)

			event.PreviousSequence = uint64(previousSequence)

			if err != nil {
				logging.LogWithFields("SQL-IP3js",
					"aggregate", event.AggregateType,
					"aggregateId", event.AggregateID,
					"aggregateType", event.AggregateType,
					"eventType", event.Type).WithError(err).Info("query failed",
					"seq", event.PreviousSequence)
				return caos_errs.ThrowInternal(err, "SQL-SBP37", "unable to create event")
			}
		}

		err = db.handleUniqueConstraints(ctx, tx, uniqueConstraints...)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil && !errors.Is(err, &caos_errs.CaosError{}) {
		err = caos_errs.ThrowInternal(err, "SQL-DjgtG", "unable to store events")
	}

	return err
}

// handleUniqueConstraints adds or removes unique constraints
func (db *CRDB) handleUniqueConstraints(ctx context.Context, tx *sql.Tx, uniqueConstraints ...*repository.UniqueConstraint) (err error) {
	if uniqueConstraints == nil || len(uniqueConstraints) == 0 || (len(uniqueConstraints) == 1 && uniqueConstraints[0] == nil) {
		return nil
	}

	for _, uniqueConstraint := range uniqueConstraints {
		if uniqueConstraint.Action == repository.UniqueConstraintAdd {
			_, err := tx.ExecContext(ctx, uniqueInsert, uniqueConstraint.UniqueType, uniqueConstraint.UniqueField)
			if err != nil {
				logging.LogWithFields("SQL-IP3js",
					"unique_type", uniqueConstraint.UniqueType,
					"unique_field", uniqueConstraint.UniqueField).WithError(err).Info("insert unique constraint failed")

				if db.isUniqueViolationError(err) {
					return caos_errs.ThrowAlreadyExists(err, "SQL-M0dsf", uniqueConstraint.ErrorMessage)
				}

				return caos_errs.ThrowInternal(err, "SQL-dM9ds", "unable to create unique constraint ")
			}
		} else if uniqueConstraint.Action == repository.UniqueConstraintRemoved {
			_, err := tx.ExecContext(ctx, uniqueDelete, uniqueConstraint.UniqueType, uniqueConstraint.UniqueField)
			if err != nil {
				logging.LogWithFields("SQL-M0vsf",
					"unique_type", uniqueConstraint.UniqueType,
					"unique_field", uniqueConstraint.UniqueField).WithError(err).Info("delete unique constraint failed")
				return caos_errs.ThrowInternal(err, "SQL-6n88i", "unable to remove unique constraint ")
			}
		}
	}
	return nil
}

// Filter returns all events matching the given search query
func (db *CRDB) Filter(ctx context.Context, searchQuery *repository.SearchQuery) (events []*repository.Event, err error) {
	events = []*repository.Event{}
	err = query(ctx, db, searchQuery, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

//LatestSequence returns the latests sequence found by the the search query
func (db *CRDB) LatestSequence(ctx context.Context, searchQuery *repository.SearchQuery) (uint64, error) {
	var seq Sequence
	err := query(ctx, db, searchQuery, &seq)
	if err != nil {
		return 0, err
	}
	return uint64(seq), nil
}

func (db *CRDB) db() *sql.DB {
	return db.client
}

func (db *CRDB) orderByEventSequence(desc bool) string {
	if desc {
		return " ORDER BY event_sequence DESC"
	}

	return " ORDER BY event_sequence"
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
	case repository.FieldAggregateID:
		return "aggregate_id"
	case repository.FieldAggregateType:
		return "aggregate_type"
	case repository.FieldSequence:
		return "event_sequence"
	case repository.FieldResourceOwner:
		return "resource_owner"
	case repository.FieldEditorService:
		return "editor_service"
	case repository.FieldEditorUser:
		return "editor_user"
	case repository.FieldEventType:
		return "event_type"
	case repository.FieldEventData:
		return "event_data"
	default:
		return ""
	}
}

func (db *CRDB) conditionFormat(operation repository.Operation) string {
	if operation == repository.OperationIn {
		return "%s %s ANY(?)"
	}
	return "%s %s ?"
}

func (db *CRDB) operation(operation repository.Operation) string {
	switch operation {
	case repository.OperationEquals, repository.OperationIn:
		return "="
	case repository.OperationGreater:
		return ">"
	case repository.OperationLess:
		return "<"
	case repository.OperationJSONContains:
		return "@>"
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

func (db *CRDB) isUniqueViolationError(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "23505" {
			return true
		}
	}
	return false
}
