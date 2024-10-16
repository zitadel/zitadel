package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
	db_mock "github.com/zitadel/zitadel/internal/database/mock"
)

var (
	testNow = time.Now()
	dayNow  = testNow.Truncate(24 * time.Hour)
)

// assertPrepare checks if the prepare func executes the correct sql query and returns the correct object
// prepareFunc must be of type
// func() (sq.SelectBuilder, func(*sql.Rows) (*struct, error))
// or
// func() (sq.SelectBuilder, func(*sql.Row) (*struct, error))
// expectedObject represents the return value of scan
// sqlExpectation represents the query executed on the database
func assertPrepare(t *testing.T, prepareFunc, expectedObject interface{}, sqlExpectation sqlExpectation, isErr checkErr, prepareArgs ...reflect.Value) bool {
	t.Helper()

	client, mock, err := sqlmock.New(sqlmock.ValueConverterOption(new(db_mock.TypeConverter)))
	if err != nil {
		t.Fatalf("failed to build mock client: %v", err)
	}

	mock = sqlExpectation(mock)

	builder, scan, err := execPrepare(prepareFunc, prepareArgs)
	if err != nil {
		t.Error(err)
		return false
	}
	errCheck := func(err error) (error, bool) {
		if isErr == nil {
			if err == nil {
				return nil, true
			} else {
				return fmt.Errorf("no error expected got: %w", err), false
			}
		}
		return isErr(err)
	}
	object, ok, didScan := execScan(t, &database.DB{DB: client}, builder, scan, errCheck)
	if !ok {
		t.Error(object)
		return false
	}
	if didScan {
		if !assert.Equal(t, expectedObject, object) {
			return false
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("sql expectations not met: %v", err)
		return false
	}

	return true
}

type checkErr func(error) (err error, ok bool)

type sqlExpectation func(sqlmock.Sqlmock) sqlmock.Sqlmock

func mockQuery(stmt string, cols []string, row []driver.Value, args ...driver.Value) func(m sqlmock.Sqlmock) sqlmock.Sqlmock {
	return func(m sqlmock.Sqlmock) sqlmock.Sqlmock {
		q := m.ExpectQuery(stmt).WithArgs(args...)
		result := m.NewRows(cols)
		if len(row) > 0 {
			result.AddRow(row...)
		}
		q.WillReturnRows(result)
		return m
	}
}

func mockQueryScanErr(stmt string, cols []string, row []driver.Value, args ...driver.Value) func(m sqlmock.Sqlmock) sqlmock.Sqlmock {
	return func(m sqlmock.Sqlmock) sqlmock.Sqlmock {
		q := m.ExpectQuery(stmt).WithArgs(args...)
		result := m.NewRows(cols)
		if len(row) > 0 {
			result.AddRow(row...)
		}
		q.WillReturnRows(result)
		return m
	}
}

func mockQueries(stmt string, cols []string, rows [][]driver.Value, args ...driver.Value) func(m sqlmock.Sqlmock) sqlmock.Sqlmock {
	return func(m sqlmock.Sqlmock) sqlmock.Sqlmock {
		q := m.ExpectQuery(stmt).WithArgs(args...)
		result := m.NewRows(cols)
		count := uint64(len(rows))
		for _, row := range rows {
			if cols[len(cols)-1] == "count" {
				row = append(row, count)
			}
			result.AddRow(row...)
		}
		q.WillReturnRows(result)
		q.RowsWillBeClosed()
		return m
	}
}

func mockQueriesScanErr(stmt string, cols []string, rows [][]driver.Value, args ...driver.Value) func(m sqlmock.Sqlmock) sqlmock.Sqlmock {
	return func(m sqlmock.Sqlmock) sqlmock.Sqlmock {
		q := m.ExpectQuery(stmt).WithArgs(args...)
		result := m.NewRows(cols)
		count := uint64(len(rows))
		for _, row := range rows {
			if cols[len(cols)-1] == "count" {
				row = append(row, count)
			}
			result.AddRow(row...)
		}
		q.WillReturnRows(result)
		q.RowsWillBeClosed()
		return m
	}
}

func mockQueryErr(stmt string, err error, args ...driver.Value) func(m sqlmock.Sqlmock) sqlmock.Sqlmock {
	return func(m sqlmock.Sqlmock) sqlmock.Sqlmock {
		q := m.ExpectQuery(stmt).WithArgs(args...)
		q.WillReturnError(err)
		return m
	}
}

func execMock(t testing.TB, exp sqlExpectation, run func(db *sql.DB)) {
	db, mock, err := sqlmock.New(sqlmock.ValueConverterOption(new(db_mock.TypeConverter)))
	require.NoError(t, err)
	defer db.Close()
	mock = exp(mock)
	run(db)
	assert.NoError(t, mock.ExpectationsWereMet())
}

var (
	rowType           = reflect.TypeOf(&sql.Row{})
	rowsType          = reflect.TypeOf(&sql.Rows{})
	selectBuilderType = reflect.TypeOf(sq.SelectBuilder{})
)

func execScan(t testing.TB, client *database.DB, builder sq.SelectBuilder, scan interface{}, errCheck checkErr) (object interface{}, ok bool, didScan bool) {
	t.Helper()

	scanType := reflect.TypeOf(scan)
	err := validateScan(scanType)
	if err != nil {
		return err, false, false
	}

	stmt, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("unexpected error from sql builder: %w", err), false, false
	}

	//resultSet represents *sql.Row or *sql.Rows,
	// depending on whats assignable to the scan function
	var res []reflect.Value

	//execute sql stmt
	// if scan(*sql.Rows)...
	if scanType.In(0).AssignableTo(rowsType) {
		err = client.Query(func(rows *sql.Rows) error {
			didScan = true
			res = reflect.ValueOf(scan).Call([]reflect.Value{reflect.ValueOf(rows)})
			if err, ok := res[1].Interface().(error); ok {
				return err
			}
			return nil
		}, stmt, args...)

		// if scan(*sql.Row)...
	} else if scanType.In(0).AssignableTo(rowType) {
		err = client.QueryRow(func(r *sql.Row) error {
			if r.Err() != nil {
				return r.Err()
			}
			didScan = true
			res = reflect.ValueOf(scan).Call([]reflect.Value{reflect.ValueOf(r)})
			if err, ok := res[1].Interface().(error); ok {
				return err
			}
			return nil
		}, stmt, args...)

	} else {
		return errors.New("scan: parameter must be *sql.Row or *sql.Rows"), false, false
	}

	if err != nil {
		err, ok := errCheck(err)
		if !ok {
			t.Fatal(err)
		}
		if didScan {
			return res[0].Interface(), ok, didScan
		}
		return err, ok, didScan
	}

	//check for error
	if res[1].Interface() != nil {
		if err, ok := errCheck(res[1].Interface().(error)); !ok {
			return fmt.Errorf("scan failed: %w", err), false, didScan
		}
	}

	return res[0].Interface(), true, didScan
}

