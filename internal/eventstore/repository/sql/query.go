package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/dialect"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type querier interface {
	columnName(field repository.Field, useV1 bool) string
	operation(repository.Operation) string
	conditionFormat(repository.Operation) string
	placeholder(query string) string
	eventQuery(useV1 bool) string
	maxSequenceQuery(useV1 bool) string
	instanceIDsQuery(useV1 bool) string
	Client() *database.DB
	orderByEventSequence(desc, shouldOrderBySequence, useV1 bool) string
	dialect.Database
}

type scan func(dest ...interface{}) error

type tx struct {
	*sql.Tx
}

func (t *tx) QueryContext(ctx context.Context, scan func(rows *sql.Rows) error, query string, args ...any) error {
	rows, err := t.Tx.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := rows.Close()
		logging.OnError(closeErr).Info("rows.Close failed")
	}()

	if err = scan(rows); err != nil {
		return err
	}
	return rows.Err()
}

func query(ctx context.Context, criteria querier, searchQuery *eventstore.SearchQueryBuilder, dest interface{}, useV1 bool) error {
	q, err := repository.QueryFromBuilder(searchQuery)
	if err != nil {
		return err
	}

	query, rowScanner := prepareColumns(criteria, q.Columns, useV1)
	where, values := prepareConditions(criteria, q, useV1)
	if where == "" || query == "" {
		return zerrors.ThrowInvalidArgument(nil, "SQL-rWeBw", "invalid query factory")
	}
	if q.Tx == nil {
		if travel := prepareTimeTravel(ctx, criteria, q.AllowTimeTravel); travel != "" {
			query += travel
		}
	}
	query += where

	// instead of using the max function of the database (which doesn't work for postgres)
	// we select the most recent row
	if q.Columns == eventstore.ColumnsMaxSequence {
		q.Limit = 1
		q.Desc = true
	}

	// if there is only one subquery we can optimize the query ordering by ordering by sequence
	var shouldOrderBySequence bool
	if len(q.SubQueries) == 1 {
		for _, filter := range q.SubQueries[0] {
			if filter.Field == repository.FieldAggregateID {
				shouldOrderBySequence = filter.Operation == repository.OperationEquals
			}
		}
	}

	switch q.Columns {
	case eventstore.ColumnsEvent,
		eventstore.ColumnsMaxSequence:
		query += criteria.orderByEventSequence(q.Desc, shouldOrderBySequence, useV1)
	}

	if q.Limit > 0 {
		values = append(values, q.Limit)
		query += " LIMIT ?"
	}

	if q.Offset > 0 {
		values = append(values, q.Offset)
		query += " OFFSET ?"
	}

	if q.LockRows {
		query += " FOR UPDATE"
		switch q.LockOption {
		case eventstore.LockOptionWait: // default behavior
		case eventstore.LockOptionNoWait:
			query += " NOWAIT"
		case eventstore.LockOptionSkipLocked:
			query += " SKIP LOCKED"

		}
	}

	query = criteria.placeholder(query)

	var contextQuerier interface {
		QueryContext(context.Context, func(rows *sql.Rows) error, string, ...interface{}) error
	}
	contextQuerier = criteria.Client()
	if q.Tx != nil {
		contextQuerier = &tx{Tx: q.Tx}
	}

	err = contextQuerier.QueryContext(ctx,
		func(rows *sql.Rows) error {
			for rows.Next() {
				err := rowScanner(rows.Scan, dest)
				if err != nil {
					return err
				}
			}
			return nil
		}, query, values...)
	if err != nil {
		logging.New().WithError(err).Info("query failed")
		return zerrors.ThrowInternal(err, "SQL-KyeAx", "unable to filter events")
	}

	return nil
}

func prepareColumns(criteria querier, columns eventstore.Columns, useV1 bool) (string, func(s scan, dest interface{}) error) {
	switch columns {
	case eventstore.ColumnsMaxSequence:
		return criteria.maxSequenceQuery(useV1), maxSequenceScanner
	case eventstore.ColumnsInstanceIDs:
		return criteria.instanceIDsQuery(useV1), instanceIDsScanner
	case eventstore.ColumnsEvent:
		return criteria.eventQuery(useV1), eventsScanner(useV1)
	default:
		return "", nil
	}
}

func prepareTimeTravel(ctx context.Context, criteria querier, allow bool) string {
	if !allow {
		return ""
	}
	took := call.Took(ctx)
	return criteria.Timetravel(took)
}

func maxSequenceScanner(row scan, dest interface{}) (err error) {
	position, ok := dest.(*sql.NullFloat64)
	if !ok {
		return zerrors.ThrowInvalidArgumentf(nil, "SQL-NBjA9", "type must be sql.NullInt64 got: %T", dest)
	}
	err = row(position)
	if err == nil || errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return zerrors.ThrowInternal(err, "SQL-bN5xg", "something went wrong")
}

func instanceIDsScanner(scanner scan, dest interface{}) (err error) {
	ids, ok := dest.(*[]string)
	if !ok {
		return zerrors.ThrowInvalidArgument(nil, "SQL-Begh2", "type must be an array of string")
	}
	var id string
	err = scanner(&id)
	if err != nil {
		logging.WithError(err).Warn("unable to scan row")
		return zerrors.ThrowInternal(err, "SQL-DEFGe", "unable to scan row")
	}
	*ids = append(*ids, id)

	return nil
}

