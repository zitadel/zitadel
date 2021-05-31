package crdb

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/id"
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

	Reducers []handler.EventReducer
}

type StatementHandler struct {
	*handler.ProjectionHandler

	client              *sql.DB
	sequenceTable       string
	maxFailureCount     uint
	failureCountStmt    string
	setFailureCountStmt string
	lockStmt            string

	aggregates []eventstore.AggregateType
	eventTypes []eventstore.EventType
	reduces    map[eventstore.EventType]handler.Reduce

	workerName string
	bulkLimit  uint64
}

func NewStatementHandler(
	ctx context.Context,
	config StatementHandlerConfig,
) StatementHandler {
	workerName, err := os.Hostname()
	if err != nil || workerName == "" {
		workerName, err = id.SonyFlakeGenerator.Next()
		logging.Log("SPOOL-bdO56").OnError(err).Panic("unable to generate lockID")
	}

	aggregateTypes := make([]eventstore.AggregateType, 0, len(config.Reducers))
	eventTypes := make([]eventstore.EventType, 0, len(config.Reducers))
	reduces := make(map[eventstore.EventType]handler.Reduce, len(config.Reducers))
	// subscriptionTopics := make(map[eventstore.AggregateType][]eventstore.EventType)
	for _, reducer := range config.Reducers {
		aggregateTypes = append(aggregateTypes, reducer.Aggregate)
		eventTypes = append(eventTypes, reducer.Event)
		reduces[reducer.Event] = reducer.Reduce
		// subscriptionTopics[reducer.Aggregate] = append(subscriptionTopics[reducer.Aggregate], reducer.Event)
	}

	h := StatementHandler{
		ProjectionHandler:   handler.NewProjectionHandler(config.ProjectionHandlerConfig),
		client:              config.Client,
		sequenceTable:       config.SequenceTable,
		maxFailureCount:     config.MaxFailureCount,
		failureCountStmt:    fmt.Sprintf(failureCountStmtFormat, config.FailedEventsTable),
		setFailureCountStmt: fmt.Sprintf(setFailureCountStmtFormat, config.FailedEventsTable),
		lockStmt:            fmt.Sprintf(lockStmtFormat, config.LockTable),
		aggregates:          aggregateTypes,
		eventTypes:          eventTypes,
		reduces:             reduces,
		workerName:          workerName,
		bulkLimit:           config.BulkLimit,
	}

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

func (h *StatementHandler) SearchQuery() (*eventstore.SearchQueryBuilder, uint64, error) {
	seq, err := h.currentSequence(h.client.QueryRow)
	if err != nil {
		return nil, 0, err
	}

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, h.aggregates...).SequenceGreater(seq).Limit(h.bulkLimit)
	// if len(h.eventTypes) > 0 {
	// 	query.EventTypes(h.eventTypes...)
	// }

	return query, h.bulkLimit, nil
}

//Update implements handler.Update
func (h *StatementHandler) Update(ctx context.Context, stmts []handler.Statement, reduce handler.Reduce) (unexecutedStmts []handler.Statement, err error) {
	tx, err := h.client.BeginTx(ctx, nil)
	if err != nil {
		return stmts, err
	}

	currentSeq, err := h.currentSequence(tx.QueryRow)
	if err != nil {
		tx.Rollback()
		return stmts, err
	}

	//checks for events between create statement and current sequence
	// because there could be events between current sequence and a creation event
	// and we cannot check via stmt.PreviousSequence
	if stmts[0].PreviousSequence == 0 {
		previousStmts, err := h.fetchPreviousStmts(ctx, stmts[0].Sequence, currentSeq, reduce)
		if err != nil {
			tx.Rollback()
			return stmts, err
		}
		stmts = append(previousStmts, stmts...)
	}

	lastSuccessfulIdx := h.executeStmts(tx, stmts, currentSeq)

	if lastSuccessfulIdx >= 0 {
		seqErr := h.updateCurrentSequence(tx, stmts[lastSuccessfulIdx])
		if seqErr != nil {
			tx.Rollback()
			return stmts, seqErr
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return stmts, commitErr
	}

	if lastSuccessfulIdx == 0 {
		return stmts, nil
	}

	unexecutedStmts = make([]handler.Statement, len(stmts)-(lastSuccessfulIdx+1))
	copy(unexecutedStmts, stmts[lastSuccessfulIdx+1:])
	stmts = nil

	return unexecutedStmts, nil
}

func (h *StatementHandler) fetchPreviousStmts(
	ctx context.Context,
	stmtSeq,
	currentSeq uint64,
	reduce handler.Reduce,
) (previousStmts []handler.Statement, err error) {

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, h.aggregates...).SequenceGreater(currentSeq).SequenceLess(stmtSeq)
	events, err := h.Eventstore.FilterEvents(ctx, query)
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		stmts, err := reduce(event)
		if err != nil {
			return nil, err
		}
		previousStmts = append(previousStmts, stmts...)
	}
	return previousStmts, nil
}

func (h *StatementHandler) executeStmts(
	tx *sql.Tx,
	stmts []handler.Statement,
	currentSeq uint64,
) int {

	lastSuccessfulIdx := -1
	for i, stmt := range stmts {
		if stmt.Sequence < currentSeq {
			continue
		}
		if stmt.PreviousSequence > currentSeq {
			logging.LogWithFields("CRDB-jJBJn", "prevSeq", stmt.PreviousSequence, "currentSeq", currentSeq).Warn("sequences do not match")
			break
		}
		err := h.executeStmt(tx, stmt)
		if err == nil {
			currentSeq, lastSuccessfulIdx = stmt.Sequence, i
			continue
		}

		shouldContinue := h.handleFailedStmt(tx, stmt, err)
		if !shouldContinue {
			break
		}

		currentSeq, lastSuccessfulIdx = stmt.Sequence, i
	}
	return lastSuccessfulIdx
}

//executeStmt handles sql statements
//an error is returned if the statement could not be inserted properly
func (h *StatementHandler) executeStmt(tx *sql.Tx, stmt handler.Statement) error {
	if stmt.IsNoop() {
		return nil
	}
	_, err := tx.Exec("SAVEPOINT push_stmt")
	if err != nil {
		return err
	}
	err = stmt.Execute(tx, h.ProjectionName)
	if err != nil {
		_, rollbackErr := tx.Exec("ROLLBACK TO SAVEPOINT push_stmt")
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	_, err = tx.Exec("RELEASE push_stmt")
	return err
}

func (h *StatementHandler) currentSequence(query func(string, ...interface{}) *sql.Row) (seq uint64, _ error) {
	row := query(`WITH seq AS (SELECT current_sequence FROM `+h.sequenceTable+` WHERE view_name = $1 FOR UPDATE)
	SELECT IF(
		EXISTS(SELECT current_sequence FROM seq),
		(SELECT current_sequence FROM seq),
		0
	) AS current_sequence`, h.ProjectionName)
	if row.Err() != nil {
		return 0, row.Err()
	}

	if err := row.Scan(&seq); err != nil {
		return 0, err
	}

	return seq, nil
}

func (h *StatementHandler) updateCurrentSequence(tx *sql.Tx, stmt handler.Statement) error {
	res, err := tx.Exec(`UPSERT INTO `+h.sequenceTable+` (view_name, current_sequence, timestamp) VALUES ($1, $2, NOW())`, h.ProjectionName, stmt.Sequence)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return errSeqNotUpdated
	}
	return nil
}
