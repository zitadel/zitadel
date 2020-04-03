package view

import (
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/model"
	"github.com/jinzhu/gorm"
	"testing"
)

func TestPrepareGetByID(t *testing.T) {
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
			"get by id",
			mockDB(t).
				expectGetByID("TESTTABLE", "test", "VALUE"),
			args{
				table: "TESTTABLE",
				key:   TestSearchKey_TEST,
				value: "VALUE",
			},
			res{
				result:  Test{id: "VALUE"},
				wantErr: false,
			},
		},
		{
			"get by id not found",
			mockDB(t).
				expectGetByIDErr("TESTTABLE", "test", "VALUE", gorm.ErrRecordNotFound),
			args{
				table: "TESTTABLE",
				key:   TestSearchKey_TEST,
				value: "VALUE",
			},
			res{
				result:  Test{id: "VALUE"},
				wantErr: true,
				errFunc: func(err error) bool {
					return caos_errs.IsNotFound(err)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &Test{}
			getByID := PrepareGetByID(tt.args.table, tt.args.key, tt.args.value)
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
			"search with equals",
			mockDB(t).
				expectGetByQuery("TESTTABLE", "test", "=", "VALUE"),
			args{
				table:       "TESTTABLE",
				searchQuery: TestSearchQuery{key: TestSearchKey_TEST, method: model.Equals, value: "VALUE"},
			},
			res{
				result:  Test{id: "VALUE"},
				wantErr: false,
			},
		},
		{
			"search with startswith",
			mockDB(t).
				expectGetByQuery("TESTTABLE", "test", "LIKE", "VALUE%"),
			args{
				table:       "TESTTABLE",
				searchQuery: TestSearchQuery{key: TestSearchKey_TEST, method: model.StartsWith, value: "VALUE"},
			},
			res{
				result:  Test{id: "VALUE"},
				wantErr: false,
			},
		},
		{
			"search with contains",
			mockDB(t).
				expectGetByQuery("TESTTABLE", "test", "LIKE", "%VALUE%"),
			args{
				table:       "TESTTABLE",
				searchQuery: TestSearchQuery{key: TestSearchKey_TEST, method: model.Contains, value: "VALUE"},
			},
			res{
				result:  Test{id: "VALUE"},
				wantErr: false,
			},
		},
		{
			"search expect err",
			mockDB(t).
				expectGetByQueryErr("TESTTABLE", "test", "LIKE", "%VALUE%", gorm.ErrRecordNotFound),
			args{
				table:       "TESTTABLE",
				searchQuery: TestSearchQuery{key: TestSearchKey_TEST, method: model.Contains, value: "VALUE"},
			},
			res{
				result:  Test{id: "VALUE"},
				wantErr: true,
				errFunc: func(err error) bool {
					return caos_errs.IsNotFound(err)
				},
			},
		},
		{
			"search with invalid column",
			mockDB(t).
				expectGetByQuery("TESTTABLE", "", "=", "VALUE"),
			args{
				table:       "TESTTABLE",
				searchQuery: TestSearchQuery{key: TestSearchKey_UNDEFINED, method: model.Equals, value: "VALUE"},
			},
			res{
				result:  Test{id: "VALUE"},
				wantErr: true,
				errFunc: func(err error) bool {
					return caos_errs.IsErrorInvalidArgument(err)
				},
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
			"put, ok",
			mockDB(t).
				expectBegin(nil).
				expectSave("TESTTABLE", Test{id: "ID", test: "VALUE"}).
				expectCommit(nil),
			args{
				table:  "TESTTABLE",
				object: &Test{id: "ID", test: "VALUE"},
			},
			res{
				result:  Test{id: "VALUE"},
				wantErr: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getPut := PreparePut(tt.args.table)
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
				result:  Test{id: "VALUE"},
				wantErr: false,
			},
		},
		{
			"delete failes",
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
				result:  Test{id: "VALUE"},
				wantErr: true,
				errFunc: func(err error) bool {
					return caos_errs.IsInternal(err)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getDelete := PrepareDelete(tt.args.table, tt.args.key, tt.args.value)
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
