package query

import (
	"errors"
	"reflect"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
	"github.com/lib/pq"
)

type testCol struct{}

func (col *testCol) FullColumnName() string {
	return "test"
}

func (col *testCol) ColumnName() string {
	return "test"
}

type testNoCol struct{}

func (col *testNoCol) FullColumnName() string {
	return ""
}

func (col *testNoCol) ColumnName() string {
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

			query := sq.Select((&testCol{}).FullColumnName()).From("test_table")
			expectedQuery, _, _ := query.ToSql()

			stmt, args, err := req.toQuery(query).ToSql()
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
		query interface{}
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
				query: sq.Eq{"test": "Hurst"},
				args:  nil,
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
				query: sq.ILike{"test": "Hurst"},
				args:  nil,
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
				query: sq.Like{"test": "Hurst%"},
				args:  nil,
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
				query: sq.ILike{"test": "Hurst%"},
				args:  nil,
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
				query: sq.Like{"test": "%Hurst"},
				args:  nil,
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
				query: sq.ILike{"test": "%Hurst"},
				args:  nil,
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
				query: sq.Like{"test": "%Hurst%"},
				args:  nil,
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
				query: sq.ILike{"test": "%Hurst%"},
				args:  nil,
			},
		},
		{
			name: "list containts",
			fields: fields{
				Column:  &testCol{},
				Text:    "Hurst",
				Compare: TextListContains,
			},
			want: want{
				query: "test @> ? ",
				args:  []interface{}{pq.StringArray{"Hurst"}},
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
			query, args := s.comp()
			if query == nil && tt.want.isNil {
				return
			} else if tt.want.isNil && query != nil {
				t.Error("query should not be nil")
			}

			if !reflect.DeepEqual(query, tt.want.query) {
				t.Errorf("wrong query: want: %v, (%T), got: %v, (%T)", tt.want.query, tt.want.query, query, query)
			}

			if !reflect.DeepEqual(args, tt.want.args) {
				t.Errorf("wrong args: want: %v, (%T), got: %v (%T)", tt.want.args, tt.want.args, args, args)
			}
		})
	}
}

func TestTextComparisonFromMethod(t *testing.T) {
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
			name: "list contains",
			args: args{
				m: domain.SearchMethodListContains,
			},
			want: TextListContains,
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
			if got := TextComparisonFromMethod(tt.args.m); got != tt.want {
				t.Errorf("TextCompareFromMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNumberQuery(t *testing.T) {
	type args struct {
		column  Column
		value   interface{}
		compare NumberComparison
	}
	tests := []struct {
		name    string
		args    args
		want    *NumberQuery
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
				compare: numberCompareMax,
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
				compare: NumberEquals,
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
				compare: NumberEquals,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "no number",
			args: args{
				column:  &testCol{},
				value:   "hurst",
				compare: NumberEquals,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrInvalidNumber)
			},
		},
		{
			name: "correct",
			args: args{
				column:  &testCol{},
				value:   5,
				compare: NumberEquals,
			},
			want: &NumberQuery{
				Column:  &testCol{},
				Number:  5,
				Compare: NumberEquals,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNumberQuery(tt.args.column, tt.args.value, tt.args.compare)
			if err != nil && tt.wantErr == nil {
				t.Errorf("NewNumberQuery() no error expected got %v", err)
				return
			} else if tt.wantErr != nil && !tt.wantErr(err) {
				t.Errorf("NewNumberQuery() unexpeted error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNumberQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberQuery_comp(t *testing.T) {
	type fields struct {
		Column  Column
		Number  interface{}
		Compare NumberComparison
	}
	type want struct {
		query interface{}
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
				Number:  42,
				Compare: NumberEquals,
			},
			want: want{
				query: sq.Eq{"test": 42},
				args:  nil,
			},
		},
		{
			name: "not equals",
			fields: fields{
				Column:  &testCol{},
				Number:  42,
				Compare: NumberNotEquals,
			},
			want: want{
				query: sq.NotEq{"test": 42},
				args:  nil,
			},
		},
		{
			name: "less",
			fields: fields{
				Column:  &testCol{},
				Number:  42,
				Compare: NumberLess,
			},
			want: want{
				query: sq.Lt{"test": 42},
				args:  nil,
			},
		},
		{
			name: "greater",
			fields: fields{
				Column:  &testCol{},
				Number:  42,
				Compare: NumberGreater,
			},
			want: want{
				query: sq.Gt{"test": 42},
				args:  nil,
			},
		},
		{
			name: "list containts",
			fields: fields{
				Column:  &testCol{},
				Number:  42,
				Compare: NumberListContains,
			},
			want: want{
				query: "test @> ? ",
				args:  []interface{}{pq.Array(42)},
			},
		},
		{
			name: "too high comparison",
			fields: fields{
				Column:  &testCol{},
				Number:  42,
				Compare: numberCompareMax,
			},
			want: want{
				isNil: true,
			},
		},
		{
			name: "too low comparison",
			fields: fields{
				Column:  &testCol{},
				Number:  42,
				Compare: -1,
			},
			want: want{
				isNil: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &NumberQuery{
				Column:  tt.fields.Column,
				Number:  tt.fields.Number,
				Compare: tt.fields.Compare,
			}
			query, args := s.comp()
			if query == nil && tt.want.isNil {
				return
			} else if tt.want.isNil && query != nil {
				t.Error("query should not be nil")
			}

			if !reflect.DeepEqual(query, tt.want.query) {
				t.Errorf("wrong query: want: %v, (%T), got: %v, (%T)", tt.want.query, tt.want.query, query, query)
			}

			if !reflect.DeepEqual(args, tt.want.args) {
				t.Errorf("wrong args: want: %v, (%T), got: %v (%T)", tt.want.args, tt.want.args, args, args)
			}
		})
	}
}

func TestNumberComparisonFromMethod(t *testing.T) {
	type args struct {
		m domain.SearchMethod
	}
	tests := []struct {
		name string
		args args
		want NumberComparison
	}{
		{
			name: "equals",
			args: args{
				m: domain.SearchMethodEquals,
			},
			want: NumberEquals,
		},
		{
			name: "not equals",
			args: args{
				m: domain.SearchMethodNotEquals,
			},
			want: NumberNotEquals,
		},
		{
			name: "less than",
			args: args{
				m: domain.SearchMethodLessThan,
			},
			want: NumberLess,
		},
		{
			name: "greater than",
			args: args{
				m: domain.SearchMethodGreaterThan,
			},
			want: NumberGreater,
		},
		{
			name: "list contains",
			args: args{
				m: domain.SearchMethodListContains,
			},
			want: NumberListContains,
		},
		{
			name: "invalid search method",
			args: args{
				m: -1,
			},
			want: numberCompareMax,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumberComparisonFromMethod(tt.args.m); got != tt.want {
				t.Errorf("TextCompareFromMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
