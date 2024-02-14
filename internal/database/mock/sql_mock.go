package mock

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type SQLMock struct {
	DB   *sql.DB
	mock sqlmock.Sqlmock
}

type expectation func(m sqlmock.Sqlmock)

func NewSQLMock(t *testing.T, expectations ...expectation) *SQLMock {
	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		t.Fatal("create mock failed", err)
	}

	for _, expectation := range expectations {
		expectation(mock)
	}

	return &SQLMock{
		DB:   db,
		mock: mock,
	}
}

func (m *SQLMock) Assert(t *testing.T) {
	t.Helper()

	if err := m.mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %v", err)
	}

	m.DB.Close()
}

func ExpectBegin(err error) expectation {
	return func(m sqlmock.Sqlmock) {
		e := m.ExpectBegin()
		if err != nil {
			e.WillReturnError(err)
		}
	}
}

func ExpectCommit(err error) expectation {
	return func(m sqlmock.Sqlmock) {
		e := m.ExpectCommit()
		if err != nil {
			e.WillReturnError(err)
		}
	}
}

type ExecOpt func(e *sqlmock.ExpectedExec) *sqlmock.ExpectedExec

func WithExecArgs(args ...driver.Value) ExecOpt {
	return func(e *sqlmock.ExpectedExec) *sqlmock.ExpectedExec {
		return e.WithArgs(args...)
	}
}

func WithExecErr(err error) ExecOpt {
	return func(e *sqlmock.ExpectedExec) *sqlmock.ExpectedExec {
		return e.WillReturnError(err)
	}
}

func WithExecNoRowsAffected() ExecOpt {
	return func(e *sqlmock.ExpectedExec) *sqlmock.ExpectedExec {
		return e.WillReturnResult(driver.ResultNoRows)
	}
}

func WithExecRowsAffected(affected driver.RowsAffected) ExecOpt {
	return func(e *sqlmock.ExpectedExec) *sqlmock.ExpectedExec {
		return e.WillReturnResult(affected)
	}
}

func ExcpectExec(stmt string, opts ...ExecOpt) expectation {
	return func(m sqlmock.Sqlmock) {
		e := m.ExpectExec(stmt)
		for _, opt := range opts {
			e = opt(e)
		}
	}
}

type QueryOpt func(e *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery

func WithQueryArgs(args ...driver.Value) QueryOpt {
	return func(e *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
		return e.WithArgs(args...)
	}
}

func WithQueryErr(err error) QueryOpt {
	return func(e *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
		return e.WillReturnError(err)
	}
}

func WithQueryResult(columns []string, rows [][]driver.Value) QueryOpt {
	return func(e *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
		mockedRows := sqlmock.NewRows(columns)
		for _, row := range rows {
			mockedRows = mockedRows.AddRow(row...)
		}
		return e.WillReturnRows(mockedRows)
	}
}

func ExpectQuery(stmt string, opts ...QueryOpt) expectation {
	return func(m sqlmock.Sqlmock) {
		e := m.ExpectQuery(stmt)
		for _, opt := range opts {
			e = opt(e)
		}
	}
}

type AnyType[T interface{}] struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyType[T]) Match(v driver.Value) bool {
	return reflect.TypeOf(new(T)).Elem().Kind().String() == reflect.TypeOf(v).Kind().String()
}
