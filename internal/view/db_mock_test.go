package view

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/caos/zitadel/internal/model"
	"github.com/jinzhu/gorm"
	"testing"
)

var (
	expectedGetByID    = `SELECT \* FROM "%s" WHERE \(%s = \$1\) LIMIT 1`
	expectedGetByQuery = `SELECT \* FROM "%s" WHERE \(LOWER\(%s\) %s LOWER\(\$1\)\) LIMIT 1`
	expectedSave       = `UPDATE "%s" SET "test" = \$1 WHERE "%s"."%s" = \$2`
	expectedRemove     = `DELETE FROM "%s" WHERE \(%s = \$1\)`
)

type TestSearchQuery struct {
	key    TestSearchKey
	method model.SearchMethod
	value  string
}

func (req TestSearchQuery) GetKey() ColumnKey {
	return req.key
}

func (req TestSearchQuery) GetMethod() model.SearchMethod {
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
	id   string `json:"-" gorm:"column:id;primary_key"`
	test string `json:"test" gorm:"column:test"`
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

func (db *dbMock) expectGetByQueryErr(table, key, method, value string, err error) *dbMock {
	query := fmt.Sprintf(expectedGetByQuery, table, key, method)
	db.mock.ExpectQuery(query).
		WithArgs(value).
		WillReturnError(err)

	return db
}

func (db *dbMock) expectSave(table string, object Test) *dbMock {
	query := fmt.Sprintf(expectedSave, table, table, "id")
	db.mock.ExpectExec(query).
		WithArgs(object.test, object.id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	return db
}

func (db *dbMock) expectSaveErr(table string, object Test, err error) *dbMock {
	query := fmt.Sprintf(expectedSave, table, table, "id")
	db.mock.ExpectExec(query).
		WithArgs(object.test, object.id).
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

func (db *dbMock) expectRemoveErr(table, key, value string, err error) *dbMock {
	query := fmt.Sprintf(expectedRemove, table, key)
	db.mock.ExpectExec(query).
		WithArgs(value).
		WillReturnError(err)

	return db
}
