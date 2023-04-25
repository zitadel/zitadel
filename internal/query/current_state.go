package query

import (
	"context"
	"database/sql"
	errs "errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

const (
	// TODO: use current_states instead
	lockStmtFormat = "UPDATE %[1]s" +
		" set locker_id = $1, locked_until = now()+$2::INTERVAL" +
		" WHERE projection_name = $3"
	lockerIDReset = "reset"
)

type LatestState struct {
	EventTimestamp time.Time
	Sequence       uint64
	LastUpdated    time.Time
}

type CurrentStates struct {
	SearchResponse
	CurrentStates []*CurrentState
}

type CurrentState struct {
	ProjectionName    string
	AggregateType     string
	AggregateID       string
	EventSequence     uint64
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

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-22H8f", "Errors.Internal")
	}
	return scan(rows)
}

func (q *Queries) latestState(ctx context.Context, projections ...table) (_ *LatestState, err error) {
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

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func (q *Queries) ClearCurrentSequence(ctx context.Context, projectionName string) (err error) {
	err = q.checkAndLock(ctx, projectionName)
	if err != nil {
		return err
	}

	tx, err := q.client.Begin()
	if err != nil {
		return errors.ThrowInternal(err, "QUERY-9iOpr", "Errors.RemoveFailed")
	}
	tables, err := tablesForReset(ctx, tx, projectionName)
	if err != nil {
		return err
	}
	err = reset(tx, tables, projectionName)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (q *Queries) checkAndLock(ctx context.Context, projectionName string) error {
	projectionQuery, args, err := sq.Select("count(*)").
		From("[show tables from projections]").
		Where(
			sq.And{
				sq.NotEq{"table_name": []string{"locks", "current_sequences", "failed_events"}},
				sq.Eq{"concat('projections.', table_name)": projectionName},
			}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return errors.ThrowInternal(err, "QUERY-Dfwf2", "Errors.ProjectionName.Invalid")
	}
	row := q.client.QueryRowContext(ctx, projectionQuery, args...)
	var count int
	if err := row.Scan(&count); err != nil || count == 0 {
		return errors.ThrowInternal(err, "QUERY-ej8fn", "Errors.ProjectionName.Invalid")
	}
	lock := fmt.Sprintf(lockStmtFormat, locksTable.identifier())
	if err != nil {
		return errors.ThrowInternal(err, "QUERY-DVfg3", "Errors.RemoveFailed")
	}
	//lock for twice the default duration (10s)
	res, err := q.client.ExecContext(ctx, lock, lockerIDReset, 20*time.Second, projectionName)
	if err != nil {
		return errors.ThrowInternal(err, "QUERY-WEfr2", "Errors.RemoveFailed")
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return errors.ThrowInternal(err, "QUERY-Bh3ws", "Errors.RemoveFailed")
	}
	time.Sleep(7 * time.Second) //more than half the default lock duration (10s)
	return nil
}

func tablesForReset(ctx context.Context, tx *sql.Tx, projectionName string) ([]string, error) {
	tablesQuery, args, err := sq.Select("concat('projections.', table_name)").
		From("[show tables from projections]").
		Where(
			sq.And{
				sq.Eq{"type": "table"},
				sq.NotEq{"table_name": []string{"locks", "current_sequences", "failed_events"}},
				sq.Like{"concat('projections.', table_name)": projectionName + "%"},
			}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-ASff2", "Errors.ProjectionName.Invalid")
	}
	var tables []string
	rows, err := tx.QueryContext(ctx, tablesQuery, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Dgfw", "Errors.ProjectionName.Invalid")
	}
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, errors.ThrowInternal(err, "QUERY-ej8fn", "Errors.ProjectionName.Invalid")
		}
		tables = append(tables, tableName)
	}
	return tables, nil
}

func reset(tx *sql.Tx, tables []string, projectionName string) error {
	for _, tableName := range tables {
		_, err := tx.Exec(fmt.Sprintf("TRUNCATE %s cascade", tableName))
		if err != nil {
			return errors.ThrowInternal(err, "QUERY-3n92f", "Errors.RemoveFailed")
		}
	}
	update, args, err := sq.Update(currentStateTable.identifier()).
		Set(CurrentStateColEventDate.name, 0).
		Where(sq.Eq{CurrentStateColProjectionName.name: projectionName}).
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
			CurrentStateColSequence.identifier(),
			CurrentStateColLastUpdated.identifier()).
			From(currentStateTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*LatestState, error) {
			state := new(LatestState)
			err := row.Scan(
				&state.EventTimestamp,
				&state.Sequence,
				&state.LastUpdated,
			)
			if err != nil && !errs.Is(err, sql.ErrNoRows) {
				return nil, errors.ThrowInternal(err, "QUERY-aAZ1D", "Errors.Internal")
			}
			return state, nil
		}
}

func prepareCurrentStateQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*CurrentStates, error)) {
	return sq.Select(
			CurrentStateColLastUpdated.identifier(),
			CurrentStateColLastAggregateType.identifier(),
			CurrentStateColLastAggregateID.identifier(),
			CurrentStateColEventDate.identifier(),
			CurrentStateColSequence.identifier(),
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
				var currentSequence sql.NullInt64
				var aggregateType sql.NullString
				var aggregateID sql.NullString

				err := rows.Scan(
					&lastRun,
					&aggregateType,
					&aggregateID,
					&eventDate,
					&currentSequence,
					&currentState.ProjectionName,
					&count,
				)
				if err != nil {
					return nil, err
				}
				currentState.EventCreationDate = eventDate.Time
				currentState.LastRun = lastRun.Time
				currentState.EventSequence = uint64(currentSequence.Int64)
				currentState.AggregateType = aggregateType.String
				currentState.AggregateID = aggregateID.String
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
	CurrentStateColSequence = Column{
		name:  "event_sequence",
		table: currentStateTable,
	}
	CurrentStateColLastUpdated = Column{
		name:  "last_updated",
		table: currentStateTable,
	}
	CurrentStateColLastAggregateType = Column{
		name:  "aggregate_type",
		table: currentStateTable,
	}
	CurrentStateColLastAggregateID = Column{
		name:  "aggregate_id",
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