func validateScan(scanType reflect.Type) error {
	if scanType.Kind() != reflect.Func {
		return errors.New("scan is not a function")
	}
	if scanType.NumIn() != 1 {
		return fmt.Errorf("scan: invalid number of inputs: want: 1 got %d", scanType.NumIn())
	}
	if scanType.NumOut() != 2 {
		return fmt.Errorf("scan: invalid number of outputs: want: 2 got %d", scanType.NumOut())
	}
	return nil
}

func execPrepare(prepare interface{}, args []reflect.Value) (builder sq.SelectBuilder, scan interface{}, err error) {
	prepareVal := reflect.ValueOf(prepare)
	if err := validatePrepare(prepareVal.Type()); err != nil {
		return sq.SelectBuilder{}, nil, err
	}
	res := prepareVal.Call(args)

	return res[0].Interface().(sq.SelectBuilder), res[1].Interface(), nil
}

func validatePrepare(prepareType reflect.Type) error {
	if prepareType.Kind() != reflect.Func {
		return errors.New("prepare is not a function")
	}
	if prepareType.NumIn() != 0 && prepareType.NumIn() != 2 {
		return fmt.Errorf("prepare: invalid number of inputs: want: 0 or 2 got %d", prepareType.NumIn())
	}
	if prepareType.NumOut() != 2 {
		return fmt.Errorf("prepare: invalid number of outputs: want: 2 got %d", prepareType.NumOut())
	}
	if prepareType.Out(0) != selectBuilderType {
		return fmt.Errorf("prepare: first return value must be: %s got %s", selectBuilderType, prepareType.Out(0))
	}
	if prepareType.Out(1).Kind() != reflect.Func {
		return fmt.Errorf("prepare: second return value must be: %s got %s", reflect.Func, prepareType.Out(1))
	}
	return nil
}

