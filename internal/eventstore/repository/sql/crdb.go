package sql

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/jackc/pgconn"
	"github.com/lib/pq"
	"github.com/zitadel/logging"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	//as soon as stored procedures are possible in crdb
	// we could move the code to migrations and coll the procedure
	// traking issue: https://github.com/cockroachdb/cockroach/issues/17511
	//
	//previous_data selects the needed data of the latest event of the aggregate
	// and buffers it (crdb inmemory)
	crdbInsert = `WITH previous_data (aggregate_type_cd, resource_owner) AS (
	SELECT agg_type.cd, agg.ro FROM (
		SELECT
			CASE WHEN $10::TIMESTAMPTZ IS NOT NULL
			THEN $10::TIMESTAMPTZ
			ELSE (
				SELECT 
					creation_date cd
				FROM 
					eventstore.events 
				WHERE 
					aggregate_type = $2 
					AND (CASE WHEN $9::TEXT IS NULL THEN instance_id IS NULL ELSE instance_id = $9::TEXT END) 
				ORDER BY
					creation_date DESC
				LIMIT 1)
			END AS cd
			, 1 join_me
	) AS agg_type
	LEFT JOIN (
		SELECT 
			resource_owner ro
			, 1 join_me 
		FROM 
			eventstore.events 
		WHERE 
			aggregate_type = $2
			AND aggregate_id = $3 
			AND (CASE WHEN $9::TEXT IS NULL THEN instance_id IS NULL ELSE instance_id = $9::TEXT END) 
		LIMIT 1
	) AS agg USING(join_me)
)
INSERT INTO eventstore.events (
	event_type,
	aggregate_type,
	aggregate_id,
	aggregate_version,
	event_data,
	editor_user,
	editor_service,
	resource_owner,
	instance_id,
	previous_event_date
) VALUES (
	$1::VARCHAR,
	$2::VARCHAR,
	$3::VARCHAR,
	$4::VARCHAR,
	$5::JSONB,
	$6::VARCHAR,
	$7::VARCHAR,
	CASE WHEN EXISTS(SELECT * FROM previous_data)
		THEN (SELECT COALESCE(resource_owner, $8) FROM previous_data)
		ELSE $8
	END,
	$9::VARCHAR,
	CASE WHEN EXISTS(SELECT * FROM previous_data)
		THEN (SELECT aggregate_type_cd FROM previous_data)
		ELSE NULL
	END
)
RETURNING id, creation_date, resource_owner, instance_id, previous_event_date`

	// " CASE WHEN EXISTS (SELECT * FROM previous_data)" +
	// " THEN (SELECT resource_owner FROM previous_data)" +
	// " ELSE $8::VARCHAR" +
	// " END," +

	uniqueInsert = `INSERT INTO eventstore.unique_constraints
					(
						unique_type,
						unique_field,
						instance_id
					) 
					VALUES (  
						$1,
						$2,
						$3
					)`

	uniqueDelete = `DELETE FROM eventstore.unique_constraints
					WHERE unique_type = $1 and unique_field = $2 and instance_id = $3`
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

		var creationDate sql.NullTime
		for _, event := range events {
			var previousEvent sql.NullTime
			err := tx.QueryRowContext(ctx, crdbInsert,
				event.Type,
				event.AggregateType,
				event.AggregateID,
				event.Version,
				Data(event.Data),
				event.EditorUser,
				event.EditorService,
				event.ResourceOwner,
				event.InstanceID,
				creationDate,
			).Scan(&event.ID, &event.CreationDate, &event.ResourceOwner, &event.InstanceID, &previousEvent)

			if err != nil {
				logging.WithFields(
					"aggregate", event.AggregateType,
					"aggregateId", event.AggregateID,
					"aggregateType", event.AggregateType,
					"eventType", event.Type,
					"instanceID", event.InstanceID,
				).WithError(err).Info("query failed")
				return caos_errs.ThrowInternal(err, "SQL-SBP37", "unable to create event")
			}

			event.PreviousEventDate = previousEvent.Time
			creationDate.Time = event.CreationDate
			creationDate.Valid = true
		}

		err := db.handleUniqueConstraints(ctx, tx, uniqueConstraints...)
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
	if len(uniqueConstraints) == 0 || (len(uniqueConstraints) == 1 && uniqueConstraints[0] == nil) {
		return nil
	}

	for _, uniqueConstraint := range uniqueConstraints {
		uniqueConstraint.UniqueField = strings.ToLower(uniqueConstraint.UniqueField)
		switch uniqueConstraint.Action {
		case repository.UniqueConstraintAdd:
			_, err := tx.ExecContext(ctx, uniqueInsert, uniqueConstraint.UniqueType, uniqueConstraint.UniqueField, uniqueConstraint.InstanceID)
			if err != nil {
				logging.WithFields(
					"unique_type", uniqueConstraint.UniqueType,
					"unique_field", uniqueConstraint.UniqueField).WithError(err).Info("insert unique constraint failed")

				if db.isUniqueViolationError(err) {
					return caos_errs.ThrowAlreadyExists(err, "SQL-M0dsf", uniqueConstraint.ErrorMessage)
				}

				return caos_errs.ThrowInternal(err, "SQL-dM9ds", "unable to create unique constraint ")
			}
		case repository.UniqueConstraintRemoved:
			_, err := tx.ExecContext(ctx, uniqueDelete, uniqueConstraint.UniqueType, uniqueConstraint.UniqueField, uniqueConstraint.InstanceID)
			if err != nil {
				logging.WithFields(
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

// LatestCreationDate returns the latest creation date found by the search query
func (db *CRDB) LatestCreationDate(ctx context.Context, searchQuery *repository.SearchQuery) (creationDate time.Time, err error) {
	err = query(ctx, db, searchQuery, &creationDate)
	return creationDate, err
}

// InstanceIDs returns the instance ids found by the search query
func (db *CRDB) InstanceIDs(ctx context.Context, searchQuery *repository.SearchQuery) ([]string, error) {
	var ids []string
	err := query(ctx, db, searchQuery, &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (db *CRDB) db() *sql.DB {
	return db.client
}

func (db *CRDB) orderByCreationDate(desc bool) string {
	if desc {
		return " ORDER BY creation_date DESC"
	}

	return " ORDER BY creation_date"
}

func (db *CRDB) eventQuery() string {
	return "SELECT" +
		" creation_date" +
		", event_type" +
		", event_data" +
		", editor_service" +
		", editor_user" +
		", resource_owner" +
		", instance_id" +
		", aggregate_type" +
		", aggregate_id" +
		", aggregate_version" +
		", previous_event_date" +
		" FROM eventstore.events" //AS OF SYSTEM TIME '-1ms'::INTERVAL"
}

func (db *CRDB) maxCreationDateQuery() string {
	return "SELECT MAX(creation_date) FROM eventstore.events"
}

func (db *CRDB) instanceIDsQuery() string {
	return "SELECT DISTINCT instance_id FROM eventstore.events AS OF SYSTEM TIME follower_read_timestamp()"
}

func (db *CRDB) columnName(col repository.Field) string {
	switch col {
	case repository.FieldAggregateID:
		return "aggregate_id"
	case repository.FieldAggregateType:
		return "aggregate_type"
	case repository.FieldResourceOwner:
		return "resource_owner"
	case repository.FieldInstanceID:
		return "instance_id"
	case repository.FieldEditorService:
		return "editor_service"
	case repository.FieldEditorUser:
		return "editor_user"
	case repository.FieldEventType:
		return "event_type"
	case repository.FieldEventData:
		return "event_data"
	case repository.FieldCreationDate:
		return "creation_date"
	default:
		panic("invalid column")
	}
}

func (db *CRDB) conditionFormat(operation repository.Operation) string {
	switch operation {
	case repository.OperationIn:
		return "%s %s ANY(?)"
	case repository.OperationNotIn:
		return "%s %s ALL(?)"
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
	case repository.OperationNotIn:
		return "<>"
	}
	return ""
}

var (
	placeholder = regexp.MustCompile(`\?`)
)

// placeholder replaces all "?" with postgres placeholders ($<NUMBER>)
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
	if pgxErr, ok := err.(*pgconn.PgError); ok {
		if pgxErr.Code == "23505" {
			return true
		}
	}
	return false
}
