package crdb

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"strings"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type mockExpectation func(sqlmock.Sqlmock)

func expectFailureCount(tableName string, projectionName string, failedSeq, failureCount uint64) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {

		m.ExpectQuery(`SELECT failure_count FROM `+tableName+` WHERE projection_name = \$1 AND failed_sequence = \$2`).
			WithArgs(projectionName, failedSeq).
			WillReturnRows(
				sqlmock.NewRows([]string{"failure_count"}).
					AddRow(failureCount),
			)
	}
}

func expectFailureCountErr(tableName string, projectionName string, failedSeq uint64) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {

		m.ExpectExec(`SELECT failure_count FROM `+tableName+` WHERE projection_name = \$1 AND failed_sequence = \$2`).
			WithArgs(projectionName, failedSeq).
			WillReturnResult(sqlmock.NewResult(1, 1))
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

func expectRollbackErr(err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectRollback().WillReturnError(err)
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

func expectCurrentSequence(tableName, projection string, seq uint64, aggregateType string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(`SELECT current_sequence, aggregate_type FROM ` + tableName + ` WHERE view_name = \$1 FOR UPDATE`).
			WithArgs(
				projection,
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"current_sequence", "aggregate_type"}).
					AddRow(seq, aggregateType),
			)
	}
}

func expectCurrentSequenceErr(tableName, projection string, err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(`SELECT current_sequence, aggregate_type FROM ` + tableName + ` WHERE view_name = \$1 FOR UPDATE`).
			WithArgs(
				projection,
			).
			WillReturnError(err)
	}
}

func expectCurrentSequenceNoRows(tableName, projection string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(`SELECT current_sequence, aggregate_type FROM ` + tableName + ` WHERE view_name = \$1 FOR UPDATE`).
			WithArgs(
				projection,
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"current_sequence", "aggregate_type"}),
			)
	}
}

func expectCurrentSequenceScanErr(tableName, projection string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(`SELECT current_sequence, aggregate_type FROM ` + tableName + ` WHERE view_name = \$1 FOR UPDATE`).
			WithArgs(
				projection,
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"current_sequence", "aggregate_type"}).
					RowError(0, sql.ErrTxDone).
					AddRow(0, "agg"),
			)
	}
}

func expectUpdateCurrentSequence(tableName, projection string, seq uint64, aggregateType string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("UPSERT INTO "+tableName+` \(view_name, aggregate_type, current_sequence, timestamp\) VALUES \(\$1, \$2, \$3, NOW\(\)\)`).
			WithArgs(
				projection,
				aggregateType,
				seq,
			).
			WillReturnResult(
				sqlmock.NewResult(1, 1),
			)
	}
}

func expectUpdateTwoCurrentSequence(tableName, projection string, seq1, seq2 uint64, aggregateType1, aggregateType2 string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		aggMatcher := &unorderedMachter{expected: []driver.Value{aggregateType1, aggregateType2}}
		seqMatcher := &unorderedMachter{expected: []driver.Value{int64(seq1), int64(seq2)}}
		m.ExpectExec("UPSERT INTO "+tableName+` \(view_name, aggregate_type, current_sequence, timestamp\) VALUES \(\$1, \$2, \$3, NOW\(\)\), \(\$4, \$5, \$6, NOW\(\)\)`).
			WithArgs(
				projection,
				aggMatcher,
				seqMatcher,
				projection,
				aggMatcher,
				seqMatcher,
			).
			WillReturnResult(
				sqlmock.NewResult(1, 1),
			)
	}
}

type unorderedMachter struct {
	expected []driver.Value
}

func (m *unorderedMachter) Match(value driver.Value) bool {
	for i := len(m.expected) - 1; i >= 0; i-- {
		if reflect.DeepEqual(value, m.expected[i]) {
			m.expected = append(m.expected[:i], m.expected[i+1:]...)
			return true
		}
	}
	return false
}

func expectUpdateCurrentSequenceErr(tableName, projection string, seq uint64, err error, aggregateType string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("UPSERT INTO "+tableName+` \(view_name, aggregate_type, current_sequence, timestamp\) VALUES \(\$1, \$2, \$3, NOW\(\)\)`).
			WithArgs(
				projection,
				aggregateType,
				seq,
			).
			WillReturnError(err)
	}
}

func expectUpdateCurrentSequenceNoRows(tableName, projection string, seq uint64, aggregateType string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("UPSERT INTO "+tableName+` \(view_name, aggregate_type, current_sequence, timestamp\) VALUES \(\$1, \$2, \$3, NOW\(\)\)`).
			WithArgs(
				projection,
				aggregateType,
				seq,
			).
			WillReturnResult(
				sqlmock.NewResult(0, 0),
			)
	}
}

func expectLock(lockTable, workerName string, d time.Duration) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec(`INSERT INTO `+lockTable+
			` \(locker_id, locked_until, view_name\) VALUES \(\$1, now\(\)\+\$2::INTERVAL, \$3\)`+
			` ON CONFLICT \(view_name\)`+
			` DO UPDATE SET locker_id = \$1, locked_until = now\(\)\+\$2::INTERVAL`+
			` WHERE `+lockTable+`\.view_name = \$3 AND \(`+lockTable+`\.locker_id = \$1 OR `+lockTable+`\.locked_until < now\(\)\)`).
			WithArgs(
				workerName,
				d,
				projectionName,
			).
			WillReturnResult(
				sqlmock.NewResult(1, 1),
			)
	}
}

func expectLockNoRows(lockTable, workerName string, d time.Duration) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec(`INSERT INTO `+lockTable+
			` \(locker_id, locked_until, view_name\) VALUES \(\$1, now\(\)\+\$2::INTERVAL, \$3\)`+
			` ON CONFLICT \(view_name\)`+
			` DO UPDATE SET locker_id = \$1, locked_until = now\(\)\+\$2::INTERVAL`+
			` WHERE `+lockTable+`\.view_name = \$3 AND \(`+lockTable+`\.locker_id = \$1 OR `+lockTable+`\.locked_until < now\(\)\)`).
			WithArgs(
				workerName,
				d,
				projectionName,
			).
			WillReturnResult(driver.ResultNoRows)
	}
}

func expectLockErr(lockTable, workerName string, d time.Duration, err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec(`INSERT INTO `+lockTable+
			` \(locker_id, locked_until, view_name\) VALUES \(\$1, now\(\)\+\$2::INTERVAL, \$3\)`+
			` ON CONFLICT \(view_name\)`+
			` DO UPDATE SET locker_id = \$1, locked_until = now\(\)\+\$2::INTERVAL`+
			` WHERE `+lockTable+`\.view_name = \$3 AND \(`+lockTable+`\.locker_id = \$1 OR `+lockTable+`\.locked_until < now\(\)\)`).
			WithArgs(
				workerName,
				d,
				projectionName,
			).
			WillReturnError(err)
	}
}
