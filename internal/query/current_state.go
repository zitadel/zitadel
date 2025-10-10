package query

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/shopspring/decimal"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Stateful interface {
	SetState(*State)
}

type State struct {
	LastRun time.Time

	Position       decimal.Decimal
	EventCreatedAt time.Time
	AggregateID    string
	AggregateType  eventstore.AggregateType
	Sequence       uint64
}

type CurrentStates struct {
	SearchResponse
	CurrentStates []*CurrentState
}

type CurrentState struct {
	ProjectionName string
	State
}

type CurrentStateSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func NewCurrentStatesInstanceIDSearchQuery(instanceID string) (SearchQuery, error) {
	return NewTextQuery(CurrentStateColInstanceID, instanceID, TextEquals)
}

func NewCurrentStatesProjectionSearchQuery(projection string) (SearchQuery, error) {
	return NewTextQuery(CurrentStateColProjectionName, projection, TextEquals)
}

func (q *CurrentStateSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) SearchCurrentStates(ctx context.Context, queries *CurrentStateSearchQueries) (currentStates *CurrentStates, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareCurrentStateQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-MmFef", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		currentStates, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-22H8f", "Errors.Internal")
	}

	return currentStates, nil
}

func (q *Queries) latestState(ctx context.Context, projections ...table) (state *State, err error) {
	return latestState(ctx, q.client, projections...)
}

func (q *Queries) ClearCurrentSequence(ctx context.Context, projectionName string) (err error) {
	tx, err := q.client.BeginTx(ctx, nil)
	if err != nil {
		return zerrors.ThrowInternal(err, "QUERY-9iOpr", "Errors.RemoveFailed")
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			logging.OnError(rollbackErr).Debug("rollback failed")
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = zerrors.ThrowInternal(commitErr, "QUERY-JGD0l", "Errors.Internal")
		}
	}()

	name, err := q.checkAndLock(tx, projectionName)
	if err != nil {
		return err
	}

	tables, err := tablesForReset(ctx, tx, name)
	if err != nil {
		return err
	}
	err = reset(ctx, tx, tables, name)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return zerrors.ThrowInternal(err, "QUERY-Sfvsc", "Errors.Internal")
	}
	return nil
}

func (q *Queries) checkAndLock(tx *sql.Tx, projectionName string) (name string, err error) {
	stmt, args, err := sq.Select(CurrentStateColProjectionName.identifier()).
		From(currentStateTable.identifier()).
		Where(sq.Eq{
			CurrentStateColProjectionName.identifier(): projectionName,
		}).Suffix("FOR NO KEY UPDATE").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", zerrors.ThrowInternal(err, "QUERY-UJTUy", "Errors.Internal")
	}
	row := tx.QueryRow(stmt, args...)
	if err := row.Scan(&name); err != nil || name == "" {
		return "", zerrors.ThrowInternal(err, "QUERY-ej8fn", "Errors.ProjectionName.Invalid")
	}
	return name, nil
}

func tablesForReset(ctx context.Context, tx *sql.Tx, projectionName string) (tables []string, err error) {
	names := strings.Split(projectionName, ".")
	if len(names) != 2 {
		return nil, zerrors.ThrowInvalidArgument(nil, "QUERY-wk1jr", "Errors.InvalidArgument")
	}
	schema := names[0]
	tablePrefix := names[1]

	tablesQuery, args, err := sq.Select("table_name").
		From("[show tables from " + schema + "]").
		Where(
			sq.And{
				sq.Eq{"type": "table"},
				sq.NotEq{"table_name": []string{"locks", "current_sequences", "current_states", "failed_events", "failed_events2"}},
				sq.Like{"table_name": tablePrefix + "%"},
			}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-ASff2", "Errors.ProjectionName.Invalid")
	}

	rows, err := tx.QueryContext(ctx, tablesQuery, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Dgfw", "Errors.ProjectionName.Invalid")
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, zerrors.ThrowInternal(err, "QUERY-ej8fn", "Errors.ProjectionName.Invalid")
		}
		tables = append(tables, schema+"."+tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func reset(ctx context.Context, tx *sql.Tx, tables []string, projectionName string) error {
	for _, tableName := range tables {
		_, err := tx.Exec(fmt.Sprintf("TRUNCATE %s cascade", tableName))
		if err != nil {
			return zerrors.ThrowInternal(err, "QUERY-3n92f", "Errors.RemoveFailed")
		}
	}
	update, args, err := sq.Update(currentStateTable.identifier()).
		Set(CurrentStateColEventDate.name, 0).
		Where(sq.Eq{
			CurrentStateColProjectionName.name: projectionName,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "QUERY-Ff3tw", "Errors.RemoveFailed")
	}
	_, err = tx.ExecContext(ctx, update, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "QUERY-NFiws", "Errors.RemoveFailed")
	}
	return nil
}

func prepareLatestState() (sq.SelectBuilder, func(*sql.Row) (*State, error)) {
	return sq.Select(
			CurrentStateColEventDate.identifier(),
			CurrentStateColPosition.identifier(),
			CurrentStateColLastUpdated.identifier()).
			From(currentStateTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*State, error) {
			var (
				creationDate sql.NullTime
				lastUpdated  sql.NullTime
				position     decimal.NullDecimal
			)
			err := row.Scan(
				&creationDate,
				&position,
				&lastUpdated,
			)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return nil, zerrors.ThrowInternal(err, "QUERY-aAZ1D", "Errors.Internal")
			}
			return &State{
				EventCreatedAt: creationDate.Time,
				LastRun:        lastUpdated.Time,
				Position:       position.Decimal,
			}, nil
		}
}

