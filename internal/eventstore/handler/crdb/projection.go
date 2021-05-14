package crdb

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/id"
)

type StatementHandler struct {
	viewName      string
	sequenceTable string
	client        *sql.DB
	eventstore    *eventstore.Eventstore
	aggregates    []eventstore.AggregateType

	workerName string
	lockStmt   string
	bulkLimit  uint64
}

func NewStatementHandler(
	client *sql.DB,
	es *eventstore.Eventstore,
	viewName,
	sequenceTable,
	lockTable string,
	bulkLimit uint64,
	aggregates ...eventstore.AggregateType,
) StatementHandler {
	workerName, err := os.Hostname()
	if err != nil || workerName == "" {
		workerName, err = id.SonyFlakeGenerator.Next()
		logging.Log("SPOOL-bdO56").OnError(err).Panic("unable to generate lockID")
	}

	return StatementHandler{
		client:        client,
		eventstore:    es,
		viewName:      viewName,
		sequenceTable: sequenceTable,
		workerName:    workerName,
		lockStmt:      fmt.Sprintf(lockStmtFormat, lockTable, lockTable, lockTable, lockTable),
		bulkLimit:     bulkLimit,
		aggregates:    aggregates,
	}
}

func (h *StatementHandler) SearchQuery() (*eventstore.SearchQueryBuilder, uint64, error) {
	seq, err := h.currentSequence(h.client.QueryRow)
	if err != nil {
		return nil, 0, err
	}
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, h.aggregates...).SequenceGreater(seq).Limit(h.bulkLimit), h.bulkLimit, nil
}

// func (h *StatementHandler) stmtError(stmt handler.Statement, err error) error {

// 	return nil
// }

func (h *StatementHandler) Update(ctx context.Context, stmts []handler.Statement, reduce handler.Reduce) error {
	if len(stmts) == 0 {
		return nil
	}

	tx, err := h.client.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	currentSeq, err := h.currentSequence(tx.QueryRow)
	if err != nil {
		tx.Rollback()
		return err
	}

	//checks for events between create statement and current sequence
	// because there could be events between current sequence and the creation event
	// and we cannot check via stmt.PreviousSequence
	if stmts[0].PreviousSequence == 0 {
		previousStmts, err := h.preparePreviousStmts(ctx, stmts[0].Sequence, currentSeq, reduce)
		if err != nil {
			return err
		}
		stmts = append(previousStmts, stmts...)
	}

	lastSuccessfulIdx := -1
	for i, stmt := range stmts {
		if stmt.PreviousSequence > 0 && stmt.PreviousSequence < currentSeq {
			continue
		}
		if stmt.PreviousSequence > currentSeq {
			break
		}
		if err = executeStmt(tx, stmt); err != nil {
			//TODO: insert into error view
			//TODO: should we retry because nothing will change
			logging.LogWithFields("CRDB-wS8Ns", "seq", stmt.Sequence, "projection", stmt.TableName).WithError(err).Warn("unable to execute statement")
			break
		}
		currentSeq = stmt.Sequence
		lastSuccessfulIdx = i
	}

	if lastSuccessfulIdx >= 0 {
		seqErr := h.updateCurrentSequence(tx, stmts[lastSuccessfulIdx])
		if seqErr != nil {
			tx.Rollback()
			return seqErr
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return commitErr
	}

	return err
}

func (h *StatementHandler) preparePreviousStmts(ctx context.Context, stmtSeq, currentSeq uint64, reduce handler.Reduce) (previousStmts []handler.Statement, err error) {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, h.aggregates...).SequenceGreater(currentSeq).SequenceLess(stmtSeq)
	events, err := h.eventstore.FilterEvents(ctx, query)
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		previousStmts, err = reduce(event)
		if err != nil {
			return nil, err
		}
	}
	return previousStmts, nil
}

//executeStmt handles sql statements
//an error is returned if the statement could not be inserted properly
func executeStmt(tx *sql.Tx, stmt handler.Statement) error {
	_, err := tx.Query("SAVEPOINT push_stmt")
	if err != nil {
		return err
	}
	err = stmt.Execute(tx)
	if err != nil {
		_, rollbackErr := tx.Query("ROLLBACK TO SAVEPOINT push_stmt")
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	_, err = tx.Query("RELEASE push_stmt")
	return err
}

const currentSequenceFormat = `with seq as (select current_sequence from %s where view_name = $1 FOR UPDATE)
select 
    if(
        count(current_sequence) > 0, 
        (select current_sequence from seq),
        0
    ) 
from seq`

func (h *StatementHandler) currentSequence(query func(string, ...interface{}) *sql.Row) (seq uint64, _ error) {
	row := query(fmt.Sprintf(currentSequenceFormat, h.sequenceTable), h.viewName)
	if row.Err() != nil {
		return 0, row.Err()
	}

	if err := row.Scan(&seq); err != nil {
		return 0, err
	}

	return seq, nil
}

const upsertCurrentSequenceFormat = `UPSERT INTO %s (view_name, current_sequence, timestamp) VALUES ($1, $2, NOW())`

func (h *StatementHandler) updateCurrentSequence(tx *sql.Tx, stmt handler.Statement) error {
	_, err := tx.Exec(fmt.Sprintf(upsertCurrentSequenceFormat, h.sequenceTable), stmt.TableName, stmt.Sequence)
	return err
}
