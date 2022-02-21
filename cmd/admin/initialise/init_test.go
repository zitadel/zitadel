package initialise

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type db struct {
	mock sqlmock.Sqlmock
	db   *sql.DB
}

func prepareDB(t *testing.T, expectations ...expectation) db {
	t.Helper()
	client, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unable to create sql mock: %v", err)
	}
	for _, expectation := range expectations {
		expectation(mock)
	}
	return db{
		mock: mock,
		db:   client,
	}
}

type expectation func(m sqlmock.Sqlmock)

func expectExists(query string, value bool, args ...driver.Value) expectation {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(value))
	}
}

func expectQueryErr(query string, err error, args ...driver.Value) expectation {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnError(err)
	}
}

func expectExec(stmt string, err error, args ...driver.Value) expectation {
	return func(m sqlmock.Sqlmock) {
		query := m.ExpectExec(regexp.QuoteMeta(stmt)).WithArgs(args...)
		if err != nil {
			query.WillReturnError(err)
			return
		}
		query.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func expectBegin(err error) expectation {
	return func(m sqlmock.Sqlmock) {
		query := m.ExpectBegin()
		if err != nil {
			query.WillReturnError(err)
		}
	}
}

func expectCommit(err error) expectation {
	return func(m sqlmock.Sqlmock) {
		query := m.ExpectCommit()
		if err != nil {
			query.WillReturnError(err)
		}
	}
}

func expectRollback(err error) expectation {
	return func(m sqlmock.Sqlmock) {
		query := m.ExpectRollback()
		if err != nil {
			query.WillReturnError(err)
		}
	}
}
