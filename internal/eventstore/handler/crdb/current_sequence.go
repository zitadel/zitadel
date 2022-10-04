package crdb

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	latestEventStmtFormat     = `SELECT id, creation_date, aggregate_type, instance_id FROM %s WHERE projection_name = $1 AND instance_id = ANY ($2) FOR UPDATE`
	updateEventIDStmtFormat   = `INSERT INTO %s (projection_name, aggregate_type, event_id, instance_id, timestamp) VALUES `
	updateEventIDConflictStmt = ` ON CONFLICT (projection_name, aggregate_type, instance_id) DO UPDATE SET event_id = EXCLUDED.event_id, timestamp = EXCLUDED.timestamp`
)

type events map[eventstore.AggregateType][]*instanceEvents

type instanceEvents struct {
	instanceID   string
	eventID      string
	creationDate time.Time
}

func (h *StatementHandler) currentSequences(ctx context.Context, query func(context.Context, string, ...interface{}) (*sql.Rows, error), instanceIDs database.StringArray) (events, error) {
	rows, err := query(ctx, h.currentSequenceStmt, h.ProjectionName, instanceIDs)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ids := make(events, len(h.aggregates))
	for rows.Next() {
		var (
			aggregateType eventstore.AggregateType
			eventID       string
			instanceID    string
			creationDate  sql.NullTime
		)

		err = rows.Scan(&eventID, &creationDate, &aggregateType, &instanceID)
		if err != nil {
			return nil, errors.ThrowInternal(err, "CRDB-dbatK", "scan failed")
		}

		ids[aggregateType] = append(ids[aggregateType], &instanceEvents{
			eventID:      eventID,
			instanceID:   instanceID,
			creationDate: creationDate.Time,
		})
	}

	if err = rows.Close(); err != nil {
		return nil, errors.ThrowInternal(err, "CRDB-h5i5m", "close rows failed")
	}

	if err = rows.Err(); err != nil {
		return nil, errors.ThrowInternal(err, "CRDB-O8zig", "errors in scanning rows")
	}

	return ids, nil
}

func (h *StatementHandler) updateCurrentSequences(tx *sql.Tx, ids events) error {
	valueQueries := make([]string, 0, len(ids))
	valueCounter := 0
	values := make([]interface{}, 0, len(ids)*3)
	for aggregate, eventID := range ids {
		for _, sequence := range eventID {
			valueQueries = append(valueQueries, "($"+strconv.Itoa(valueCounter+1)+", $"+strconv.Itoa(valueCounter+2)+", $"+strconv.Itoa(valueCounter+3)+", $"+strconv.Itoa(valueCounter+4)+", NOW())")
			valueCounter += 4
			values = append(values, h.ProjectionName, aggregate, sequence.eventID, sequence.instanceID)
		}
	}

	res, err := tx.Exec(h.updateSequencesBaseStmt+strings.Join(valueQueries, ", ")+updateEventIDConflictStmt, values...)
	if err != nil {
		return errors.ThrowInternal(err, "CRDB-TrH2Z", "unable to exec update sequence")
	}
	if rows, _ := res.RowsAffected(); rows < 1 {
		return errSeqNotUpdated
	}
	return nil
}
