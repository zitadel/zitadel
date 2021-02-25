package repository

import (
	"database/sql/driver"
	"fmt"
	"github.com/caos/zitadel/internal/domain"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

var (
	expectedGetByID                 = `SELECT \* FROM "%s" WHERE \(%s = \$1\) LIMIT 1`
	expectedGetByQuery              = `SELECT \* FROM "%s" WHERE \(LOWER\(%s\) %s LOWER\(\$1\)\) LIMIT 1`
	expectedGetByQueryCaseSensitive = `SELECT \* FROM "%s" WHERE \(%s %s \$1\) LIMIT 1`
	expectedSave                    = `UPDATE "%s" SET "test" = \$1 WHERE "%s"."%s" = \$2`
	expectedRemove                  = `DELETE FROM "%s" WHERE \(%s = \$1\)`
	expectedRemoveByKeys            = func(i int, table string) string {
		sql := fmt.Sprintf(`DELETE FROM "%s"`, table)
		sql += ` WHERE \(%s = \$1\)`
		for j := 1; j < i; j++ {
			sql = sql + ` AND \(%s = \$` + strconv.Itoa(j+1) + `\)`
		}
		return sql
	}
	expectedRemoveByObject           = `DELETE FROM "%s" WHERE "%s"."%s" = \$1`
	expectedRemoveByObjectMultiplePK = `DELETE FROM "%s" WHERE "%s"."%s" = \$1 AND "%s"."%s" = \$2`
	expectedTruncate                 = `TRUNCATE %s;`
	expectedSearch                   = `SELECT \* FROM "%s" OFFSET 0`
	expectedSearchCount              = `SELECT count\(\*\) FROM "%s"`
	expectedSearchLimit              = `SELECT \* FROM "%s" LIMIT %v OFFSET 0`
	expectedSearchLimitCount         = `SELECT count\(\*\) FROM "%s"`
	expectedSearchOffset             = `SELECT \* FROM "%s" OFFSET %v`
	expectedSearchOffsetCount        = `SELECT count\(\*\) FROM "%s"`
	expectedSearchSorting            = `SELECT \* FROM "%s" ORDER BY %s %s OFFSET 0`
	expectedSearchSortingCount       = `SELECT count\(\*\) FROM "%s"`
	expectedSearchQuery              = `SELECT \* FROM "%s" WHERE \(LOWER\(%s\) %s LOWER\(\$1\)\) OFFSET 0`
	expectedSearchQueryCount         = `SELECT count\(\*\) FROM "%s" WHERE \(LOWER\(%s\) %s LOWER\(\$1\)\)`
	expectedSearchQueryAllParams     = `SELECT \* FROM "%s" WHERE \(LOWER\(%s\) %s LOWER\(\$1\)\) ORDER BY %s %s LIMIT %v OFFSET %v`
	expectedSearchQueryAllParamCount = `SELECT count\(\*\) FROM "%s" WHERE \(LOWER\(%s\) %s LOWER\(\$1\)\)`
)

type TestSearchRequest struct {
	limit         uint64
	offset        uint64
	sortingColumn ColumnKey
	asc           bool
	queries       []SearchQuery
}

func (req TestSearchRequest) GetLimit() uint64 {
	return req.limit
}

func (req TestSearchRequest) GetOffset() uint64 {
	return req.offset
}

func (req TestSearchRequest) GetSortingColumn() ColumnKey {
	return req.sortingColumn
}

func (req TestSearchRequest) GetAsc() bool {
	return req.asc
}

func (req TestSearchRequest) GetQueries() []SearchQuery {
	return req.queries
}

type TestSearchQuery struct {
	key    TestSearchKey
	method domain.SearchMethod
	value  string
}

func (req TestSearchQuery) GetKey() ColumnKey {
	return req.key
}

func (req TestSearchQuery) GetMethod() domain.SearchMethod {
	return req.method
}

func (req TestSearchQuery) GetValue() interface{} {
	return req.value
}

type TestSearchKey int32

const (
	TestSearchKey_UNDEFINED TestSearchKey = iota
	TestSearchKey_TEST
	TestSearchKey_ID
)

func (key TestSearchKey) ToColumnName() string {
	switch TestSearchKey(key) {
	case TestSearchKey_TEST:
		return "test"
	case TestSearchKey_ID:
		return "id"
	default:
		return ""
	}
}

type Test struct {
	ID   string `json:"-" gorm:"column:primary_id;primary_key"`
	Test string `json:"test" gorm:"column:test"`
}

type TestMultiplePK struct {
	TestID  string `gorm:"column:testId;primary_key"`
	HodorID string `gorm:"column:hodorId;primary_key"`
	Test    string `gorm:"column:test"`
}

type dbMock struct {
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

func (db *dbMock) close() {
	db.db.Close()
}

func mockDB(t *testing.T) *dbMock {
	mockDB := dbMock{}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error occured while creating stub db %v", err)
	}

	mockDB.mock = mock
	mockDB.db, err = gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("error occured while connecting to stub db: %v", err)
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

func (db *dbMock) expectGetByID(table, key, value string) *dbMock {
	query := fmt.Sprintf(expectedGetByID, table, key)
	db.mock.ExpectQuery(query).
		WithArgs(value).
		WillReturnRows(sqlmock.NewRows([]string{key}).
			AddRow(key))

	return db
}

func (db *dbMock) expectGetByIDErr(table, key, value string, err error) *dbMock {
	query := fmt.Sprintf(expectedGetByID, table, key)
	db.mock.ExpectQuery(query).
		WithArgs(value).
		WillReturnError(err)

	return db
}

func (db *dbMock) expectGetByQuery(table, key, method, value string) *dbMock {
	query := fmt.Sprintf(expectedGetByQuery, table, key, method)
	db.mock.ExpectQuery(query).
		WithArgs(value).
		WillReturnRows(sqlmock.NewRows([]string{key}).
			AddRow(key))

	return db
}

func (db *dbMock) expectGetByQueryCaseSensitive(table, key, method, value string) *dbMock {
	query := fmt.Sprintf(expectedGetByQueryCaseSensitive, table, key, method)
	db.mock.ExpectQuery(query).
		WithArgs(value).
		WillReturnRows(sqlmock.NewRows([]string{key}).
			AddRow(key))

	return db
}

func (db *dbMock) expectGetByQueryErr(table, key, method, value string, err error) *dbMock {
	query := fmt.Sprintf(expectedGetByQuery, table, key, method)
	db.mock.ExpectQuery(query).
		WithArgs(value).
		WillReturnError(err)

	return db
}

func (db *dbMock) expectSave(table string, object Test) *dbMock {
	query := fmt.Sprintf(expectedSave, table, table, "primary_id")
	db.mock.ExpectExec(query).
		WithArgs(object.Test, object.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	return db
}

func (db *dbMock) expectSaveErr(table string, object Test, err error) *dbMock {
	query := fmt.Sprintf(expectedSave, table, table, "id")
	db.mock.ExpectExec(query).
		WithArgs(object.Test, object.ID).
		WillReturnError(err)

	return db
}

func (db *dbMock) expectRemove(table, key, value string) *dbMock {
	query := fmt.Sprintf(expectedRemove, table, key)
	db.mock.ExpectExec(query).
		WithArgs(value).
		WillReturnResult(sqlmock.NewResult(1, 1))

	return db
}

func (db *dbMock) expectRemoveKeys(table string, keys ...Key) *dbMock {
	keynames := make([]interface{}, len(keys))
	keyvalues := make([]driver.Value, len(keys))
	for i, key := range keys {
		keynames[i] = key.Key.ToColumnName()
		keyvalues[i] = key.Value
	}
	query := fmt.Sprintf(expectedRemoveByKeys(len(keys), table), keynames...)
	db.mock.ExpectExec(query).
		WithArgs(keyvalues...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	return db
}

func (db *dbMock) expectRemoveByObject(table string, object Test) *dbMock {
	query := fmt.Sprintf(expectedRemoveByObject, table, table, "primary_id")
	db.mock.ExpectExec(query).
		WithArgs(object.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	return db
}

func (db *dbMock) expectRemoveByObjectMultiplePKs(table string, object TestMultiplePK) *dbMock {
	query := fmt.Sprintf(expectedRemoveByObjectMultiplePK, table, table, "testId", table, "hodorId")
	db.mock.ExpectExec(query).
		WithArgs(object.TestID, object.HodorID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	return db
}

func (db *dbMock) expectRemoveErr(table, key, value string, err error) *dbMock {
	query := fmt.Sprintf(expectedRemove, table, key)
	db.mock.ExpectExec(query).
		WithArgs(value).
		WillReturnError(err)

	return db
}

func (db *dbMock) expectTruncate(table string) *dbMock {
	query := fmt.Sprintf(expectedTruncate, table)
	db.mock.ExpectExec(query).
		WillReturnResult(sqlmock.NewResult(1, 1))

	return db
}
func (db *dbMock) expectTruncateErr(table string, err error) *dbMock {
	query := fmt.Sprintf(expectedTruncate, table)
	db.mock.ExpectExec(query).
		WillReturnError(err)

	return db
}
func (db *dbMock) expectGetSearchRequestNoParams(table string, resultAmount, total int) *dbMock {
	query := fmt.Sprintf(expectedSearch, table)
	queryCount := fmt.Sprintf(expectedSearchCount, table)

	rows := sqlmock.NewRows([]string{"id"})
	for i := 0; i < resultAmount; i++ {
		rows.AddRow(fmt.Sprintf("hodor-%d", i))
	}

	db.mock.ExpectQuery(queryCount).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(total))
	db.mock.ExpectQuery(query).
		WillReturnRows(rows)
	return db
}

func (db *dbMock) expectGetSearchRequestWithLimit(table string, limit, resultAmount, total int) *dbMock {
	query := fmt.Sprintf(expectedSearchLimit, table, limit)
	queryCount := fmt.Sprintf(expectedSearchLimitCount, table)

	rows := sqlmock.NewRows([]string{"id"})
	for i := 0; i < resultAmount; i++ {
		rows.AddRow(fmt.Sprintf("hodor-%d", i))
	}

	db.mock.ExpectQuery(queryCount).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(total))
	db.mock.ExpectQuery(query).
		WillReturnRows(rows)
	return db
}

func (db *dbMock) expectGetSearchRequestWithOffset(table string, offset, resultAmount, total int) *dbMock {
	query := fmt.Sprintf(expectedSearchOffset, table, offset)
	queryCount := fmt.Sprintf(expectedSearchOffsetCount, table)

	rows := sqlmock.NewRows([]string{"id"})
	for i := 0; i < resultAmount; i++ {
		rows.AddRow(fmt.Sprintf("hodor-%d", i))
	}

	db.mock.ExpectQuery(queryCount).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(total))
	db.mock.ExpectQuery(query).
		WillReturnRows(rows)
	return db
}

func (db *dbMock) expectGetSearchRequestWithSorting(table, sorting string, sortingColumn ColumnKey, resultAmount, total int) *dbMock {
	query := fmt.Sprintf(expectedSearchSorting, table, sortingColumn.ToColumnName(), sorting)
	queryCount := fmt.Sprintf(expectedSearchSortingCount, table)

	rows := sqlmock.NewRows([]string{"id"})
	for i := 0; i < resultAmount; i++ {
		rows.AddRow(fmt.Sprintf("hodor-%d", i))
	}

	db.mock.ExpectQuery(queryCount).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(total))
	db.mock.ExpectQuery(query).
		WillReturnRows(rows)
	return db
}

func (db *dbMock) expectGetSearchRequestWithSearchQuery(table, key, method, value string, resultAmount, total int) *dbMock {
	query := fmt.Sprintf(expectedSearchQuery, table, key, method)
	queryCount := fmt.Sprintf(expectedSearchQueryCount, table, key, method)

	rows := sqlmock.NewRows([]string{"id"})
	for i := 0; i < resultAmount; i++ {
		rows.AddRow(fmt.Sprintf("hodor-%d", i))
	}

	db.mock.ExpectQuery(queryCount).
		WithArgs(value).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(total))
	db.mock.ExpectQuery(query).
		WithArgs(value).
		WillReturnRows(rows)
	return db
}

func (db *dbMock) expectGetSearchRequestWithAllParams(table, key, method, value, sorting string, sortingColumn ColumnKey, limit, offset, resultAmount, total int) *dbMock {
	query := fmt.Sprintf(expectedSearchQueryAllParams, table, key, method, sortingColumn.ToColumnName(), sorting, limit, offset)
	queryCount := fmt.Sprintf(expectedSearchQueryAllParamCount, table, key, method)

	rows := sqlmock.NewRows([]string{"id"})
	for i := 0; i < resultAmount; i++ {
		rows.AddRow(fmt.Sprintf("hodor-%d", i))
	}

	db.mock.ExpectQuery(queryCount).
		WithArgs(value).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(total))
	db.mock.ExpectQuery(query).
		WithArgs(value).
		WillReturnRows(rows)
	return db
}

func (db *dbMock) expectGetSearchRequestErr(table string, resultAmount, total int, err error) *dbMock {
	query := fmt.Sprintf(expectedSearch, table)
	queryCount := fmt.Sprintf(expectedSearchCount, table)

	rows := sqlmock.NewRows([]string{"id"})
	for i := 0; i < resultAmount; i++ {
		rows.AddRow(fmt.Sprintf("hodor-%d", i))
	}

	db.mock.ExpectQuery(queryCount).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(total))
	db.mock.ExpectQuery(query).
		WillReturnError(err)
	return db
}
