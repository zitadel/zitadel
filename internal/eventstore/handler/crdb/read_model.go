package crdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/handler"
)

type StatementHandler struct {
	viewName      string
	sequenceTable string
	lockTable     string
	client        *sql.DB
}

func NewStatementHandler(
	client *sql.DB,
	viewName,
	sequenceTable,
	lockTable string,
) StatementHandler {
	return StatementHandler{
		client:        client,
		viewName:      viewName,
		sequenceTable: sequenceTable,
		lockTable:     lockTable,
	}
}

func (h *StatementHandler) Lock() error {
	query := "INSERT INTO " + h.lockTable + " (view_name) VALUES ($1)"
	_, err := h.client.Exec(query, h.viewName)
	return err
}

func (h *StatementHandler) Unlock() error {
	query := "DELETE FROM " + h.lockTable + " WHERE view_name = $1"
	_, err := h.client.Exec(query, h.viewName)
	return err
}

func (h *StatementHandler) Update(ctx context.Context, stmts []handler.Statement) error {
	if len(stmts) == 0 {
		return nil
	}

	tx, err := h.client.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	currentSeq, err := h.currentSequence(tx, stmts[0])
	if err != nil {
		tx.Rollback()
		return err
	}

	lastSuccessfulIdx := -1

	for i, stmt := range stmts {
		start := time.Now()
		if stmt.PreviousSequence > 0 && stmt.PreviousSequence < currentSeq {
			continue
		}
		if stmt.PreviousSequence > currentSeq {
			break
		}
		if err = executeStmt(tx, stmt); err != nil {
			break
		}
		currentSeq = stmt.Sequence
		lastSuccessfulIdx = i
		logging.LogWithFields("HANDL-j5vuD", "start", start, "end", time.Now(), "diff", time.Now().Sub(start), "iter", i).Warn("stmt")
	}

	if lastSuccessfulIdx >= 0 {
		err = h.updateCurrentSequence(tx, stmts[lastSuccessfulIdx])
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
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

func (h *StatementHandler) currentSequence(tx *sql.Tx, stmt handler.Statement) (seq uint64, _ error) {
	row := tx.QueryRow(fmt.Sprintf(currentSequenceFormat, h.sequenceTable), stmt.TableName)
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
