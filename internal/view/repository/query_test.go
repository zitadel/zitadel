package repository

import (
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
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
				errFunc: zerrors.IsInternal,
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

func TestSetQuery(t *testing.T) {
	query := mockDB(t).db.Select("test_field").Table("test_table")
	exprPrefix := `(SELECT test_field FROM "test_table"  WHERE `
	type args struct {
		key    ColumnKey
		value  interface{}
		method domain.SearchMethod
	}
	type want struct {
		isErr func(t *testing.T, got error)
		query *gorm.SqlExpr
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "contains",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "asdf",
				method: domain.SearchMethodContains,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(test LIKE ?))", "%asdf%"),
			},
		},
		{
			name: "contains _ wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "as_df",
				method: domain.SearchMethodContains,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(test LIKE ?))", "%as\\_df%"),
			},
		},
		{
			name: "contains % wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "as%df",
				method: domain.SearchMethodContains,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(test LIKE ?))", "%as\\%df%"),
			},
		},
		{
			name: "contains % wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "a_s%d_f",
				method: domain.SearchMethodContains,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(test LIKE ?))", "%a\\_s\\%d\\_f%"),
			},
		},
		{
			name: "starts with",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "asdf",
				method: domain.SearchMethodStartsWith,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(test LIKE ?))", "asdf%"),
			},
		},
		{
			name: "starts with _ wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "as_df",
				method: domain.SearchMethodStartsWith,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(test LIKE ?))", "as\\_df%"),
			},
		},
		{
			name: "starts with % wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "as%df",
				method: domain.SearchMethodStartsWith,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(test LIKE ?))", "as\\%df%"),
			},
		},
		{
			name: "starts with % wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "a_s%d_f",
				method: domain.SearchMethodStartsWith,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(test LIKE ?))", "a\\_s\\%d\\_f%"),
			},
		},
		{
			name: "ends with",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "asdf",
				method: domain.SearchMethodEndsWith,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(test LIKE ?))", "%asdf"),
			},
		},
		{
			name: "ends with _ wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "as_df",
				method: domain.SearchMethodEndsWith,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(test LIKE ?))", "%as\\_df"),
			},
		},
		{
			name: "ends with % wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "as%df",
				method: domain.SearchMethodEndsWith,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(test LIKE ?))", "%as\\%df"),
			},
		},
		{
			name: "ends with % wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "a_s%d_f",
				method: domain.SearchMethodEndsWith,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(test LIKE ?))", "%a\\_s\\%d\\_f"),
			},
		},
		{
			name: "starts with ignore case",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "asdf",
				method: domain.SearchMethodStartsWithIgnoreCase,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(LOWER(test) LIKE LOWER(?)))", "asdf%"),
			},
		},
		{
			name: "starts with ignore case _ wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "as_df",
				method: domain.SearchMethodStartsWithIgnoreCase,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(LOWER(test) LIKE LOWER(?)))", "as\\_df%"),
			},
		},
		{
			name: "starts with ignore case % wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "as%df",
				method: domain.SearchMethodStartsWithIgnoreCase,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(LOWER(test) LIKE LOWER(?)))", "as\\%df%"),
			},
		},
		{
			name: "starts with ignore case % wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "a_s%d_f",
				method: domain.SearchMethodStartsWithIgnoreCase,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(LOWER(test) LIKE LOWER(?)))", "a\\_s\\%d\\_f%"),
			},
		},
		{
			name: "ends with ignore case",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "asdf",
				method: domain.SearchMethodEndsWithIgnoreCase,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(LOWER(test) LIKE LOWER(?)))", "%asdf"),
			},
		},
		{
			name: "ends with ignore case _ wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "as_df",
				method: domain.SearchMethodEndsWithIgnoreCase,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(LOWER(test) LIKE LOWER(?)))", "%as\\_df"),
			},
		},
		{
			name: "ends with ignore case % wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "as%df",
				method: domain.SearchMethodEndsWithIgnoreCase,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(LOWER(test) LIKE LOWER(?)))", "%as\\%df"),
			},
		},
		{
			name: "ends with ignore case % wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "a_s%d_f",
				method: domain.SearchMethodEndsWithIgnoreCase,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(LOWER(test) LIKE LOWER(?)))", "%a\\_s\\%d\\_f"),
			},
		},
		{
			name: "contains ignore case",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "asdf",
				method: domain.SearchMethodContainsIgnoreCase,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(LOWER(test) LIKE LOWER(?)))", "%asdf%"),
			},
		},
		{
			name: "contains ignore case _ wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "as_df",
				method: domain.SearchMethodContainsIgnoreCase,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(LOWER(test) LIKE LOWER(?)))", "%as\\_df%"),
			},
		},
		{
			name: "contains ignore case % wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "as%df",
				method: domain.SearchMethodContainsIgnoreCase,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(LOWER(test) LIKE LOWER(?)))", "%as\\%df%"),
			},
		},
		{
			name: "contains ignore case % wildcard",
			args: args{
				key:    TestSearchKey_TEST,
				value:  "a_s%d_f",
				method: domain.SearchMethodContainsIgnoreCase,
			},
			want: want{
				query: gorm.Expr(exprPrefix+"(LOWER(test) LIKE LOWER(?)))", "%a\\_s\\%d\\_f%"),
			},
		},
	}
	for _, tt := range tests {
		if tt.want.isErr == nil {
			tt.want.isErr = func(t *testing.T, got error) {
				if got == nil {
					return
				}
				t.Errorf("no error expected got: %v", got)
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetQuery(query, tt.args.key, tt.args.value, tt.args.method)
			tt.want.isErr(t, err)
			if !reflect.DeepEqual(got.SubQuery(), tt.want.query) {
				t.Errorf("unexpected query: \nwant: %v\n got: %v", *tt.want.query, *got.SubQuery())
			}
		})
	}
}
