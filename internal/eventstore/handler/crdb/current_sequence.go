package crdb

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	currentSequenceStmtFormat        = `SELECT current_sequence, aggregate_type, instance_id FROM %s WHERE projection_name = $1 FOR UPDATE`
	updateCurrentSequencesStmtFormat = `UPSERT INTO %s (projection_name, aggregate_type, current_sequence, instance_id, timestamp) VALUES `
)

type currentSequences map[eventstore.AggregateType][]*instanceSequence

type instanceSequence struct {
	instanceID string
	sequence   uint64
}

func (h *StatementHandler) currentSequences(ctx context.Context, query func(context.Context, string, ...interface{}) (*sql.Rows, error)) (currentSequences, error) {
	rows, err := query(ctx, h.currentSequenceStmt, h.ProjectionName)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sequences := make(currentSequences, len(h.aggregates))
	for rows.Next() {
		var (
			aggregateType eventstore.AggregateType
			sequence      uint64
			instanceID    string
		)

		err = rows.Scan(&sequence, &aggregateType, &instanceID)
		if err != nil {
			return nil, errors.ThrowInternal(err, "CRDB-dbatK", "scan failed")
		}

		sequences[aggregateType] = append(sequences[aggregateType], &instanceSequence{
			sequence:   sequence,
			instanceID: instanceID,
		})
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
	for aggregate, instanceSequence := range sequences {
		for _, sequence := range instanceSequence {
			valueQueries = append(valueQueries, "($"+strconv.Itoa(valueCounter+1)+", $"+strconv.Itoa(valueCounter+2)+", $"+strconv.Itoa(valueCounter+3)+", $"+strconv.Itoa(valueCounter+4)+", NOW())")
			valueCounter += 4
			values = append(values, h.ProjectionName, aggregate, sequence.sequence, sequence.instanceID)
		}
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
