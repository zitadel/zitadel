package query

import (
	"context"
	"database/sql"
	_ "embed"
	errs "errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type LatestState struct {
	EventTimestamp time.Time
	Position       uint64
	LastUpdated    time.Time
}

type CurrentStates struct {
	SearchResponse
	CurrentStates []*CurrentState
}

type CurrentState struct {
	ProjectionName    string
	CurrentPosition   uint64
	EventCreationDate time.Time
	LastRun           time.Time
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

	query, scan := prepareCurrentStateQuery(ctx, q.client)
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-MmFef", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		currentStates, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-22H8f", "Errors.Internal")
	}

	return currentStates, nil
}

func (q *Queries) latestState(ctx context.Context, projections ...table) (state *LatestState, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareLatestState(ctx, q.client)
	or := make(sq.Or, len(projections))
	for i, projection := range projections {
		or[i] = sq.Eq{CurrentStateColProjectionName.identifier(): projection.name}
	}
	stmt, args, err := query.
		Where(or).
		Where(sq.Eq{CurrentStateColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}).
		OrderBy(CurrentStateColEventDate.identifier() + " DESC").
		ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-5CfX9", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		state, err = scan(row)
		return err
	}, stmt, args...)

	return state, err
}

func (q *Queries) ClearCurrentSequence(ctx context.Context, projectionName string) (err error) {
	tx, err := q.client.Begin()
	if err != nil {
		return errors.ThrowInternal(err, "QUERY-9iOpr", "Errors.RemoveFailed")
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			logging.OnError(rollbackErr).Debug("rollback failed")
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			err = errors.ThrowInternal(commitErr, "QUERY-JGD0l", "Errors.Internal")
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
		return errors.ThrowInternal(err, "QUERY-Sfvsc", "Errors.Internal")
	}
	return nil
}

func (q *Queries) checkAndLock(tx *sql.Tx, projectionName string) (name string, err error) {
	stmt, args, err := sq.Select(CurrentStateColProjectionName.identifier()).
		From(currentStateTable.identifier()).
		Where(sq.Eq{
			CurrentStateColProjectionName.identifier(): projectionName,
		}).Suffix("FOR UPDATE").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", errors.ThrowInternal(err, "QUERY-UJTUy", "Errors.Internal")
	}
	row := tx.QueryRow(stmt, args...)
	if err := row.Scan(&name); err != nil || name == "" {
		return "", errors.ThrowInternal(err, "QUERY-ej8fn", "Errors.ProjectionName.Invalid")
	}
	return name, nil
}

func tablesForReset(ctx context.Context, tx *sql.Tx, projectionName string) (tables []string, err error) {
	names := strings.Split(projectionName, ".")
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
		return nil, errors.ThrowInternal(err, "QUERY-ASff2", "Errors.ProjectionName.Invalid")
	}

	rows, err := tx.QueryContext(ctx, tablesQuery, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Dgfw", "Errors.ProjectionName.Invalid")
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, errors.ThrowInternal(err, "QUERY-ej8fn", "Errors.ProjectionName.Invalid")
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
			return errors.ThrowInternal(err, "QUERY-3n92f", "Errors.RemoveFailed")
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
		return errors.ThrowInternal(err, "QUERY-Ff3tw", "Errors.RemoveFailed")
	}
	_, err = tx.Exec(update, args...)
	if err != nil {
		return errors.ThrowInternal(err, "QUERY-NFiws", "Errors.RemoveFailed")
	}
	return nil
}

func prepareLatestState(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*LatestState, error)) {
	return sq.Select(
			CurrentStateColEventDate.identifier(),
			CurrentStateColPosition.identifier(),
			CurrentStateColLastUpdated.identifier()).
			From(currentStateTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*LatestState, error) {
			var (
				creationDate sql.NullTime
				lastUpdated  sql.NullTime
				position     sql.NullInt64
			)
			err := row.Scan(
				&creationDate,
				&position,
				&lastUpdated,
			)
			if err != nil && !errs.Is(err, sql.ErrNoRows) {
				return nil, errors.ThrowInternal(err, "QUERY-aAZ1D", "Errors.Internal")
			}
			return &LatestState{
				EventTimestamp: creationDate.Time,
				LastUpdated:    lastUpdated.Time,
				Position:       uint64(position.Int64),
			}, nil
		}
}

func prepareCurrentStateQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*CurrentStates, error)) {
	return sq.Select(
			CurrentStateColLastUpdated.identifier(),
			CurrentStateColEventDate.identifier(),
			CurrentStateColPosition.identifier(),
			CurrentStateColProjectionName.identifier(),
			countColumn.identifier()).
			From(currentStateTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*CurrentStates, error) {
			states := make([]*CurrentState, 0)
			var count uint64
			for rows.Next() {
				currentState := new(CurrentState)
				var lastRun sql.NullTime
				var eventDate sql.NullTime
				var currentPosition sql.NullInt64

				err := rows.Scan(
					&lastRun,
					&eventDate,
					&currentPosition,
					&currentState.ProjectionName,
					&count,
				)
				if err != nil {
					return nil, err
				}
				currentState.EventCreationDate = eventDate.Time
				currentState.LastRun = lastRun.Time
				currentState.CurrentPosition = uint64(currentPosition.Int64)
				states = append(states, currentState)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-jbJ77", "Errors.Query.CloseRows")
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
