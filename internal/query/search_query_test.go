package query

import (
	"errors"
	"reflect"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
)

type testCol struct{}

func (col *testCol) toColumnName() string {
	return "test"
}

type testNoCol struct{}

func (col *testNoCol) toColumnName() string {
	return ""
}

func TestSearchRequest_ToQuery(t *testing.T) {
	type fields struct {
		Offset        uint64
		Limit         uint64
		SortingColumn Column
		Asc           bool
	}
	type want struct {
		stmtAddition string
		args         []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name:   "no queries",
			fields: fields{},
			want: want{
				stmtAddition: "",
				args:         nil,
			},
		},
		{
			name: "offset",
			fields: fields{
				Offset: 5,
			},
			want: want{
				stmtAddition: "OFFSET 5",
				args:         nil,
			},
		},
		{
			name: "limit",
			fields: fields{
				Limit: 5,
			},
			want: want{
				stmtAddition: "LIMIT 5",
				args:         nil,
			},
		},
		{
			name: "sort asc",
			fields: fields{
				SortingColumn: &testCol{},
				Asc:           true,
			},
			want: want{
				stmtAddition: "ORDER BY LOWER(?)",
				args:         []interface{}{"test"},
			},
		},
		{
			name: "sort desc",
			fields: fields{
				SortingColumn: &testCol{},
			},
			want: want{
				stmtAddition: "ORDER BY LOWER(?) DESC",
				args:         []interface{}{"test"},
			},
		},
		{
			name: "all",
			fields: fields{
				Offset:        5,
				Limit:         10,
				SortingColumn: &testCol{},
				Asc:           true,
			},
			want: want{
				stmtAddition: "ORDER BY LOWER(?) LIMIT 10 OFFSET 5",
				args:         []interface{}{"test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &SearchRequest{
				Offset:        tt.fields.Offset,
				Limit:         tt.fields.Limit,
				SortingColumn: tt.fields.SortingColumn,
				Asc:           tt.fields.Asc,
			}

			query := sq.Select((&testCol{}).toColumnName()).From("test_table")
			expectedQuery, _, _ := query.ToSql()

			stmt, args, err := req.ToQuery(query).ToSql()
			if len(tt.want.stmtAddition) > 0 {
				expectedQuery += " " + tt.want.stmtAddition
			}
			if expectedQuery != stmt {
				t.Errorf("stmt = %q, want %q", stmt, expectedQuery)
			}

			if !reflect.DeepEqual(args, tt.want.args) {
				t.Errorf("args = %v, want %v", args, tt.want.stmtAddition)
			}

			if err != nil {
				t.Errorf("no error expected but got %v", err)
			}
		})
	}
}

