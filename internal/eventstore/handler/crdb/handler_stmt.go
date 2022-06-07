package crdb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var (
	errSeqNotUpdated = errors.ThrowInternal(nil, "CRDB-79GWt", "current sequence not updated")
)

type StatementHandlerConfig struct {
	handler.ProjectionHandlerConfig

	Client            *sql.DB
	SequenceTable     string
	LockTable         string
	FailedEventsTable string
	MaxFailureCount   uint
	BulkLimit         uint64

	Reducers []handler.AggregateReducer
}

type StatementHandler struct {
	*handler.ProjectionHandler
	Locker

	client                  *sql.DB
	sequenceTable           string
	currentSequenceStmt     string
	updateSequencesBaseStmt string
	maxFailureCount         uint
	failureCountStmt        string
	setFailureCountStmt     string

	aggregates []eventstore.AggregateType
	reduces    map[eventstore.EventType]handler.Reduce

	bulkLimit uint64
}

func NewStatementHandler(
	ctx context.Context,
	config StatementHandlerConfig,
) StatementHandler {
	aggregateTypes := make([]eventstore.AggregateType, 0, len(config.Reducers))
	reduces := make(map[eventstore.EventType]handler.Reduce, len(config.Reducers))
	for _, aggReducer := range config.Reducers {
		aggregateTypes = append(aggregateTypes, aggReducer.Aggregate)
		for _, eventReducer := range aggReducer.EventRedusers {
			reduces[eventReducer.Event] = eventReducer.Reduce
		}
	}

	h := StatementHandler{
		client:                  config.Client,
		sequenceTable:           config.SequenceTable,
		maxFailureCount:         config.MaxFailureCount,
		currentSequenceStmt:     fmt.Sprintf(currentSequenceStmtFormat, config.SequenceTable),
		updateSequencesBaseStmt: fmt.Sprintf(updateCurrentSequencesStmtFormat, config.SequenceTable),
		failureCountStmt:        fmt.Sprintf(failureCountStmtFormat, config.FailedEventsTable),
		setFailureCountStmt:     fmt.Sprintf(setFailureCountStmtFormat, config.FailedEventsTable),
		aggregates:              aggregateTypes,
		reduces:                 reduces,
		bulkLimit:               config.BulkLimit,
		Locker:                  NewLocker(config.Client, config.LockTable, config.ProjectionName),
	}

	h.ProjectionHandler = handler.NewProjectionHandler(config.ProjectionHandlerConfig, h.reduce, h.Update, h.SearchQuery)

	go h.Process(
		ctx,
		h.reduce,
		h.Update,
		h.Lock,
		h.Unlock,
	)

	h.Subscribe(h.aggregates...)

	return h
}

func (h *StatementHandler) TriggerBulk(ctx context.Context) {
	ctx, span := tracing.NewSpan(ctx)
	var err error
	defer span.EndWithError(err)

	err = h.ProjectionHandler.TriggerBulk(ctx, h.Lock, h.Unlock)
	logging.OnError(err).WithField("projection", h.ProjectionName).Warn("unable to trigger bulk")
}

func (h *StatementHandler) SearchQuery() (*eventstore.SearchQueryBuilder, uint64, error) {
	sequences, err := h.currentSequences(h.client.Query)
	if err != nil {
		return nil, 0, err
	}

	queryBuilder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).Limit(h.bulkLimit)
	for _, aggregateType := range h.aggregates {
		queryBuilder.
			AddQuery().
			AggregateTypes(aggregateType).
			SequenceGreater(sequences[aggregateType])
	}

	return queryBuilder, h.bulkLimit, nil
}

