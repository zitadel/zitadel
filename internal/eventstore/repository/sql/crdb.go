package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	//as soon as stored procedures are possible in crdb
	// we could move the code to migrations and call the procedure
	// traking issue: https://github.com/cockroachdb/cockroach/issues/17511
	//
	//previous_data selects the needed data of the latest event of the aggregate
	// and buffers it (crdb inmemory)
	crdbInsert = "WITH previous_data (aggregate_type_sequence, aggregate_sequence, resource_owner) AS (" +
		"SELECT agg_type.seq, agg.seq, agg.ro FROM " +
		"(" +
		//max sequence of requested aggregate type
		" SELECT MAX(event_sequence) seq, 1 join_me" +
		" FROM eventstore.events" +
		" WHERE aggregate_type = $2" +
		" AND (CASE WHEN $9::TEXT IS NULL THEN instance_id is null else instance_id = $9::TEXT END)" +
		") AS agg_type " +
		// combined with
		"LEFT JOIN " +
		"(" +
		// max sequence and resource owner of aggregate root
		" SELECT event_sequence seq, resource_owner ro, 1 join_me" +
		" FROM eventstore.events" +
		" WHERE aggregate_type = $2 AND aggregate_id = $3" +
		" AND (CASE WHEN $9::TEXT IS NULL THEN instance_id is null else instance_id = $9::TEXT END)" +
		" ORDER BY event_sequence DESC" +
		" LIMIT 1" +
		") AS agg USING(join_me)" +
		") " +
		"INSERT INTO eventstore.events (" +
		" event_type," +
		" aggregate_type," +
		" aggregate_id," +
		" aggregate_version," +
		" creation_date," +
		" position," +
		" event_data," +
		" editor_user," +
		" editor_service," +
		" resource_owner," +
		" instance_id," +
		" event_sequence," +
		" previous_aggregate_sequence," +
		" previous_aggregate_type_sequence," +
		" in_tx_order" +
		") " +
		// defines the data to be inserted
		"SELECT" +
		" $1::VARCHAR AS event_type," +
		" $2::VARCHAR AS aggregate_type," +
		" $3::VARCHAR AS aggregate_id," +
		" $4::VARCHAR AS aggregate_version," +
		" hlc_to_timestamp(cluster_logical_timestamp()) AS creation_date," +
		" cluster_logical_timestamp() AS position," +
		" $5::JSONB AS event_data," +
		" $6::VARCHAR AS editor_user," +
		" $7::VARCHAR AS editor_service," +
		" COALESCE((resource_owner), $8::VARCHAR) AS resource_owner," +
		" $9::VARCHAR AS instance_id," +
		" COALESCE(aggregate_sequence, 0)+1," +
		" aggregate_sequence AS previous_aggregate_sequence," +
		" aggregate_type_sequence AS previous_aggregate_type_sequence," +
		" $10 AS in_tx_order " +
		"FROM previous_data " +
		"RETURNING id, event_sequence, creation_date, resource_owner, instance_id"

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
	uniqueDeleteInstance = `DELETE FROM eventstore.unique_constraints
					WHERE instance_id = $1`
)

// awaitOpenTransactions ensures event ordering, so we don't events younger that open transactions
var (
	awaitOpenTransactionsV1 string
	awaitOpenTransactionsV2 string
)

func awaitOpenTransactions(useV1 bool) string {
	if useV1 {
		return awaitOpenTransactionsV1
	}
	return awaitOpenTransactionsV2
}

type CRDB struct {
	*database.DB
}

func NewCRDB(client *database.DB) *CRDB {
	switch client.Type() {
	case "cockroach":
		awaitOpenTransactionsV1 = " AND creation_date::TIMESTAMP < (SELECT COALESCE(MIN(start), NOW())::TIMESTAMP FROM crdb_internal.cluster_transactions where application_name = ANY(?))"
		awaitOpenTransactionsV2 = ` AND hlc_to_timestamp("position") < (SELECT COALESCE(MIN(start), NOW())::TIMESTAMP FROM crdb_internal.cluster_transactions where application_name = ANY(?))`
	case "postgres":
		awaitOpenTransactionsV1 = ` AND EXTRACT(EPOCH FROM created_at) < (SELECT COALESCE(EXTRACT(EPOCH FROM min(xact_start)), EXTRACT(EPOCH FROM now())) FROM pg_stat_activity WHERE datname = current_database() AND application_name = ANY(?) AND state <> 'idle')`
		awaitOpenTransactionsV2 = ` AND "position" < (SELECT COALESCE(EXTRACT(EPOCH FROM min(xact_start)), EXTRACT(EPOCH FROM now())) FROM pg_stat_activity WHERE datname = current_database() AND application_name = ANY(?) AND state <> 'idle')`
	}

	return &CRDB{client}
}

