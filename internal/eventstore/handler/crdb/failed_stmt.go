package crdb

import (
	"database/sql"

	"github.com/caos/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
)

const (
	setFailureCountStmtFormat = "UPSERT INTO %s" +
		" (projection_name, failed_sequence, failure_count, error, instance_id)" +
		" VALUES ($1, $2, $3, $4, $5)"
	failureCountStmtFormat = "WITH failures AS (SELECT failure_count FROM %s WHERE projection_name = $1 AND failed_sequence = $2 AND instance_id = $3)" +
		" SELECT IF(" +
		"EXISTS(SELECT failure_count FROM failures)," +
		" (SELECT failure_count FROM failures)," +
		" 0" +
		") AS failure_count"
)

func (h *StatementHandler) handleFailedStmt(tx *sql.Tx, stmt *handler.Statement, execErr error) (shouldContinue bool) {
	failureCount, err := h.failureCount(tx, stmt.Sequence, stmt.InstanceID)
	if err != nil {
		logging.WithFields("projection", h.ProjectionName, "sequence", stmt.Sequence).WithError(err).Warn("unable to get failure count")
		return false
	}
	failureCount += 1
	err = h.setFailureCount(tx, stmt.Sequence, failureCount, execErr, stmt.InstanceID)
	logging.WithFields("projection", h.ProjectionName, "sequence", stmt.Sequence).OnError(err).Warn("unable to update failure count")

	return failureCount >= h.maxFailureCount
}

func (h *StatementHandler) failureCount(tx *sql.Tx, seq uint64, instanceID string) (count uint, err error) {
	row := tx.QueryRow(h.failureCountStmt, h.ProjectionName, seq, instanceID)
	if err = row.Err(); err != nil {
		return 0, errors.ThrowInternal(err, "CRDB-Unnex", "unable to update failure count")
	}
	if err = row.Scan(&count); err != nil {
		return 0, errors.ThrowInternal(err, "CRDB-RwSMV", "unable to scan count")
	}
	return count, nil
}

func (h *StatementHandler) setFailureCount(tx *sql.Tx, seq uint64, count uint, err error, instanceID string) error {
	_, dbErr := tx.Exec(h.setFailureCountStmt, h.ProjectionName, seq, count, err.Error(), instanceID)
	if dbErr != nil {
		return errors.ThrowInternal(dbErr, "CRDB-4Ht4x", "set failure count failed")
	}
	return nil
}
