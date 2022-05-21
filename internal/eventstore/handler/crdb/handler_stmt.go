package crdb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
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

	Reducers  []handler.AggregateReducer
	InitCheck *handler.Check
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
		ProjectionHandler:       handler.NewProjectionHandler(config.ProjectionHandlerConfig),
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
		Locker:                  NewLocker(config.Client, config.LockTable, config.ProjectionHandlerConfig.ProjectionName),
	}

	err := h.Init(ctx, config.InitCheck)
	logging.OnError(err).Fatal("unable to initialize projections")

	go h.ProjectionHandler.Process(
		ctx,
		h.reduce,
		h.Update,
		h.Lock,
		h.Unlock,
		h.SearchQuery,
	)

	h.ProjectionHandler.Handler.Subscribe(h.aggregates...)

	return h
}

func (h *StatementHandler) SearchQuery(ctx context.Context) (*eventstore.SearchQueryBuilder, uint64, error) {
	sequences, err := h.currentSequences(ctx, h.client.QueryContext)
	if err != nil {
		return nil, 0, err
	}

	queryBuilder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).Limit(h.bulkLimit)
	for _, aggregateType := range h.aggregates {
		instances := make([]string, 0)
		for _, sequence := range sequences[aggregateType] {
			instances = appendToIgnoredInstances(instances, sequence.instanceID)
			queryBuilder.
				AddQuery().
				AggregateTypes(aggregateType).
				SequenceGreater(sequence.sequence).
				InstanceID(sequence.instanceID)
		}
		queryBuilder.
			AddQuery().
			AggregateTypes(aggregateType).
			SequenceGreater(0).
			ExcludedInstanceID(instances...)
	}

	return queryBuilder, h.bulkLimit, nil
}

func appendToIgnoredInstances(instances []string, id string) []string {
	for _, instance := range instances {
		if instance == id {
			return instances
		}
	}
	return append(instances, id)
}

//Update implements handler.Update
func (h *StatementHandler) Update(ctx context.Context, stmts []*handler.Statement, reduce handler.Reduce) (unexecutedStmts []*handler.Statement, err error) {
	if len(stmts) == 0 {
		return nil, nil
	}
	tx, err := h.client.BeginTx(ctx, nil)
	if err != nil {
		return stmts, errors.ThrowInternal(err, "CRDB-e89Gq", "begin failed")
	}

	sequences, err := h.currentSequences(ctx, tx.QueryContext)
	if err != nil {
		tx.Rollback()
		return stmts, err
	}

	//checks for events between create statement and current sequence
	// because there could be events between current sequence and a creation event
	// and we cannot check via stmt.PreviousSequence
	if stmts[0].PreviousSequence == 0 {
		previousStmts, err := h.fetchPreviousStmts(ctx, tx, stmts[0].Sequence, stmts[0].InstanceID, sequences, reduce)
		if err != nil {
			tx.Rollback()
			return stmts, err
		}
		stmts = append(previousStmts, stmts...)
	}

	lastSuccessfulIdx := h.executeStmts(tx, &stmts, sequences)

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

	if lastSuccessfulIdx == -1 && len(stmts) > 0 {
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

func (h *StatementHandler) fetchPreviousStmts(ctx context.Context, tx *sql.Tx, stmtSeq uint64, instanceID string, sequences currentSequences, reduce handler.Reduce) (previousStmts []*handler.Statement, err error) {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).SetTx(tx)
	queriesAdded := false
	for _, aggregateType := range h.aggregates {
		for _, sequence := range sequences[aggregateType] {
			if stmtSeq <= sequence.sequence && instanceID == sequence.instanceID {
				continue
			}

			query.
				AddQuery().
				AggregateTypes(aggregateType).
				SequenceGreater(sequence.sequence).
				SequenceLess(stmtSeq).
				InstanceID(sequence.instanceID)

			queriesAdded = true
		}
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
	stmts *[]*handler.Statement,
	sequences currentSequences,
) int {

	lastSuccessfulIdx := -1
stmts:
	for i := 0; i < len(*stmts); i++ {
		stmt := (*stmts)[i]
		for _, sequence := range sequences[stmt.AggregateType] {
			if stmt.Sequence <= sequence.sequence && stmt.InstanceID == sequence.instanceID {
				logging.WithFields("statement", stmt, "currentSequence", sequence).Debug("statement dropped")
				if i < len(*stmts)-1 {
					copy((*stmts)[i:], (*stmts)[i+1:])
				}
				*stmts = (*stmts)[:len(*stmts)-1]
				i--
				continue stmts
			}
			if stmt.PreviousSequence > 0 && stmt.PreviousSequence != sequence.sequence && stmt.InstanceID == sequence.instanceID {
				logging.WithFields("projection", h.ProjectionName, "aggregateType", stmt.AggregateType, "sequence", stmt.Sequence, "prevSeq", stmt.PreviousSequence, "currentSeq", sequence.sequence).Warn("sequences do not match")
				break stmts
			}
		}
		err := h.executeStmt(tx, stmt)
		if err == nil {
			updateSequences(sequences, stmt)
			lastSuccessfulIdx = i
			continue
		}

		shouldContinue := h.handleFailedStmt(tx, stmt, err)
		if !shouldContinue {
			break
		}

		updateSequences(sequences, stmt)
		lastSuccessfulIdx = i
		continue
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

func updateSequences(sequences currentSequences, stmt *handler.Statement) {
	for _, sequence := range sequences[stmt.AggregateType] {
		if sequence.instanceID == stmt.InstanceID {
			sequence.sequence = stmt.Sequence
			return
		}
	}
	sequences[stmt.AggregateType] = append(sequences[stmt.AggregateType], &instanceSequence{
		instanceID: stmt.InstanceID,
		sequence:   stmt.Sequence,
	})
}