func (db *CRDB) Health(ctx context.Context) error { return db.Ping() }

// Push adds all events to the eventstreams of the aggregates.
// This call is transaction save. The transaction will be rolled back if one event fails
func (db *CRDB) Push(ctx context.Context, commands ...eventstore.Command) (events []eventstore.Event, err error) {
	events = make([]eventstore.Event, len(commands))

	err = crdb.ExecuteTx(ctx, db.DB.DB, nil, func(tx *sql.Tx) error {

		var uniqueConstraints []*eventstore.UniqueConstraint

		for i, command := range commands {
			if command.Aggregate().InstanceID == "" {
				command.Aggregate().InstanceID = authz.GetInstance(ctx).InstanceID()
			}

			var payload []byte
			if command.Payload() != nil {
				payload, err = json.Marshal(command.Payload())
				if err != nil {
					return err
				}
			}
			e := &repository.Event{
				Typ:           command.Type(),
				Data:          payload,
				EditorUser:    command.Creator(),
				Version:       command.Aggregate().Version,
				AggregateID:   command.Aggregate().ID,
				AggregateType: command.Aggregate().Type,
				ResourceOwner: sql.NullString{String: command.Aggregate().ResourceOwner, Valid: command.Aggregate().ResourceOwner != ""},
				InstanceID:    command.Aggregate().InstanceID,
			}

			err := tx.QueryRowContext(ctx, crdbInsert,
				e.Type(),
				e.Aggregate().Type,
				e.Aggregate().ID,
				e.Aggregate().Version,
				payload,
				e.Creator(),
				"zitadel",
				e.Aggregate().ResourceOwner,
				e.Aggregate().InstanceID,
				i,
			).Scan(&e.ID, &e.Seq, &e.CreationDate, &e.ResourceOwner, &e.InstanceID)

			if err != nil {
				logging.WithFields(
					"aggregate", e.Aggregate().Type,
					"aggregateId", e.Aggregate().ID,
					"aggregateType", e.Aggregate().Type,
					"eventType", e.Type(),
					"instanceID", e.Aggregate().InstanceID,
				).WithError(err).Debug("query failed")
				return zerrors.ThrowInternal(err, "SQL-SBP37", "unable to create event")
			}

			uniqueConstraints = append(uniqueConstraints, command.UniqueConstraints()...)
			events[i] = e
		}

		return db.handleUniqueConstraints(ctx, tx, uniqueConstraints...)
	})
	if err != nil && !errors.Is(err, &zerrors.ZitadelError{}) {
		err = zerrors.ThrowInternal(err, "SQL-DjgtG", "unable to store events")
	}

	return events, err
}

// handleUniqueConstraints adds or removes unique constraints
func (db *CRDB) handleUniqueConstraints(ctx context.Context, tx *sql.Tx, uniqueConstraints ...*eventstore.UniqueConstraint) (err error) {
	if len(uniqueConstraints) == 0 || (len(uniqueConstraints) == 1 && uniqueConstraints[0] == nil) {
		return nil
	}

	for _, uniqueConstraint := range uniqueConstraints {
		uniqueConstraint.UniqueField = strings.ToLower(uniqueConstraint.UniqueField)
		switch uniqueConstraint.Action {
		case eventstore.UniqueConstraintAdd:
			_, err := tx.ExecContext(ctx, uniqueInsert, uniqueConstraint.UniqueType, uniqueConstraint.UniqueField, authz.GetInstance(ctx).InstanceID())
			if err != nil {
				logging.WithFields(
					"unique_type", uniqueConstraint.UniqueType,
					"unique_field", uniqueConstraint.UniqueField).WithError(err).Info("insert unique constraint failed")

				if db.isUniqueViolationError(err) {
					return zerrors.ThrowAlreadyExists(err, "SQL-wHcEq", uniqueConstraint.ErrorMessage)
				}

				return zerrors.ThrowInternal(err, "SQL-dM9ds", "unable to create unique constraint")
			}
		case eventstore.UniqueConstraintRemove:
			_, err := tx.ExecContext(ctx, uniqueDelete, uniqueConstraint.UniqueType, uniqueConstraint.UniqueField, authz.GetInstance(ctx).InstanceID())
			if err != nil {
				logging.WithFields(
					"unique_type", uniqueConstraint.UniqueType,
					"unique_field", uniqueConstraint.UniqueField).WithError(err).Info("delete unique constraint failed")
				return zerrors.ThrowInternal(err, "SQL-6n88i", "unable to remove unique constraint")
			}
		case eventstore.UniqueConstraintInstanceRemove:
			_, err := tx.ExecContext(ctx, uniqueDeleteInstance, authz.GetInstance(ctx).InstanceID())
			if err != nil {
				logging.WithFields(
					"instance_id", authz.GetInstance(ctx).InstanceID()).WithError(err).Info("delete instance unique constraints failed")
				return zerrors.ThrowInternal(err, "SQL-6n88i", "unable to remove unique constraints of instance")
			}
		}
	}
	return nil
}

