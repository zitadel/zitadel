package repository

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/jinzhu/gorm"
	"testing"
)

func TestPrepareGetByKey(t *testing.T) {
	type args struct {
		table string
		key   ColumnKey
		value string
	}
	type res struct {
		result  Test
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		db   *dbMock
		args args
		res  res
	}{
		{
			"ok",
			mockDB(t).
				expectGetByID("TESTTABLE", "test", "VALUE"),
			args{
				table: "TESTTABLE",
				key:   TestSearchKey_TEST,
				value: "VALUE",
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: false,
			},
		},
		{
			"not found",
			mockDB(t).
				expectGetByIDErr("TESTTABLE", "test", "VALUE", gorm.ErrRecordNotFound),
			args{
				table: "TESTTABLE",
				key:   TestSearchKey_TEST,
				value: "VALUE",
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			"db err",
			mockDB(t).
				expectGetByIDErr("TESTTABLE", "test", "VALUE", gorm.ErrUnaddressable),
			args{
				table: "TESTTABLE",
				key:   TestSearchKey_TEST,
				value: "VALUE",
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: true,
				errFunc: caos_errs.IsInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &Test{}
			getByID := PrepareGetByKey(tt.args.table, tt.args.key, tt.args.value)
			err := getByID(tt.db.db, res)

			if !tt.res.wantErr && err != nil {
				t.Errorf("got wrong err should be nil: %v ", err)
			}

			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if err := tt.db.mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			tt.db.close()
		})
	}
}

func TestPrepareGetByQuery(t *testing.T) {
	type args struct {
		table       string
		searchQuery SearchQuery
	}
	type res struct {
		result  Test
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		db   *dbMock
		args args
		res  res
	}{
		{
			"search with equals case insensitive",
			mockDB(t).
				expectGetByQuery("TESTTABLE", "test", "=", "VALUE"),
			args{
				table:       "TESTTABLE",
				searchQuery: TestSearchQuery{key: TestSearchKey_TEST, method: domain.SearchMethodEqualsIgnoreCase, value: "VALUE"},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: false,
			},
		},
		{
			"search with equals case sensitive",
			mockDB(t).
				expectGetByQueryCaseSensitive("TESTTABLE", "test", "=", "VALUE"),
			args{
				table:       "TESTTABLE",
				searchQuery: TestSearchQuery{key: TestSearchKey_TEST, method: domain.SearchMethodEquals, value: "VALUE"},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: false,
			},
		},
		{
			"search with startswith, case insensitive",
			mockDB(t).
				expectGetByQuery("TESTTABLE", "test", "LIKE", "VALUE%"),
			args{
				table:       "TESTTABLE",
				searchQuery: TestSearchQuery{key: TestSearchKey_TEST, method: domain.SearchMethodStartsWithIgnoreCase, value: "VALUE"},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: false,
			},
		},
		{
			"search with startswith case sensitive",
			mockDB(t).
				expectGetByQueryCaseSensitive("TESTTABLE", "test", "LIKE", "VALUE%"),
			args{
				table:       "TESTTABLE",
				searchQuery: TestSearchQuery{key: TestSearchKey_TEST, method: domain.SearchMethodStartsWith, value: "VALUE"},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: false,
			},
		},
		{
			"search with contains case insensitive",
			mockDB(t).
				expectGetByQuery("TESTTABLE", "test", "LIKE", "%VALUE%"),
			args{
				table:       "TESTTABLE",
				searchQuery: TestSearchQuery{key: TestSearchKey_TEST, method: domain.SearchMethodContainsIgnoreCase, value: "VALUE"},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: false,
			},
		},
		{
			"search with contains case sensitive",
			mockDB(t).
				expectGetByQueryCaseSensitive("TESTTABLE", "test", "LIKE", "%VALUE%"),
			args{
				table:       "TESTTABLE",
				searchQuery: TestSearchQuery{key: TestSearchKey_TEST, method: domain.SearchMethodContains, value: "VALUE"},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: false,
			},
		},
		{
			"search expect not found err",
			mockDB(t).
				expectGetByQueryErr("TESTTABLE", "test", "LIKE", "%VALUE%", gorm.ErrRecordNotFound),
			args{
				table:       "TESTTABLE",
				searchQuery: TestSearchQuery{key: TestSearchKey_TEST, method: domain.SearchMethodContainsIgnoreCase, value: "VALUE"},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			"search expect internal err",
			mockDB(t).
				expectGetByQueryErr("TESTTABLE", "test", "LIKE", "%VALUE%", gorm.ErrUnaddressable),
			args{
				table:       "TESTTABLE",
				searchQuery: TestSearchQuery{key: TestSearchKey_TEST, method: domain.SearchMethodContainsIgnoreCase, value: "VALUE"},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: true,
				errFunc: caos_errs.IsInternal,
			},
		},
		{
			"search with invalid column",
			mockDB(t).
				expectGetByQuery("TESTTABLE", "", "=", "VALUE"),
			args{
				table:       "TESTTABLE",
				searchQuery: TestSearchQuery{key: TestSearchKey_UNDEFINED, method: domain.SearchMethodEqualsIgnoreCase, value: "VALUE"},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: true,
				errFunc: caos_errs.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &Test{}
			getByQuery := PrepareGetByQuery(tt.args.table, tt.args.searchQuery)
			err := getByQuery(tt.db.db, res)

			if !tt.res.wantErr && err != nil {
				t.Errorf("got wrong err should be nil: %v ", err)
			}

			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if err := tt.db.mock.ExpectationsWereMet(); !tt.res.wantErr && err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			tt.db.close()
		})
	}
}

func TestPreparePut(t *testing.T) {
	type args struct {
		table  string
		object *Test
	}
	type res struct {
		result  Test
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		db   *dbMock
		args args
		res  res
	}{
		{
			"ok",
			mockDB(t).
				expectBegin(nil).
				expectSave("TESTTABLE", Test{ID: "AggregateID", Test: "VALUE"}).
				expectCommit(nil),
			args{
				table:  "TESTTABLE",
				object: &Test{ID: "AggregateID", Test: "VALUE"},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: false,
			},
		},
		{
			"db error",
			mockDB(t).
				expectBegin(nil).
				expectSaveErr("TESTTABLE", Test{ID: "AggregateID", Test: "VALUE"}, gorm.ErrUnaddressable).
				expectCommit(nil),
			args{
				table:  "TESTTABLE",
				object: &Test{ID: "AggregateID", Test: "VALUE"},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: true,
				errFunc: caos_errs.IsInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getPut := PrepareSave(tt.args.table)
			err := getPut(tt.db.db, tt.args.object)

			if !tt.res.wantErr && err != nil {
				t.Errorf("got wrong err should be nil: %v ", err)
			}

			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if err := tt.db.mock.ExpectationsWereMet(); !tt.res.wantErr && err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			tt.db.close()
		})
	}
}

func TestPrepareDelete(t *testing.T) {
	type args struct {
		table string
		key   ColumnKey
		value string
	}
	type res struct {
		result  Test
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		db   *dbMock
		args args
		res  res
	}{
		{
			"delete",
			mockDB(t).
				expectBegin(nil).
				expectRemove("TESTTABLE", "id", "VALUE").
				expectCommit(nil),
			args{
				table: "TESTTABLE",
				key:   TestSearchKey_ID,
				value: "VALUE",
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: false,
			},
		},
		{
			"db error",
			mockDB(t).
				expectBegin(nil).
				expectRemoveErr("TESTTABLE", "id", "VALUE", gorm.ErrUnaddressable).
				expectCommit(nil),
			args{
				table: "TESTTABLE",
				key:   TestSearchKey_ID,
				value: "VALUE",
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: true,
				errFunc: caos_errs.IsInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getDelete := PrepareDeleteByKey(tt.args.table, tt.args.key, tt.args.value)
			err := getDelete(tt.db.db)

			if !tt.res.wantErr && err != nil {
				t.Errorf("got wrong err should be nil: %v ", err)
			}

			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if err := tt.db.mock.ExpectationsWereMet(); !tt.res.wantErr && err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			tt.db.close()
		})
	}
}

func TestPrepareDeleteByKeys(t *testing.T) {
	type args struct {
		table string
		keys  []Key
	}
	type res struct {
		result  Test
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		db   *dbMock
		args args
		res  res
	}{
		{
			"delete single key",
			mockDB(t).
				expectBegin(nil).
				expectRemoveKeys("TESTTABLE", Key{Key: TestSearchKey_ID, Value: "VALUE"}).
				expectCommit(nil),
			args{
				table: "TESTTABLE",
				keys: []Key{
					{Key: TestSearchKey_ID, Value: "VALUE"},
				},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: false,
			},
		},
		{
			"delete multiple keys",
			mockDB(t).
				expectBegin(nil).
				expectRemoveKeys("TESTTABLE", Key{Key: TestSearchKey_ID, Value: "VALUE"}, Key{Key: TestSearchKey_TEST, Value: "VALUE2"}).
				expectCommit(nil),
			args{
				table: "TESTTABLE",
				keys: []Key{
					{Key: TestSearchKey_ID, Value: "VALUE"},
					{Key: TestSearchKey_TEST, Value: "VALUE2"},
				},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: false,
			},
		},
		{
			"db error",
			mockDB(t).
				expectBegin(nil).
				expectRemoveErr("TESTTABLE", "id", "VALUE", gorm.ErrUnaddressable).
				expectCommit(nil),
			args{
				table: "TESTTABLE",
				keys: []Key{
					{Key: TestSearchKey_ID, Value: "VALUE"},
				},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: true,
				errFunc: caos_errs.IsInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getDelete := PrepareDeleteByKeys(tt.args.table, tt.args.keys...)
			err := getDelete(tt.db.db)

			if !tt.res.wantErr && err != nil {
				t.Errorf("got wrong err should be nil: %v ", err)
			}

			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if err := tt.db.mock.ExpectationsWereMet(); !tt.res.wantErr && err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			tt.db.close()
		})
	}
}

func TestPrepareDeleteByObject(t *testing.T) {
	type args struct {
		table  string
		object interface{}
	}
	type res struct {
		result  Test
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		db   *dbMock
		args args
		res  res
	}{
		{
			"delete",
			mockDB(t).
				expectBegin(nil).
				expectRemoveByObject("TESTTABLE", Test{ID: "VALUE", Test: "TEST"}).
				expectCommit(nil),
			args{
				table:  "TESTTABLE",
				object: &Test{ID: "VALUE", Test: "TEST"},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: false,
			},
		},
		{
			"delete multiple PK",
			mockDB(t).
				expectBegin(nil).
				expectRemoveByObjectMultiplePKs("TESTTABLE", TestMultiplePK{TestID: "TESTID", HodorID: "HODORID", Test: "TEST"}).
				expectCommit(nil),
			args{
				table:  "TESTTABLE",
				object: &TestMultiplePK{TestID: "TESTID", HodorID: "HODORID", Test: "TEST"},
			},
			res{
				wantErr: false,
			},
		},
		{
			"db error",
			mockDB(t).
				expectBegin(nil).
				expectRemoveErr("TESTTABLE", "id", "VALUE", gorm.ErrUnaddressable).
				expectCommit(nil),
			args{
				table:  "TESTTABLE",
				object: &Test{ID: "VALUE", Test: "TEST"},
			},
			res{
				result:  Test{ID: "VALUE"},
				wantErr: true,
				errFunc: caos_errs.IsInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getDelete := PrepareDeleteByObject(tt.args.table, tt.args.object)
			err := getDelete(tt.db.db)

			if !tt.res.wantErr && err != nil {
				t.Errorf("got wrong err should be nil: %v ", err)
			}

			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if err := tt.db.mock.ExpectationsWereMet(); !tt.res.wantErr && err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			tt.db.close()
		})
	}
}
