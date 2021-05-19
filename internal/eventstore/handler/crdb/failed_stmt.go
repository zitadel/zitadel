package crdb

import (
	"database/sql"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/handler"
)

const (
	setFailureCountStmtFormat = "UPSERT INTO %s" +
		" (projection_name, failed_sequence, failure_count, error)" +
		" VALUES ($1, $2, $3, $4)"
	failureCountStmtFormat = "SELECT failure_count FROM %s" +
		" WHERE projection_name = $1 AND failed_sequence = $2"
)

func (h *StatementHandler) handleFailedStmt(tx *sql.Tx, stmt handler.Statement, execErr error) (shouldContinue bool) {
	failureCount, err := h.failureCount(tx, stmt.Sequence)
	if err != nil {
		logging.LogWithFields("CRDB-WJaFk", "seq", stmt.Sequence, "projection").WithError(err).Warn("unable to get failure count")
		return false
	}
	if failureCount < h.maxFailureCount {
		err = h.setFailureCount(tx, stmt.Sequence, failureCount+1, execErr)
		logging.LogWithFields("CRDB-cI0dB", "seq", stmt.Sequence, "projection").OnError(err).Warn("unable to update failure count")
		return false
	}
	return true
}

func (h *StatementHandler) failureCount(tx *sql.Tx, seq uint64) (count uint, err error) {
	row := tx.QueryRow(h.failureCountStmt, h.projectionName, seq)
	if err = row.Err(); err != nil {
		return 0, err
	}
	err = row.Scan(&count)
	return count, err
}

func (h *StatementHandler) setFailureCount(tx *sql.Tx, seq uint64, count uint, err error) error {
	_, dbErr := tx.Exec(h.setFailureCountStmt, h.projectionName, seq, count, err.Error())
	return dbErr
}