func TestValidateScan(t *testing.T) {
	tests := []struct {
		name      string
		t         reflect.Type
		expectErr bool
	}{
		{
			name:      "not a func",
			t:         reflect.TypeOf(&struct{}{}),
			expectErr: true,
		},
		{
			name: "wong input count",
			t: reflect.TypeOf(func() (*struct{}, error) {
				log.Fatal("should not be executed")
				return nil, nil
			}),
			expectErr: true,
		},
		{
			name: "wrong output count",
			t: reflect.TypeOf(func(interface{}) error {
				log.Fatal("should not be executed")
				return nil
			}),
			expectErr: true,
		},
		{
			name: "correct",
			t: reflect.TypeOf(func(interface{}) (*struct{}, error) {
				log.Fatal("should not be executed")
				return nil, nil
			}),
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateScan(tt.t)
			if (err != nil) != tt.expectErr {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

func TestValidatePrepare(t *testing.T) {
	tests := []struct {
		name      string
		t         reflect.Type
		expectErr bool
	}{
		{
			name:      "not a func",
			t:         reflect.TypeOf(&struct{}{}),
			expectErr: true,
		},
		{
			name: "wong input count",
			t: reflect.TypeOf(func(int) (sq.SelectBuilder, func(*sql.Rows) (interface{}, error)) {
				log.Fatal("should not be executed")
				return sq.SelectBuilder{}, nil
			}),
			expectErr: true,
		},
		{
			name: "wrong output count",
			t: reflect.TypeOf(func() sq.SelectBuilder {
				log.Fatal("should not be executed")
				return sq.SelectBuilder{}
			}),
			expectErr: true,
		},
		{
			name: "first output type wrong",
			t: reflect.TypeOf(func() (*struct{}, func(*sql.Rows) (interface{}, error)) {
				log.Fatal("should not be executed")
				return nil, nil
			}),
			expectErr: true,
		},
		{
			name: "second output type wrong",
			t: reflect.TypeOf(func() (sq.SelectBuilder, *struct{}) {
				log.Fatal("should not be executed")
				return sq.SelectBuilder{}, nil
			}),
			expectErr: true,
		},
		{
			name: "correct",
			t: reflect.TypeOf(func(context.Context, prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (interface{}, error)) {
				log.Fatal("should not be executed")
				return sq.SelectBuilder{}, nil
			}),
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePrepare(tt.t)
			if (err != nil) != tt.expectErr {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}

type prepareDB struct{}

const asOfSystemTime = " AS OF SYSTEM TIME '-1 ms' "

func (*prepareDB) Timetravel(time.Duration) string { return asOfSystemTime }

var defaultPrepareArgs = []reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(new(prepareDB))}

func (*prepareDB) DatabaseName() string { return "db" }

func (*prepareDB) Username() string { return "user" }

func (*prepareDB) Type() string { return "type" }