func prepareCurrentStateQuery() (sq.SelectBuilder, func(*sql.Rows) (*CurrentStates, error)) {
	return sq.Select(
			CurrentStateColLastUpdated.identifier(),
			CurrentStateColEventDate.identifier(),
			CurrentStateColPosition.identifier(),
			CurrentStateColProjectionName.identifier(),
			CurrentStateColAggregateType.identifier(),
			CurrentStateColAggregateID.identifier(),
			CurrentStateColSequence.identifier(),
			countColumn.identifier()).
			From(currentStateTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*CurrentStates, error) {
			states := make([]*CurrentState, 0)
			var count uint64
			for rows.Next() {
				currentState := new(CurrentState)
				var (
					lastRun         sql.NullTime
					eventDate       sql.NullTime
					currentPosition decimal.NullDecimal
					aggregateType   sql.NullString
					aggregateID     sql.NullString
					sequence        sql.NullInt64
				)

				err := rows.Scan(
					&lastRun,
					&eventDate,
					&currentPosition,
					&currentState.ProjectionName,
					&aggregateType,
					&aggregateID,
					&sequence,
					&count,
				)
				if err != nil {
					return nil, err
				}
				currentState.State.EventCreatedAt = eventDate.Time
				currentState.State.LastRun = lastRun.Time
				currentState.Position = currentPosition.Decimal
				currentState.AggregateType = eventstore.AggregateType(aggregateType.String)
				currentState.AggregateID = aggregateID.String
				currentState.Sequence = uint64(sequence.Int64)
				states = append(states, currentState)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-jbJ77", "Errors.Query.CloseRows")
			}

			return &CurrentStates{
				CurrentStates: states,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

var (
	currentStateTable = table{
		name:          projection.CurrentStateTable,
		instanceIDCol: "instance_id",
	}
	CurrentStateColEventDate = Column{
		name:  "event_date",
		table: currentStateTable,
	}
	CurrentStateColPosition = Column{
		name:  "position",
		table: currentStateTable,
	}
	CurrentStateColAggregateType = Column{
		name:  "aggregate_type",
		table: currentStateTable,
	}
	CurrentStateColAggregateID = Column{
		name:  "aggregate_id",
		table: currentStateTable,
	}
	CurrentStateColSequence = Column{
		name:  "sequence",
		table: currentStateTable,
	}
	CurrentStateColLastUpdated = Column{
		name:  "last_updated",
		table: currentStateTable,
	}
	CurrentStateColProjectionName = Column{
		name:  "projection_name",
		table: currentStateTable,
	}
	CurrentStateColInstanceID = Column{
		name:  "instance_id",
		table: currentStateTable,
	}
)

var (
	locksTable = table{
		name:          projection.LocksTable,
		instanceIDCol: "instance_id",
	}
	LocksColLockerID = Column{
		name:  "locker_id",
		table: locksTable,
	}
	LocksColUntil = Column{
		name:  "locked_until",
		table: locksTable,
	}
	LocksColProjectionName = Column{
		name:  "projection_name",
		table: locksTable,
	}
)
