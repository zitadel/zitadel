package sql

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	selectEscaped = `SELECT creation_date, event_type, event_sequence, previous_aggregate_sequence, previous_aggregate_root_sequence, event_data, editor_service, editor_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore\.events WHERE aggregate_type = \$1`
)

var (
	eventColumns                             = []string{"creation_date", "event_type", "event_sequence", "previous_aggregate_sequence", "previous_aggregate_root_sequence", "event_data", "editor_service", "editor_user", "resource_owner", "aggregate_type", "aggregate_id", "aggregate_version"}
	expectedFilterEventsLimitFormat          = regexp.MustCompile(selectEscaped + ` ORDER BY event_sequence LIMIT \$2`).String()
	expectedFilterEventsDescFormat           = regexp.MustCompile(selectEscaped + ` ORDER BY event_sequence DESC`).String()
	expectedFilterEventsAggregateIDLimit     = regexp.MustCompile(selectEscaped + ` AND aggregate_id = \$2 ORDER BY event_sequence LIMIT \$3`).String()
	expectedFilterEventsAggregateIDTypeLimit = regexp.MustCompile(selectEscaped + ` AND aggregate_id = \$2 ORDER BY event_sequence LIMIT \$3`).String()
	expectedGetAllEvents                     = regexp.MustCompile(selectEscaped + ` ORDER BY event_sequence`).String()

	expectedInsertStatement = regexp.MustCompile(`INSERT INTO eventstore\.events ` +
		`\(event_type, aggregate_type, aggregate_id, aggregate_version, creation_date, event_data, editor_user, editor_service, resource_owner, previous_aggregate_sequence, previous_aggregate_root_sequence\) ` +
		`SELECT \$1, \$2, \$3, \$4, COALESCE\(\$5, now\(\)\), \$6, \$7, \$8, \$9, \$10 ` +
		`WHERE EXISTS \(` +
		`SELECT 1 FROM eventstore\.events WHERE aggregate_type = \$11 AND aggregate_id = \$12 HAVING MAX\(event_sequence\) = \$13 OR \(\$14::BIGINT IS NULL AND COUNT\(\*\) = 0\)\) ` +
		`RETURNING event_sequence, creation_date`).String()
)

type dbMock struct {
	sqlClient *sql.DB
	mock      sqlmock.Sqlmock
}

func (db *dbMock) close() {
	db.sqlClient.Close()
}

func mockDB(t *testing.T) *dbMock {
	mockDB := dbMock{}
	var err error
	mockDB.sqlClient, mockDB.mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("error occured while creating stub db %v", err)
	}

	mockDB.mock.MatchExpectationsInOrder(true)

	return &mockDB
}

func (db *dbMock) expectBegin(err error) *dbMock {
	if err != nil {
		db.mock.ExpectBegin().WillReturnError(err)
	} else {
		db.mock.ExpectBegin()
	}
	return db
}

func (db *dbMock) expectSavepoint() *dbMock {
	db.mock.ExpectExec("SAVEPOINT").WillReturnResult(sqlmock.NewResult(1, 1))

	return db
}

func (db *dbMock) expectReleaseSavepoint(err error) *dbMock {
	expectation := db.mock.ExpectExec("RELEASE SAVEPOINT")
	if err == nil {
		expectation.WillReturnResult(sqlmock.NewResult(1, 1))
	} else {
		expectation.WillReturnError(err)
	}

	return db
}

func (db *dbMock) expectCommit(err error) *dbMock {
	if err != nil {
		db.mock.ExpectCommit().WillReturnError(err)
	} else {
		db.mock.ExpectCommit()
	}
	return db
}

func (db *dbMock) expectRollback(err error) *dbMock {
	if err != nil {
		db.mock.ExpectRollback().WillReturnError(err)
	} else {
		db.mock.ExpectRollback()
	}
	return db
}

func (db *dbMock) expectInsertEvent(e *models.Event, returnedSequence uint64) *dbMock {
	db.mock.ExpectQuery(expectedInsertStatement).
		WithArgs(
			e.Type, e.AggregateType, e.AggregateID, e.AggregateVersion, sqlmock.AnyArg(), Data(e.Data), e.EditorUser, e.EditorService, e.ResourceOwner, Sequence(e.PreviousSequence),
			e.AggregateType, e.AggregateID, Sequence(e.PreviousSequence), Sequence(e.PreviousSequence),
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"event_sequence", "creation_date"}).
				AddRow(returnedSequence, time.Now().UTC()),
		)

	return db
}

