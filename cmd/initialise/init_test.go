package initialise

import (
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/zitadel/zitadel/internal/database"
	db_mock "github.com/zitadel/zitadel/internal/database/mock"
)

type db struct {
	mock sqlmock.Sqlmock
	db   *database.DB
}

func prepareDB(t *testing.T, expectations ...expectation) db {
	t.Helper()
	client, mock, err := sqlmock.New(sqlmock.ValueConverterOption(new(db_mock.TypeConverter)))
	if err != nil {
		t.Fatalf("unable to create sql mock: %v", err)
	}
	for _, expectation := range expectations {
		expectation(mock)
	}
	return db{
		mock: mock,
		db:   &database.DB{DB: client},
	}
}

type expectation func(m sqlmock.Sqlmock)

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

func expectQuery(stmt string, err error, columns []string, rows [][]driver.Value, args ...driver.Value) expectation {
	return func(m sqlmock.Sqlmock) {
		res := m.NewRows(columns)
		for _, row := range rows {
			res.AddRow(row...)
		}
		query := m.ExpectQuery(regexp.QuoteMeta(stmt)).WithArgs(args...).WillReturnRows(res)
		if err != nil {
			query.WillReturnError(err)
			return
		}
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