//Update implements handler.Update
func (h *StatementHandler) Update(ctx context.Context, stmts []*handler.Statement, reduce handler.Reduce) (unexecutedStmts []*handler.Statement, err error) {
	tx, err := h.client.BeginTx(ctx, nil)
	if err != nil {
		return stmts, errors.ThrowInternal(err, "CRDB-e89Gq", "begin failed")
	}

	sequences, err := h.currentSequences(tx.Query)
	if err != nil {
		tx.Rollback()
		return stmts, err
	}

	//checks for events between create statement and current sequence
	// because there could be events between current sequence and a creation event
	// and we cannot check via stmt.PreviousSequence
	if stmts[0].PreviousSequence == 0 {
		previousStmts, err := h.fetchPreviousStmts(ctx, stmts[0].Sequence, sequences, reduce)
		if err != nil {
			tx.Rollback()
			return stmts, err
		}
		stmts = append(previousStmts, stmts...)
	}

	lastSuccessfulIdx := h.executeStmts(tx, stmts, sequences)

	if lastSuccessfulIdx >= 0 {
		err = h.updateCurrentSequences(tx, sequences)
		if err != nil {
			tx.Rollback()
			return stmts, err
		}
	}

	if err = tx.Commit(); err != nil {
		return stmts, err
	}

	if lastSuccessfulIdx == -1 {
		return stmts, handler.ErrSomeStmtsFailed
	}

	unexecutedStmts = make([]*handler.Statement, len(stmts)-(lastSuccessfulIdx+1))
	copy(unexecutedStmts, stmts[lastSuccessfulIdx+1:])
	stmts = nil

	if len(unexecutedStmts) > 0 {
		return unexecutedStmts, handler.ErrSomeStmtsFailed
	}

	return unexecutedStmts, nil
}

func (h *StatementHandler) fetchPreviousStmts(
	ctx context.Context,
	stmtSeq uint64,
	sequences currentSequences,
	reduce handler.Reduce,
) (previousStmts []*handler.Statement, err error) {

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent)
	queriesAdded := false
	for _, aggregateType := range h.aggregates {
		if stmtSeq <= sequences[aggregateType] {
			continue
		}

		query.
			AddQuery().
			AggregateTypes(aggregateType).
			SequenceGreater(sequences[aggregateType]).
			SequenceLess(stmtSeq)

		queriesAdded = true
	}

	if !queriesAdded {
		return nil, nil
	}

	events, err := h.Eventstore.Filter(ctx, query)
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		stmt, err := reduce(event)
		if err != nil {
			return nil, err
		}
		previousStmts = append(previousStmts, stmt)
	}
	return previousStmts, nil
}

func (h *StatementHandler) executeStmts(
	tx *sql.Tx,
	stmts []*handler.Statement,
	sequences currentSequences,
) int {

	lastSuccessfulIdx := -1
	for i, stmt := range stmts {
		if stmt.Sequence <= sequences[stmt.AggregateType] {
			continue
		}
		if stmt.PreviousSequence > 0 && stmt.PreviousSequence != sequences[stmt.AggregateType] {
			logging.LogWithFields("CRDB-jJBJn", "projection", h.ProjectionName, "aggregateType", stmt.AggregateType, "seq", stmt.Sequence, "prevSeq", stmt.PreviousSequence, "currentSeq", sequences[stmt.AggregateType]).Warn("sequences do not match")
			break
		}
		err := h.executeStmt(tx, stmt)
		if err == nil {
			sequences[stmt.AggregateType], lastSuccessfulIdx = stmt.Sequence, i
			continue
		}

		shouldContinue := h.handleFailedStmt(tx, stmt, err)
		if !shouldContinue {
			break
		}

		sequences[stmt.AggregateType], lastSuccessfulIdx = stmt.Sequence, i
	}
	return lastSuccessfulIdx
}

//executeStmt handles sql statements
//an error is returned if the statement could not be inserted properly
func (h *StatementHandler) executeStmt(tx *sql.Tx, stmt *handler.Statement) error {
	if stmt.IsNoop() {
		return nil
	}
	_, err := tx.Exec("SAVEPOINT push_stmt")
	if err != nil {
		return errors.ThrowInternal(err, "CRDB-i1wp6", "unable to create savepoint")
	}
	err = stmt.Execute(tx, h.ProjectionName)
	if err != nil {
		_, rollbackErr := tx.Exec("ROLLBACK TO SAVEPOINT push_stmt")
		if rollbackErr != nil {
			return errors.ThrowInternal(rollbackErr, "CRDB-zzp3P", "rollback to savepoint failed")
		}
		return errors.ThrowInternal(err, "CRDB-oRkaN", "unable execute stmt")
	}
	_, err = tx.Exec("RELEASE push_stmt")
	if err != nil {
		return errors.ThrowInternal(err, "CRDB-qWgwT", "unable to release savepoint")
	}
	return nil
}
