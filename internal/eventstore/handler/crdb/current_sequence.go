package crdb

// import (
// 	"context"
// 	"database/sql"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/zitadel/zitadel/internal/database"
// 	"github.com/zitadel/zitadel/internal/errors"
// )

// const (
// 	currentSequenceStmtFormat          = `SELECT event_timestamp, instance_id FROM %s WHERE projection_name = $1 AND instance_id = ANY ($2) FOR UPDATE`
// 	updateCurrentSequencesStmtFormat   = `INSERT INTO %s (projection_name, instance_id, event_timestamp, last_updated) VALUES `
// 	updateCurrentSequencesConflictStmt = ` ON CONFLICT (projection_name, instance_id) DO UPDATE SET timestamp = EXCLUDED.timestamp, last_updated = excluded.last_updated`
// )

// type currentState []*instanceState

// type instanceState struct {
// 	instanceID        string
// 	eventCreationDate time.Time
// }

// func (h *StatementHandler) currentSequences(ctx context.Context, query func(context.Context, string, ...interface{}) (*sql.Rows, error), instanceIDs database.StringArray) (currentState, error) {
// 	rows, err := query(ctx, h.currentSequenceStmt, h.ProjectionName, instanceIDs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	sequences := make(currentState, 0, len(instanceIDs))
// 	for rows.Next() {
// 		var (
// 			eventTimestamp time.Time
// 			instanceID     string
// 		)

// 		err = rows.Scan(&eventTimestamp, &instanceID)
// 		if err != nil {
// 			return nil, errors.ThrowInternal(err, "CRDB-dbatK", "scan failed")
// 		}

// 		sequences = append(sequences, &instanceState{
// 			instanceID:        instanceID,
// 			eventCreationDate: eventTimestamp,
// 		})

// 	}

// 	if err = rows.Close(); err != nil {
// 		return nil, errors.ThrowInternal(err, "CRDB-h5i5m", "close rows failed")
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, errors.ThrowInternal(err, "CRDB-O8zig", "errors in scanning rows")
// 	}

// 	return sequences, nil
// }

// func (h *StatementHandler) updateCurrentSequences(tx *sql.Tx, sequences currentState) error {
// 	valueQueries := make([]string, 0, len(sequences))
// 	valueCounter := 0
// 	values := make([]interface{}, 0, len(sequences)*3)

// 	for _, sequence := range sequences {
// 		valueQueries = append(valueQueries, "($"+strconv.Itoa(valueCounter+1)+", $"+strconv.Itoa(valueCounter+2)+", $"+strconv.Itoa(valueCounter+3)+", NOW())")
// 		valueCounter += 3
// 		values = append(values, h.ProjectionName, sequence.instanceID, sequence.eventCreationDate)
// 	}

// 	res, err := tx.Exec(h.updateSequencesBaseStmt+strings.Join(valueQueries, ", ")+updateCurrentSequencesConflictStmt, values...)
// 	if err != nil {
// 		return errors.ThrowInternal(err, "CRDB-TrH2Z", "unable to exec update sequence")
// 	}
// 	if rows, _ := res.RowsAffected(); rows < 1 {
// 		return errSeqNotUpdated
// 	}
// 	return nil
// }
