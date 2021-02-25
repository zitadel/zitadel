package repository

import (
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/jinzhu/gorm"
	"testing"
)

func TestPrepareSearchQuery(t *testing.T) {
	type args struct {
		table         string
		searchRequest SearchRequest
	}
	type res struct {
		count   uint64
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
			"search with no params",
			mockDB(t).
				expectGetSearchRequestNoParams("TESTTABLE", 1, 1),
			args{
				table:         "TESTTABLE",
				searchRequest: TestSearchRequest{},
			},
			res{
				count:   1,
				wantErr: false,
			},
		},
		{
			"search with limit",
			mockDB(t).
				expectGetSearchRequestWithLimit("TESTTABLE", 2, 2, 5),
			args{
				table:         "TESTTABLE",
				searchRequest: TestSearchRequest{limit: 2},
			},
			res{
				count:   5,
				wantErr: false,
			},
		},
		{
			"search with offset",
			mockDB(t).
				expectGetSearchRequestWithOffset("TESTTABLE", 2, 2, 2),
			args{
				table:         "TESTTABLE",
				searchRequest: TestSearchRequest{offset: 2},
			},
			res{
				count:   2,
				wantErr: false,
			},
		},
		{
			"search with sorting asc",
			mockDB(t).
				expectGetSearchRequestWithSorting("TESTTABLE", "ASC", TestSearchKey_ID, 2, 2),
			args{
				table:         "TESTTABLE",
				searchRequest: TestSearchRequest{sortingColumn: TestSearchKey_ID, asc: true},
			},
			res{
				count:   2,
				wantErr: false,
			},
		},
		{
			"search with sorting asc",
			mockDB(t).
				expectGetSearchRequestWithSorting("TESTTABLE", "DESC", TestSearchKey_ID, 2, 2),
			args{
				table:         "TESTTABLE",
				searchRequest: TestSearchRequest{sortingColumn: TestSearchKey_ID},
			},
			res{
				count:   2,
				wantErr: false,
			},
		},
		{
			"search with search query",
			mockDB(t).
				expectGetSearchRequestWithSearchQuery("TESTTABLE", TestSearchKey_ID.ToColumnName(), "=", "AggregateID", 2, 2),
			args{
				table:         "TESTTABLE",
				searchRequest: TestSearchRequest{queries: []SearchQuery{TestSearchQuery{key: TestSearchKey_ID, method: domain.SearchMethodEqualsIgnoreCase, value: "AggregateID"}}},
			},
			res{
				count:   2,
				wantErr: false,
			},
		},
		{
			"search with all params",
			mockDB(t).
				expectGetSearchRequestWithAllParams("TESTTABLE", TestSearchKey_ID.ToColumnName(), "=", "AggregateID", "ASC", TestSearchKey_ID, 2, 2, 2, 5),
			args{
				table:         "TESTTABLE",
				searchRequest: TestSearchRequest{limit: 2, offset: 2, sortingColumn: TestSearchKey_ID, asc: true, queries: []SearchQuery{TestSearchQuery{key: TestSearchKey_ID, method: domain.SearchMethodEqualsIgnoreCase, value: "AggregateID"}}},
			},
			res{
				count:   5,
				wantErr: false,
			},
		},
		{
			"search db error",
			mockDB(t).
				expectGetSearchRequestErr("TESTTABLE", 1, 1, gorm.ErrUnaddressable),
			args{
				table:         "TESTTABLE",
				searchRequest: TestSearchRequest{},
			},
			res{
				count:   1,
				wantErr: true,
				errFunc: caos_errs.IsInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &Test{}
			getQuery := PrepareSearchQuery(tt.args.table, tt.args.searchRequest)
			count, err := getQuery(tt.db.db, res)

			if !tt.res.wantErr && err != nil {
				t.Errorf("got wrong err should be nil: %v ", err)
			}

			if !tt.res.wantErr && count != tt.res.count {
				t.Errorf("got wrong count: %v ", err)
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
