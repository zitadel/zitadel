package crdb

import (
	"database/sql"
	"database/sql/driver"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type mockExpectation func(sqlmock.Sqlmock)

func expectFailureCount(tableName string, projectionName, instanceID string, failedSeq, failureCount uint64) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(`WITH failures AS \(SELECT failure_count FROM `+tableName+` WHERE projection_name = \$1 AND failed_sequence = \$2 AND instance_id = \$3\) SELECT COALESCE\(\(SELECT failure_count FROM failures\), 0\) AS failure_count`).
			WithArgs(projectionName, failedSeq, instanceID).
			WillReturnRows(
				sqlmock.NewRows([]string{"failure_count"}).
					AddRow(failureCount),
			)
	}
}

func expectUpdateFailureCount(tableName string, projectionName, instanceID string, seq, failureCount uint64) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec(`INSERT INTO `+tableName+` \(projection_name, failed_sequence, failure_count, error, instance_id, last_failed\) VALUES \(\$1, \$2, \$3, \$4\, \$5\, \$6\) ON CONFLICT \(projection_name, failed_sequence, instance_id\) DO UPDATE SET failure_count = EXCLUDED\.failure_count, error = EXCLUDED\.error, last_failed = EXCLUDED\.last_failed`).
			WithArgs(projectionName, seq, failureCount, sqlmock.AnyArg(), instanceID, "NOW()").WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func expectCreate(projectionName string, columnNames, placeholders []string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		args := make([]driver.Value, len(columnNames))
		for i := 0; i < len(columnNames); i++ {
			args[i] = sqlmock.AnyArg()
			placeholders[i] = `\` + placeholders[i]
		}
		m.ExpectExec("INSERT INTO " + projectionName + ` \(` + strings.Join(columnNames, ", ") + `\) VALUES \(` + strings.Join(placeholders, ", ") + `\)`).
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func expectCreateErr(projectionName string, columnNames, placeholders []string, err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		args := make([]driver.Value, len(columnNames))
		for i := 0; i < len(columnNames); i++ {
			args[i] = sqlmock.AnyArg()
			placeholders[i] = `\` + placeholders[i]
		}
		m.ExpectExec("INSERT INTO " + projectionName + ` \(` + strings.Join(columnNames, ", ") + `\) VALUES \(` + strings.Join(placeholders, ", ") + `\)`).
			WithArgs(args...).
			WillReturnError(err)
	}
}

func expectBegin() func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectBegin()
	}
}

func expectBeginErr(err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectBegin().WillReturnError(err)
	}
}

func expectCommit() func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectCommit()
	}
}

func expectCommitErr(err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectCommit().WillReturnError(err)
	}
}

func expectRollback() func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectRollback()
	}
}

func expectSavePoint() func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("SAVEPOINT push_stmt").
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func expectSavePointErr(err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("SAVEPOINT push_stmt").
			WillReturnError(err)
	}
}

func expectSavePointRollback() func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("ROLLBACK TO SAVEPOINT push_stmt").
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func expectSavePointRollbackErr(err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("ROLLBACK TO SAVEPOINT push_stmt").
			WillReturnError(err)
	}
}

func expectSavePointRelease() func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("RELEASE push_stmt").
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func expectCurrentSequence(isTx bool, tableName, projection string, seq uint64, aggregateType string, instanceIDs []string) func(sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"current_sequence", "aggregate_type", "instance_id"})
	for _, instanceID := range instanceIDs {
		rows.AddRow(seq, aggregateType, instanceID)
	}
	return func(m sqlmock.Sqlmock) {
		stmt := `SELECT current_sequence, aggregate_type, instance_id FROM ` + tableName + ` WHERE projection_name = \$1 AND instance_id = ANY \(\$2\)`
		if isTx {
			stmt += " FOR UPDATE"
		}
		m.ExpectQuery(stmt).
			WithArgs(
				projection,
				database.StringArray(instanceIDs),
			).
			WillReturnRows(
				rows,
			)
	}
}

func expectCurrentSequenceErr(isTx bool, tableName, projection string, instanceIDs []string, err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		stmt := `SELECT current_sequence, aggregate_type, instance_id FROM ` + tableName + ` WHERE projection_name = \$1 AND instance_id = ANY \(\$2\)`
		if isTx {
			stmt += " FOR UPDATE"
		}
		m.ExpectQuery(stmt).
			WithArgs(
				projection,
				database.StringArray(instanceIDs),
			).
			WillReturnError(err)
	}
}

func expectCurrentSequenceNoRows(tableName, projection string, instanceIDs []string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(`SELECT current_sequence, aggregate_type, instance_id FROM `+tableName+` WHERE projection_name = \$1 AND instance_id = ANY \(\$2\) FOR UPDATE`).
			WithArgs(
				projection,
				database.StringArray(instanceIDs),
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"current_sequence", "aggregate_type", "instance_id"}),
			)
	}
}

func expectCurrentSequenceScanErr(tableName, projection string, instanceIDs []string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(`SELECT current_sequence, aggregate_type, instance_id FROM `+tableName+` WHERE projection_name = \$1 AND instance_id = ANY \(\$2\) FOR UPDATE`).
			WithArgs(
				projection,
				database.StringArray(instanceIDs),
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"current_sequence", "aggregate_type", "instance_id"}).
					RowError(0, sql.ErrTxDone).
					AddRow(0, "agg", "instanceID"),
			)
	}
}

func expectUpdateCurrentSequence(tableName, projection string, seq uint64, aggregateType, instanceID string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("INSERT INTO "+tableName+` \(projection_name, aggregate_type, current_sequence, instance_id, timestamp\) VALUES \(\$1, \$2, \$3, \$4, NOW\(\)\) ON CONFLICT \(projection_name, aggregate_type, instance_id\) DO UPDATE SET current_sequence = EXCLUDED.current_sequence, timestamp = EXCLUDED.timestamp`).
			WithArgs(
				projection,
				aggregateType,
				seq,
				instanceID,
			).
			WillReturnResult(
				sqlmock.NewResult(1, 1),
			)
	}
}

func expectUpdateThreeCurrentSequence(t *testing.T, tableName, projection string, sequences currentSequences) func(sqlmock.Sqlmock) {
	args := make([][]interface{}, 0)
	for aggregateType, instanceSequences := range sequences {
		for _, sequence := range instanceSequences {
			args = append(args, []interface{}{
				projection,
				aggregateType,
				sequence.sequence,
				sequence.instanceID,
			})
		}
	}
	matcher := &currentSequenceMatcher{t: t, seq: args}
	matchers := make([]driver.Value, len(args)*4)
	for i := 0; i < len(args)*4; i++ {
		matchers[i] = matcher
	}
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("INSERT INTO " + tableName + ` \(projection_name, aggregate_type, current_sequence, instance_id, timestamp\) VALUES \(\$1, \$2, \$3, \$4, NOW\(\)\), \(\$5, \$6, \$7, \$8, NOW\(\)\), \(\$9, \$10, \$11, \$12, NOW\(\)\) ON CONFLICT \(projection_name, aggregate_type, instance_id\) DO UPDATE SET current_sequence = EXCLUDED.current_sequence, timestamp = EXCLUDED.timestamp`).
			WithArgs(
				matchers...,
			).
			WillReturnResult(
				sqlmock.NewResult(1, 1),
			)
	}
}

type currentSequenceMatcher struct {
	seq [][]interface{}
	i   int
	t   *testing.T
}

func (m *currentSequenceMatcher) Match(value driver.Value) bool {
	if m.i%4 == 0 {
		m.i = 0
	}
	left := make([]interface{}, 0, len(m.seq))
	for _, seq := range m.seq {
		found := seq[m.i]
		if found == nil {
			continue
		}
		switch v := value.(type) {
		case string:
			if found == v || found == eventstore.AggregateType(v) {
				seq[m.i] = nil
				m.i++
				return true
			}
		case int64:
			if found == uint64(v) {
				seq[m.i] = nil
				m.i++
				return true
			}
		}
		left = append(left, found)
	}
	m.t.Errorf("expected: %v, possible left values: %v", value, left)
	m.t.FailNow()
	return false
}

func expectUpdateCurrentSequenceErr(tableName, projection string, seq uint64, err error, aggregateType, instanceID string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("INSERT INTO "+tableName+` \(projection_name, aggregate_type, current_sequence, instance_id, timestamp\) VALUES \(\$1, \$2, \$3, \$4, NOW\(\)\) ON CONFLICT \(projection_name, aggregate_type, instance_id\) DO UPDATE SET current_sequence = EXCLUDED.current_sequence, timestamp = EXCLUDED.timestamp`).
			WithArgs(
				projection,
				aggregateType,
				seq,
				instanceID,
			).
			WillReturnError(err)
	}
}

func expectUpdateCurrentSequenceNoRows(tableName, projection string, seq uint64, aggregateType, instanceID string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("INSERT INTO "+tableName+` \(projection_name, aggregate_type, current_sequence, instance_id, timestamp\) VALUES \(\$1, \$2, \$3, \$4, NOW\(\)\) ON CONFLICT \(projection_name, aggregate_type, instance_id\) DO UPDATE SET current_sequence = EXCLUDED.current_sequence, timestamp = EXCLUDED.timestamp`).
			WithArgs(
				projection,
				aggregateType,
				seq,
				instanceID,
			).
			WillReturnResult(
				sqlmock.NewResult(0, 0),
			)
	}
}

func expectLock(lockTable, workerName string, d time.Duration, instanceID string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec(`INSERT INTO `+lockTable+
			` \(locker_id, locked_until, projection_name, instance_id\) VALUES \(\$1, now\(\)\+\$2::INTERVAL, \$3\, \$4\)`+
			` ON CONFLICT \(projection_name, instance_id\)`+
			` DO UPDATE SET locker_id = \$1, locked_until = now\(\)\+\$2::INTERVAL`+
			` WHERE `+lockTable+`\.projection_name = \$3 AND `+lockTable+`\.instance_id = ANY \(\$5\) AND \(`+lockTable+`\.locker_id = \$1 OR `+lockTable+`\.locked_until < now\(\)\)`).
			WithArgs(
				workerName,
				d,
				projectionName,
				instanceID,
				database.StringArray{instanceID},
			).
			WillReturnResult(
				sqlmock.NewResult(1, 1),
			)
	}
}

func expectLockMultipleInstances(lockTable, workerName string, d time.Duration, instanceID1, instanceID2 string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec(`INSERT INTO `+lockTable+
			` \(locker_id, locked_until, projection_name, instance_id\) VALUES \(\$1, now\(\)\+\$2::INTERVAL, \$3\, \$4\), \(\$1, now\(\)\+\$2::INTERVAL, \$3\, \$5\)`+
			` ON CONFLICT \(projection_name, instance_id\)`+
			` DO UPDATE SET locker_id = \$1, locked_until = now\(\)\+\$2::INTERVAL`+
			` WHERE `+lockTable+`\.projection_name = \$3 AND `+lockTable+`\.instance_id = ANY \(\$6\) AND \(`+lockTable+`\.locker_id = \$1 OR `+lockTable+`\.locked_until < now\(\)\)`).
			WithArgs(
				workerName,
				d,
				projectionName,
				instanceID1,
				instanceID2,
				database.StringArray{instanceID1, instanceID2},
			).
			WillReturnResult(
				sqlmock.NewResult(1, 1),
			)
	}
}

func expectLockNoRows(lockTable, workerName string, d time.Duration, instanceID string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec(`INSERT INTO `+lockTable+
			` \(locker_id, locked_until, projection_name, instance_id\) VALUES \(\$1, now\(\)\+\$2::INTERVAL, \$3\, \$4\)`+
			` ON CONFLICT \(projection_name, instance_id\)`+
			` DO UPDATE SET locker_id = \$1, locked_until = now\(\)\+\$2::INTERVAL`+
			` WHERE `+lockTable+`\.projection_name = \$3 AND `+lockTable+`\.instance_id = ANY \(\$5\) AND \(`+lockTable+`\.locker_id = \$1 OR `+lockTable+`\.locked_until < now\(\)\)`).
			WithArgs(
				workerName,
				d,
				projectionName,
				instanceID,
				database.StringArray{instanceID},
			).
			WillReturnResult(driver.ResultNoRows)
	}
}

func expectLockErr(lockTable, workerName string, d time.Duration, instanceID string, err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec(`INSERT INTO `+lockTable+
			` \(locker_id, locked_until, projection_name, instance_id\) VALUES \(\$1, now\(\)\+\$2::INTERVAL, \$3\, \$4\)`+
			` ON CONFLICT \(projection_name, instance_id\)`+
			` DO UPDATE SET locker_id = \$1, locked_until = now\(\)\+\$2::INTERVAL`+
			` WHERE `+lockTable+`\.projection_name = \$3 AND `+lockTable+`\.instance_id = ANY \(\$5\) AND \(`+lockTable+`\.locker_id = \$1 OR `+lockTable+`\.locked_until < now\(\)\)`).
			WithArgs(
				workerName,
				d,
				projectionName,
				instanceID,
				database.StringArray{instanceID},
			).
			WillReturnError(err)
	}
}