func TestNewTextQuery(t *testing.T) {
	type args struct {
		column  Column
		value   string
		compare TextComparison
	}
	tests := []struct {
		name    string
		args    args
		want    *TextQuery
		wantErr func(error) bool
	}{
		{
			name: "too low compare",
			args: args{
				column:  &testCol{},
				value:   "hurst",
				compare: -1,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrInvalidCompare)
			},
		},
		{
			name: "too high compare",
			args: args{
				column:  &testCol{},
				value:   "hurst",
				compare: textCompareMax,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrInvalidCompare)
			},
		},
		{
			name: "no column",
			args: args{
				column:  nil,
				value:   "hurst",
				compare: TextEquals,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "no column name",
			args: args{
				column:  &testNoCol{},
				value:   "hurst",
				compare: TextEquals,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "correct",
			args: args{
				column:  &testCol{},
				value:   "hurst",
				compare: TextEquals,
			},
			want: &TextQuery{
				Column:  &testCol{},
				Text:    "hurst",
				Compare: TextEquals,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTextQuery(tt.args.column, tt.args.value, tt.args.compare)
			if err != nil && tt.wantErr == nil {
				t.Errorf("NewTextQuery() no error expected got %v", err)
				return
			} else if tt.wantErr != nil && !tt.wantErr(err) {
				t.Errorf("NewTextQuery() unexpeted error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTextQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextQuery_comp(t *testing.T) {
	type fields struct {
		Column  Column
		Text    string
		Compare TextComparison
	}
	type want struct {
		stmt  string
		args  []interface{}
		isNil bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "equals",
			fields: fields{
				Column:  &testCol{},
				Text:    "Hurst",
				Compare: TextEquals,
			},
			want: want{
				stmt: "test = ?",
				args: []interface{}{"Hurst"},
			},
		},
		{
			name: "equals ignore case",
			fields: fields{
				Column:  &testCol{},
				Text:    "Hurst",
				Compare: TextEqualsIgnoreCase,
			},
			want: want{
				stmt: "LOWER(test) = ?",
				args: []interface{}{"hurst"},
			},
		},
		{
			name: "starts with",
			fields: fields{
				Column:  &testCol{},
				Text:    "Hurst",
				Compare: TextStartsWith,
			},
			want: want{
				stmt: "test LIKE ?",
				args: []interface{}{"Hurst%"},
			},
		},
		{
			name: "starts with ignore case",
			fields: fields{
				Column:  &testCol{},
				Text:    "Hurst",
				Compare: TextStartsWithIgnoreCase,
			},
			want: want{
				stmt: "LOWER(test) LIKE ?",
				args: []interface{}{"hurst%"},
			},
		},
		{
			name: "ends with",
			fields: fields{
				Column:  &testCol{},
				Text:    "Hurst",
				Compare: TextEndsWith,
			},
			want: want{
				stmt: "test LIKE ?",
				args: []interface{}{"%Hurst"},
			},
		},
		{
			name: "ends with ignore case",
			fields: fields{
				Column:  &testCol{},
				Text:    "Hurst",
				Compare: TextEndsWithIgnoreCase,
			},
			want: want{
				stmt: "LOWER(test) LIKE ?",
				args: []interface{}{"%hurst"},
			},
		},
		{
			name: "contains",
			fields: fields{
				Column:  &testCol{},
				Text:    "Hurst",
				Compare: TextContains,
			},
			want: want{
				stmt: "test LIKE ?",
				args: []interface{}{"%Hurst%"},
			},
		},
		{
			name: "containts ignore case",
			fields: fields{
				Column:  &testCol{},
				Text:    "Hurst",
				Compare: TextContainsIgnoreCase,
			},
			want: want{
				stmt: "LOWER(test) LIKE ?",
				args: []interface{}{"%hurst%"},
			},
		},
		{
			name: "too high comparison",
			fields: fields{
				Column:  &testCol{},
				Text:    "Hurst",
				Compare: textCompareMax,
			},
			want: want{
				isNil: true,
			},
		},
		{
			name: "too low comparison",
			fields: fields{
				Column:  &testCol{},
				Text:    "Hurst",
				Compare: -1,
			},
			want: want{
				isNil: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TextQuery{
				Column:  tt.fields.Column,
				Text:    tt.fields.Text,
				Compare: tt.fields.Compare,
			}
			query := s.comp()
			if query == nil && tt.want.isNil {
				return
			} else if tt.want.isNil && query != nil {
				t.Error("query should not be nil")
			}
			stmt, args, err := query.ToSql()
			if err != nil {
				t.Errorf("no err expected: %v", err)
			}
			if stmt != tt.want.stmt {
				t.Errorf("stmt = %v, want %v", stmt, tt.want.stmt)
			}
			if !reflect.DeepEqual(args, tt.want.args) {
				t.Errorf("args = %v, want %v", args, tt.want.args)
			}
		})
	}
}

func TestTextCompareFromMethod(t *testing.T) {
	type args struct {
		m domain.SearchMethod
	}
	tests := []struct {
		name string
		args args
		want TextComparison
	}{
		{
			name: "equals",
			args: args{
				m: domain.SearchMethodEquals,
			},
			want: TextEquals,
		},
		{
			name: "equals ignore case",
			args: args{
				m: domain.SearchMethodEqualsIgnoreCase,
			},
			want: TextEqualsIgnoreCase,
		},
		{
			name: "starts with",
			args: args{
				m: domain.SearchMethodStartsWith,
			},
			want: TextStartsWith,
		},
		{
			name: "starts with ignore case",
			args: args{
				m: domain.SearchMethodStartsWithIgnoreCase,
			},
			want: TextStartsWithIgnoreCase,
		},
		{
			name: "ends with",
			args: args{
				m: domain.SearchMethodEndsWith,
			},
			want: TextEndsWith,
		},
		{
			name: "ends with ignore case",
			args: args{
				m: domain.SearchMethodEndsWithIgnoreCase,
			},
			want: TextEndsWithIgnoreCase,
		},
		{
			name: "contains",
			args: args{
				m: domain.SearchMethodContains,
			},
			want: TextContains,
		},
		{
			name: "containts ignore case",
			args: args{
				m: domain.SearchMethodContainsIgnoreCase,
			},
			want: TextContainsIgnoreCase,
		},
		{
			name: "invalid search method",
			args: args{
				m: -1,
			},
			want: textCompareMax,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TextCompareFromMethod(tt.args.m); got != tt.want {
				t.Errorf("TextCompareFromMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