// FilterToReducer finds all events matching the given search query and passes them to the reduce function.
func (crdb *CRDB) FilterToReducer(ctx context.Context, searchQuery *eventstore.SearchQueryBuilder, reduce eventstore.Reducer) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = query(ctx, crdb, searchQuery, reduce, false)
	if err == nil {
		return nil
	}
	pgErr := new(pgconn.PgError)
	// check events2 not exists
	if errors.As(err, &pgErr) && pgErr.Code == "42P01" {
		return query(ctx, crdb, searchQuery, reduce, true)
	}
	return err
}

// LatestSequence returns the latest sequence found by the search query
func (db *CRDB) LatestSequence(ctx context.Context, searchQuery *eventstore.SearchQueryBuilder) (float64, error) {
	var position sql.NullFloat64
	err := query(ctx, db, searchQuery, &position, false)
	return position.Float64, err
}

// InstanceIDs returns the instance ids found by the search query
func (db *CRDB) InstanceIDs(ctx context.Context, searchQuery *eventstore.SearchQueryBuilder) ([]string, error) {
	var ids []string
	err := query(ctx, db, searchQuery, &ids, false)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (db *CRDB) Client() *database.DB {
	return db.DB
}

func (db *CRDB) orderByEventSequence(desc, shouldOrderBySequence, useV1 bool) string {
	if useV1 {
		if desc {
			return ` ORDER BY event_sequence DESC`
		}
		return ` ORDER BY event_sequence`
	}
	if shouldOrderBySequence {
		if desc {
			return ` ORDER BY "sequence" DESC`
		}
		return ` ORDER BY "sequence"`
	}

	if desc {
		return ` ORDER BY "position" DESC, in_tx_order DESC`
	}
	return ` ORDER BY "position", in_tx_order`
}

func (db *CRDB) eventQuery(useV1 bool) string {
	if useV1 {
		return "SELECT" +
			" creation_date" +
			", event_type" +
			", event_sequence" +
			", event_data" +
			", editor_user" +
			", resource_owner" +
			", instance_id" +
			", aggregate_type" +
			", aggregate_id" +
			", aggregate_version" +
			" FROM eventstore.events"
	}
	return "SELECT" +
		" created_at" +
		", event_type" +
		`, "sequence"` +
		`, "position"` +
		", payload" +
		", creator" +
		`, "owner"` +
		", instance_id" +
		", aggregate_type" +
		", aggregate_id" +
		", revision" +
		" FROM eventstore.events2"
}

func (db *CRDB) maxSequenceQuery(useV1 bool) string {
	if useV1 {
		return `SELECT event_sequence FROM eventstore.events`
	}
	return `SELECT "position" FROM eventstore.events2`
}

func (db *CRDB) instanceIDsQuery(useV1 bool) string {
	table := "eventstore.events2"
	if useV1 {
		table = "eventstore.events"
	}
	return "SELECT DISTINCT instance_id FROM " + table
}

func (db *CRDB) columnName(col repository.Field, useV1 bool) string {
	switch col {
	case repository.FieldAggregateID:
		return "aggregate_id"
	case repository.FieldAggregateType:
		return "aggregate_type"
	case repository.FieldSequence:
		if useV1 {
			return "event_sequence"
		}
		return `"sequence"`
	case repository.FieldResourceOwner:
		if useV1 {
			return "resource_owner"
		}
		return `"owner"`
	case repository.FieldInstanceID:
		return "instance_id"
	case repository.FieldEditorService:
		if useV1 {
			return "editor_service"
		}
		return ""
	case repository.FieldEditorUser:
		if useV1 {
			return "editor_user"
		}
		return "creator"
	case repository.FieldEventType:
		return "event_type"
	case repository.FieldEventData:
		if useV1 {
			return "event_data"
		}
		return "payload"
	case repository.FieldCreationDate:
		if useV1 {
			return "creation_date"
		}
		return "created_at"
	case repository.FieldPosition:
		return `"position"`
	default:
		return ""
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
	if pgxErr, ok := err.(*pgconn.PgError); ok {
		if pgxErr.Code == "23505" {
			return true
		}
	}
	return false
}
