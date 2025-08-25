package sql

import (
	"context"
	"errors"
	"regexp"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

// awaitOpenTransactions ensures event ordering, so we don't events younger that open transactions
var (
	awaitOpenTransactionsV1 = ` AND created_at <= now()`
	awaitOpenTransactionsV2 = ` AND "position" <= EXTRACT(EPOCH FROM now())`
)

func awaitOpenTransactions(useV1 bool) string {
	if useV1 {
		return awaitOpenTransactionsV1
	}
	return awaitOpenTransactionsV2
}

type Postgres struct {
	*database.DB
}

func NewPostgres(client *database.DB) *Postgres {
	return &Postgres{client}
}

func (db *Postgres) Health(ctx context.Context) error { return db.Ping() }

// FilterToReducer finds all events matching the given search query and passes them to the reduce function.
func (psql *Postgres) FilterToReducer(ctx context.Context, searchQuery *eventstore.SearchQueryBuilder, reduce eventstore.Reducer) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = query(ctx, psql, searchQuery, reduce, false)
	if err == nil {
		return nil
	}
	pgErr := new(pgconn.PgError)
	// check events2 not exists
	if errors.As(err, &pgErr) && pgErr.Code == "42P01" {
		return query(ctx, psql, searchQuery, reduce, true)
	}
	return err
}

// LatestPosition returns the latest position found by the search query
func (db *Postgres) LatestPosition(ctx context.Context, searchQuery *eventstore.SearchQueryBuilder) (decimal.Decimal, error) {
	var position decimal.Decimal
	err := query(ctx, db, searchQuery, &position, false)
	return position, err
}

// InstanceIDs returns the instance ids found by the search query
func (db *Postgres) InstanceIDs(ctx context.Context, searchQuery *eventstore.SearchQueryBuilder) ([]string, error) {
	var ids []string
	err := query(ctx, db, searchQuery, &ids, false)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (db *Postgres) Client() *database.DB {
	return db.DB
}

func (db *Postgres) orderByEventSequence(desc, shouldOrderBySequence, useV1 bool) string {
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

func (db *Postgres) eventQuery(useV1 bool) string {
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

func (db *Postgres) maxPositionQuery(useV1 bool) string {
	if useV1 {
		return `SELECT event_sequence FROM eventstore.events`
	}
	return `SELECT "position" FROM eventstore.events2`
}

func (db *Postgres) instanceIDsQuery(useV1 bool) string {
	table := "eventstore.events2"
	if useV1 {
		table = "eventstore.events"
	}
	return "SELECT DISTINCT instance_id FROM " + table
}

func (db *Postgres) columnName(col repository.Field, useV1 bool) string {
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

func (db *Postgres) conditionFormat(operation repository.Operation) string {
	switch operation {
	case repository.OperationIn:
		return "%s %s ANY(?)"
	case repository.OperationNotIn:
		return "%s %s ALL(?)"
	case repository.OperationEquals, repository.OperationGreater, repository.OperationLess, repository.OperationJSONContains:
		fallthrough
	default:
		return "%s %s ?"
	}
}

func (db *Postgres) operation(operation repository.Operation) string {
	switch operation {
	case repository.OperationEquals, repository.OperationIn:
		return "="
	case repository.OperationGreater:
		return ">"
	case repository.OperationGreaterOrEquals:
		return ">="
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
func (db *Postgres) placeholder(query string) string {
	occurrences := placeholder.FindAllStringIndex(query, -1)
	if len(occurrences) == 0 {
		return query
	}
	replaced := query[:occurrences[0][0]]

	for i, l := range occurrences {
		nextIDX := len(query)
		if i < len(occurrences)-1 {
			nextIDX = occurrences[i+1][0]
		}
		replaced = replaced + "$" + strconv.Itoa(i+1) + query[l[1]:nextIDX]
	}
	return replaced
}
