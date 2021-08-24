package handler

import (
	"database/sql"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/handler"
)

const (
	setFailureCountStmt = "UPSERT INTO zitadel.projections.failed_events" +
		" (projection_name, failed_sequence, failure_count, error)" +
		" VALUES ($1, $2, $3, $4)"
	failureCountStmt = "WITH failures AS (SELECT failure_count FROM zitadel.projections.failed_events WHERE projection_name = $1 AND failed_sequence = $2)" +
		" SELECT IF(" +
		"EXISTS(SELECT failure_count FROM failures)," +
		" (SELECT failure_count FROM failures)," +
		" 0" +
		") AS failure_count"
)

func (p *Pusher) handleFailedStmt(tx *sql.Tx, stmt handler.Statement, execErr error) (shouldContinue bool) {
	failureCount, err := p.failureCount(tx, stmt.Sequence)
	if err != nil {
		logging.LogWithFields("CRDB-WJaFk", "projection", p.projectionName, "seq", stmt.Sequence).WithError(err).Warn("unable to get failure count")
		return false
	}
	failureCount += 1
	err = p.setFailureCount(tx, stmt.Sequence, failureCount, execErr)
	logging.LogWithFields("CRDB-cI0dB", "projection", p.projectionName, "seq", stmt.Sequence).OnError(err).Warn("unable to update failure count")

	return failureCount >= p.maxFailureCount
}

func (p *Pusher) failureCount(tx *sql.Tx, seq uint64) (count uint, err error) {
	row := tx.QueryRow(failureCountStmt, p.projectionName, seq)
	if err = row.Err(); err != nil {
		return 0, errors.ThrowInternal(err, "CRDB-Unnex", "unable to update failure count")
	}
	if err = row.Scan(&count); err != nil {
		return 0, errors.ThrowInternal(err, "CRDB-RwSMV", "unable to scann count")
	}
	return count, nil
}

func (p *Pusher) setFailureCount(tx *sql.Tx, seq uint64, count uint, err error) error {
	_, dbErr := tx.Exec(setFailureCountStmt, p.projectionName, seq, count, err.Error())
	if dbErr != nil {
		return errors.ThrowInternal(dbErr, "CRDB-4Ht4x", "set failure count failed")
	}
	return nil
}
