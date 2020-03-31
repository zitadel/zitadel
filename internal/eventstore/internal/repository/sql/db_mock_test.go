package sql

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/lib/pq"
)

var (
	expectedFilterEventsLimitFormat          = regexp.MustCompile(`SELECT id, creation_date, event_type, event_sequence, previous_sequence, event_data, modifier_service, modifier_tenant, modifier_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore\.events ORDER BY event_sequence LIMIT \$1`).String()
	expectedFilterEventsDescFormat           = regexp.MustCompile(`SELECT id, creation_date, event_type, event_sequence, previous_sequence, event_data, modifier_service, modifier_tenant, modifier_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore\.events ORDER BY event_sequence DESC`).String()
	expectedFilterEventsAggregateIDLimit     = regexp.MustCompile(`SELECT id, creation_date, event_type, event_sequence, previous_sequence, event_data, modifier_service, modifier_tenant, modifier_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore\.events WHERE aggregate_id = \$1 ORDER BY event_sequence LIMIT \$2`).String()
	expectedFilterEventsAggregateIDTypeLimit = regexp.MustCompile(`SELECT id, creation_date, event_type, event_sequence, previous_sequence, event_data, modifier_service, modifier_tenant, modifier_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore\.events WHERE aggregate_id = \$1 AND aggregate_type IN \(\$2\) ORDER BY event_sequence LIMIT \$3`).String()
	expectedGetAllEvents                     = regexp.MustCompile(`SELECT id, creation_date, event_type, event_sequence, previous_sequence, event_data, modifier_service, modifier_tenant, modifier_user, resource_owner, aggregate_type, aggregate_id, aggregate_version FROM eventstore\.events ORDER BY event_sequence`).String()

	expectedInsertStatement = regexp.MustCompile(`insert into eventstore\.events ` +
		`\(event_type, aggregate_type, aggregate_id, aggregate_version, creation_date, event_data, modifier_user, modifier_service, modifier_tenant, resource_owner, previous_sequence\) ` +
		`select \$1, \$2, \$3, \$4, coalesce\(\$5, now\(\)\), \$6, \$7, \$8, \$9, \$10, ` +
		`case \(select exists\(select event_sequence from eventstore\.events where aggregate_type = \$11 AND aggregate_id = \$12\)\) ` +
		`WHEN true then \(select event_sequence from eventstore\.events where aggregate_type = \$13 AND aggregate_id = \$14 order by event_sequence desc limit 1\) ` +
		`ELSE NULL ` +
		`end ` +
		`where \(` +
		`\(select count\(id\) from eventstore\.events where event_sequence >= \$15 AND aggregate_type = \$16 AND aggregate_id = \$17\) = 1 OR ` +
		`\(\(select count\(id\) from eventstore\.events where aggregate_type = \$18 and aggregate_id = \$19\) = 0 AND \$20 = 0\)\) RETURNING id, event_sequence, creation_date`).String()
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

func (db *dbMock) expectInsertEvent(e *models.Event, returnedID string, returnedSequence uint64) *dbMock {
	db.mock.ExpectQuery(expectedInsertStatement).
		WithArgs(
			e.Typ, e.AggregateType, e.AggregateID, e.AggregateVersion, sqlmock.AnyArg(), e.Data, e.ModifierUser, e.ModifierService, e.ModifierTenant, e.ResourceOwner,
			e.AggregateType, e.AggregateID,
			e.AggregateType, e.AggregateID,
			e.PreviousSequence, e.AggregateType, e.AggregateID,
			e.AggregateType, e.AggregateID, e.PreviousSequence,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "event_sequence", "creation_date"}).
				AddRow(returnedID, returnedSequence, time.Now().UTC()),
		)

	return db
}

func (db *dbMock) expectInsertEventError(e *models.Event) *dbMock {
	db.mock.ExpectQuery(expectedInsertStatement).
		WithArgs(
			e.Typ, e.AggregateType, e.AggregateID, e.AggregateVersion, sqlmock.AnyArg(), e.Data, e.ModifierUser, e.ModifierService, e.ModifierTenant, e.ResourceOwner,
			e.AggregateType, e.AggregateID,
			e.AggregateType, e.AggregateID,
			e.PreviousSequence, e.AggregateType, e.AggregateID,
			e.AggregateType, e.AggregateID, e.PreviousSequence,
		).
		WillReturnError(sql.ErrTxDone)

	return db
}

func (db *dbMock) expectFilterEventsLimit(limit uint64, eventCount int) *dbMock {
	rows := sqlmock.NewRows([]string{"id", "creation_date"})
	for i := 0; i < eventCount; i++ {
		rows.AddRow(fmt.Sprint("event", i), time.Now())
	}
	db.mock.ExpectQuery(expectedFilterEventsLimitFormat).
		WithArgs(limit).
		WillReturnRows(rows)
	return db
}

func (db *dbMock) expectFilterEventsDesc(eventCount int) *dbMock {
	rows := sqlmock.NewRows([]string{"id", "creation_date"})
	for i := eventCount; i > 0; i-- {
		rows.AddRow(fmt.Sprint("event", i), time.Now())
	}
	db.mock.ExpectQuery(expectedFilterEventsDescFormat).
		WillReturnRows(rows)
	return db
}

func (db *dbMock) expectFilterEventsAggregateIDLimit(aggregateID string, limit uint64) *dbMock {
	rows := sqlmock.NewRows([]string{"id", "creation_date"})
	for i := limit; i > 0; i-- {
		rows.AddRow(fmt.Sprint("event", i), time.Now())
	}
	db.mock.ExpectQuery(expectedFilterEventsAggregateIDLimit).
		WithArgs(aggregateID, limit).
		WillReturnRows(rows)
	return db
}

func (db *dbMock) expectFilterEventsAggregateIDTypeLimit(aggregateID, aggregateType string, limit uint64) *dbMock {
	rows := sqlmock.NewRows([]string{"id", "creation_date"})
	for i := limit; i > 0; i-- {
		rows.AddRow(fmt.Sprint("event", i), time.Now())
	}
	db.mock.ExpectQuery(expectedFilterEventsAggregateIDTypeLimit).
		WithArgs(aggregateID, pq.Array([]string{aggregateType}), limit).
		WillReturnRows(rows)
	return db
}

func (db *dbMock) expectFilterEventsError(returnedErr error) *dbMock {
	db.mock.ExpectQuery(expectedGetAllEvents).
		WillReturnError(returnedErr)
	return db
}

func (db *dbMock) expectPrepareInsert() *dbMock {
	db.mock.ExpectPrepare(expectedInsertStatement)

	return db
}