func (db *dbMock) expectInsertEventError(e *models.Event) *dbMock {
	db.mock.ExpectQuery(expectedInsertStatement).
		WithArgs(
			e.Type, e.AggregateType, e.AggregateID, e.AggregateVersion, sqlmock.AnyArg(), Data(e.Data), e.EditorUser, e.EditorService, e.ResourceOwner, Sequence(e.PreviousSequence),
			e.AggregateType, e.AggregateID, Sequence(e.PreviousSequence), Sequence(e.PreviousSequence),
		).
		WillReturnError(sql.ErrTxDone)

	return db
}

func (db *dbMock) expectFilterEventsLimit(aggregateType string, limit uint64, eventCount int) *dbMock {
	rows := sqlmock.NewRows(eventColumns)
	for i := 0; i < eventCount; i++ {
		rows.AddRow(time.Now(), "eventType", Sequence(i+1), Sequence(i), Sequence(i), nil, "svc", "hodor", "org", "aggType", "aggID", "v1.0.0")
	}
	db.mock.ExpectQuery(expectedFilterEventsLimitFormat).
		WithArgs(aggregateType, limit).
		WillReturnRows(rows)
	return db
}

func (db *dbMock) expectFilterEventsDesc(aggregateType string, eventCount int) *dbMock {
	rows := sqlmock.NewRows(eventColumns)
	for i := eventCount; i > 0; i-- {
		rows.AddRow(time.Now(), "eventType", Sequence(i+1), Sequence(i), Sequence(i), nil, "svc", "hodor", "org", "aggType", "aggID", "v1.0.0")
	}
	db.mock.ExpectQuery(expectedFilterEventsDescFormat).
		WillReturnRows(rows)
	return db
}

func (db *dbMock) expectFilterEventsAggregateIDLimit(aggregateType, aggregateID string, limit uint64) *dbMock {
	rows := sqlmock.NewRows(eventColumns)
	for i := limit; i > 0; i-- {
		rows.AddRow(time.Now(), "eventType", Sequence(i+1), Sequence(i), Sequence(i), nil, "svc", "hodor", "org", "aggType", "aggID", "v1.0.0")
	}
	db.mock.ExpectQuery(expectedFilterEventsAggregateIDLimit).
		WithArgs(aggregateType, aggregateID, limit).
		WillReturnRows(rows)
	return db
}

func (db *dbMock) expectFilterEventsAggregateIDTypeLimit(aggregateType, aggregateID string, limit uint64) *dbMock {
	rows := sqlmock.NewRows(eventColumns)
	for i := limit; i > 0; i-- {
		rows.AddRow(time.Now(), "eventType", Sequence(i+1), Sequence(i), Sequence(i), nil, "svc", "hodor", "org", "aggType", "aggID", "v1.0.0")
	}
	db.mock.ExpectQuery(expectedFilterEventsAggregateIDTypeLimit).
		WithArgs(aggregateType, aggregateID, limit).
		WillReturnRows(rows)
	return db
}

func (db *dbMock) expectFilterEventsError(returnedErr error) *dbMock {
	db.mock.ExpectQuery(expectedGetAllEvents).
		WillReturnError(returnedErr)
	return db
}

func (db *dbMock) expectLatestSequenceFilter(aggregateType string, sequence Sequence) *dbMock {
	db.mock.ExpectQuery(`SELECT MAX\(event_sequence\) FROM eventstore\.events WHERE aggregate_type = \$1`).
		WithArgs(aggregateType).
		WillReturnRows(sqlmock.NewRows([]string{"max_sequence"}).AddRow(sequence))
	return db
}

func (db *dbMock) expectLatestSequenceFilterError(aggregateType string, err error) *dbMock {
	db.mock.ExpectQuery(`SELECT MAX\(event_sequence\) FROM eventstore\.events WHERE aggregate_type = \$1`).
		WithArgs(aggregateType).WillReturnError(err)
	return db
}

func (db *dbMock) expectPrepareInsert(err error) *dbMock {
	prepare := db.mock.ExpectPrepare(expectedInsertStatement)
	if err != nil {
		prepare.WillReturnError(err)
	}

	return db
}
