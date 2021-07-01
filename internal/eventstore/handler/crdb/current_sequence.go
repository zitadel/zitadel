package crdb

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
)

const (
	currentSequenceStmtFormat        = `SELECT current_sequence, aggregate_type FROM %s WHERE view_name = $1 FOR UPDATE`
	updateCurrentSequencesStmtFormat = `UPSERT INTO %s (view_name, aggregate_type, current_sequence, timestamp) VALUES `
)

type currentSequences map[eventstore.AggregateType]uint64

func (h *StatementHandler) currentSequences(query func(string, ...interface{}) (*sql.Rows, error)) (currentSequences, error) {
	rows, err := query(h.currentSequenceStmt, h.ProjectionName)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sequences := make(currentSequences, len(h.aggregates))
	for rows.Next() {
		var (
			aggregateType eventstore.AggregateType
			sequence      uint64
		)

		err = rows.Scan(&sequence, &aggregateType)
		if err != nil {
			return nil, errors.ThrowInternal(err, "CRDB-dbatK", "scan failed")
		}

		sequences[aggregateType] = sequence
	}

	if err = rows.Close(); err != nil {
		return nil, errors.ThrowInternal(err, "CRDB-h5i5m", "close rows failed")
	}

	if err = rows.Err(); err != nil {
		return nil, errors.ThrowInternal(err, "CRDB-O8zig", "errors in scanning rows")
	}

	return sequences, nil
}

func (h *StatementHandler) updateCurrentSequences(tx *sql.Tx, sequences currentSequences) error {
	valueQueries := make([]string, 0, len(sequences))
	valueCounter := 0
	values := make([]interface{}, 0, len(sequences)*3)
	for aggregate, sequence := range sequences {
		valueQueries = append(valueQueries, "($"+strconv.Itoa(valueCounter+1)+", $"+strconv.Itoa(valueCounter+2)+", $"+strconv.Itoa(valueCounter+3)+", NOW())")
		valueCounter += 3
		values = append(values, h.ProjectionName, aggregate, sequence)
	}

	res, err := tx.Exec(h.updateSequencesBaseStmt+strings.Join(valueQueries, ", "), values...)
	if err != nil {
		return errors.ThrowInternal(err, "CRDB-TrH2Z", "unable to exec update sequence")
	}
	if rows, _ := res.RowsAffected(); rows < 1 {
		return errSeqNotUpdated
	}
	return nil
}
