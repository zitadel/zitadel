package crdb

import (
	"database/sql/driver"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/zitadel/zitadel/internal/database"
)

type mockExpectation func(sqlmock.Sqlmock)

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
				database.TextArray[string]{instanceID},
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
				database.TextArray[string]{instanceID1, instanceID2},
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
				database.TextArray[string]{instanceID},
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
				database.TextArray[string]{instanceID},
			).
			WillReturnError(err)
	}
}