func eventsScanner(useV1 bool) func(scanner scan, dest interface{}) (err error) {
	return func(scanner scan, dest interface{}) (err error) {
		reduce, ok := dest.(eventstore.Reducer)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "SQL-4GP6F", "events scanner: invalid type %T", dest)
		}
		event := new(repository.Event)
		position := new(sql.NullFloat64)

		if useV1 {
			err = scanner(
				&event.CreationDate,
				&event.Typ,
				&event.Seq,
				&event.Data,
				&event.EditorUser,
				&event.ResourceOwner,
				&event.InstanceID,
				&event.AggregateType,
				&event.AggregateID,
				&event.Version,
			)
		} else {
			var revision uint8
			err = scanner(
				&event.CreationDate,
				&event.Typ,
				&event.Seq,
				position,
				&event.Data,
				&event.EditorUser,
				&event.ResourceOwner,
				&event.InstanceID,
				&event.AggregateType,
				&event.AggregateID,
				&revision,
			)
			event.Version = eventstore.Version("v" + strconv.Itoa(int(revision)))
		}

		if err != nil {
			logging.New().WithError(err).Warn("unable to scan row")
			return zerrors.ThrowInternal(err, "SQL-M0dsf", "unable to scan row")
		}
		event.Pos = position.Float64
		return reduce(event)
	}
}

func prepareConditions(criteria querier, query *repository.SearchQuery, useV1 bool) (_ string, args []any) {
	clauses, args := prepareQuery(criteria, useV1, query.InstanceID, query.InstanceIDs, query.ExcludedInstances)
	if clauses != "" && len(query.SubQueries) > 0 {
		clauses += " AND "
	}
	subClauses := make([]string, len(query.SubQueries))
	for i, filters := range query.SubQueries {
		var subArgs []any
		subClauses[i], subArgs = prepareQuery(criteria, useV1, filters...)
		// an error is thrown in [query]
		if subClauses[i] == "" {
			return "", nil
		}
		if len(query.SubQueries) > 1 && len(subArgs) > 1 {
			subClauses[i] = "(" + subClauses[i] + ")"
		}
		args = append(args, subArgs...)
	}
	if len(subClauses) == 1 {
		clauses += subClauses[0]
	} else if len(subClauses) > 1 {
		clauses += "(" + strings.Join(subClauses, " OR ") + ")"
	}

	additionalClauses, additionalArgs := prepareQuery(criteria, useV1,
		query.Position,
		query.Owner,
		query.Sequence,
		query.CreatedAfter,
		query.CreatedBefore,
		query.Creator,
	)
	if additionalClauses != "" {
		if clauses != "" {
			clauses += " AND "
		}
		clauses += additionalClauses
		args = append(args, additionalArgs...)
	}

	excludeAggregateIDs := query.ExcludeAggregateIDs
	if len(excludeAggregateIDs) > 0 {
		excludeAggregateIDs = append(excludeAggregateIDs, query.InstanceID, query.InstanceIDs, query.Position, query.CreatedAfter, query.CreatedBefore)
	}
	excludeAggregateIDsClauses, excludeAggregateIDsArgs := prepareQuery(criteria, useV1, excludeAggregateIDs...)
	if excludeAggregateIDsClauses != "" {
		if clauses != "" {
			clauses += " AND "
		}
		if useV1 {
			clauses += "aggregate_id NOT IN (SELECT aggregate_id FROM eventstore.events WHERE " + excludeAggregateIDsClauses + ")"
		} else {
			clauses += "aggregate_id NOT IN (SELECT aggregate_id FROM eventstore.events2 WHERE " + excludeAggregateIDsClauses + ")"
		}
		args = append(args, excludeAggregateIDsArgs...)
	}

	if query.AwaitOpenTransactions {
		instanceIDs := make(database.TextArray[string], 0, 3)
		if query.InstanceID != nil {
			instanceIDs = append(instanceIDs, query.InstanceID.Value.(string))
		} else if query.InstanceIDs != nil {
			instanceIDs = append(instanceIDs, query.InstanceIDs.Value.(database.TextArray[string])...)
		}

		for i := range instanceIDs {
			instanceIDs[i] = dialect.DBPurposeEventPusher.AppName() + "_" + instanceIDs[i]
		}

		clauses += awaitOpenTransactions(useV1)
		args = append(args, instanceIDs)
	}

	if clauses == "" {
		return "", nil
	}

	return " WHERE " + clauses, args
}

func prepareQuery(criteria querier, useV1 bool, filters ...*repository.Filter) (_ string, args []any) {
	clauses := make([]string, 0, len(filters))
	args = make([]any, 0, len(filters))
	for _, filter := range filters {
		if filter == nil {
			continue
		}
		arg := filter.Value

		// marshal if payload filter
		if filter.Field == repository.FieldEventData {
			var err error
			arg, err = json.Marshal(arg)
			if err != nil {
				logging.WithError(err).Warn("unable to marshal search value")
				continue
			}

		}

		clauses = append(clauses, getCondition(criteria, filter, useV1))
		// if mapping failed an error is thrown in [query]
		if clauses[len(clauses)-1] == "" {
			return "", nil
		}
		args = append(args, arg)
	}

	return strings.Join(clauses, " AND "), args
}

func getCondition(cond querier, filter *repository.Filter, useV1 bool) (condition string) {
	field := cond.columnName(filter.Field, useV1)
	operation := cond.operation(filter.Operation)
	if field == "" || operation == "" {
		return ""
	}
	format := cond.conditionFormat(filter.Operation)

	return fmt.Sprintf(format, field, operation)
}
