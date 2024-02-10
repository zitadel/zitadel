package query

import (
	"errors"
	"reflect"
	"testing"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/domain"
)

var (
	testTable = table{
		name:          "test_table",
		instanceIDCol: "instance_id",
	}
	testTableAlias = table{
		name:          "test_table",
		alias:         "test_alias",
		instanceIDCol: "instance_id",
	}
	testTable2 = table{
		name:  "test_table2",
		alias: "test_table2",
	}
	testCol = Column{
		name:  "test_col",
		table: testTable,
	}
	testColAlias = Column{
		name:  "test_col",
		table: testTableAlias,
	}
	testCol2 = Column{
		name:  "test_col2",
		table: testTable2,
	}
	testLowerCol = Column{
		name:           "test_lower_col",
		table:          testTable,
		isOrderByLower: true,
	}
	testNoCol = Column{
		name:  "",
		table: testTable,
	}
)

func TestSearchRequest_ToQuery(t *testing.T) {
	type fields struct {
		Offset        uint64
		Limit         uint64
		SortingColumn Column
		Asc           bool
	}
	type want struct {
		stmtAddition string
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
			},
		},
		{
			name: "offset",
			fields: fields{
				Offset: 5,
			},
			want: want{
				stmtAddition: "OFFSET 5",
			},
		},
		{
			name: "limit",
			fields: fields{
				Limit: 5,
			},
			want: want{
				stmtAddition: "LIMIT 5",
			},
		},
		{
			name: "sort asc",
			fields: fields{
				SortingColumn: testCol,
				Asc:           true,
			},
			want: want{
				stmtAddition: "ORDER BY test_table.test_col",
			},
		},
		{
			name: "sort desc",
			fields: fields{
				SortingColumn: testCol,
			},
			want: want{
				stmtAddition: "ORDER BY test_table.test_col DESC",
			},
		},
		{
			name: "sort lower asc",
			fields: fields{
				SortingColumn: testLowerCol,
				Asc:           true,
			},
			want: want{
				stmtAddition: "ORDER BY LOWER(test_table.test_lower_col)",
			},
		},
		{
			name: "all",
			fields: fields{
				Offset:        5,
				Limit:         10,
				SortingColumn: testCol,
				Asc:           true,
			},
			want: want{
				stmtAddition: "ORDER BY test_table.test_col LIMIT 10 OFFSET 5",
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

			query := sq.Select((testCol).identifier()).From(testTable.identifier())
			expectedQuery, _, _ := query.ToSql()

			stmt, _, err := req.toQuery(query).ToSql()
			if len(tt.want.stmtAddition) > 0 {
				expectedQuery += " " + tt.want.stmtAddition
			}
			if expectedQuery != stmt {
				t.Errorf("stmt = %q, want %q", stmt, expectedQuery)
			}

			if err != nil {
				t.Errorf("no error expected but got %v", err)
			}
		})
	}
}

func TestNewSubSelect(t *testing.T) {
	type args struct {
		column  Column
		queries []SearchQuery
	}
	tests := []struct {
		name    string
		args    args
		want    *SubSelect
		wantErr func(error) bool
	}{
		{
			name: "no query nil",
			args: args{
				column:  testCol,
				queries: nil,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrNothingSelected)
			},
		},
		{
			name: "no query zero",
			args: args{
				column:  testCol,
				queries: []SearchQuery{},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrNothingSelected)
			},
		},
		{
			name: "no column 1",
			args: args{
				column:  Column{},
				queries: []SearchQuery{&textQuery{testCol, "horst", TextEquals}},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "no column name 1",
			args: args{
				column:  testNoCol,
				queries: []SearchQuery{&textQuery{testCol, "horst", TextEquals}},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "correct 1",
			args: args{
				column:  testCol,
				queries: []SearchQuery{&textQuery{testCol, "horst", TextEquals}},
			},
			want: &SubSelect{
				Column:  testCol,
				Queries: []SearchQuery{&textQuery{testCol, "horst", TextEquals}},
			},
		},
		{
			name: "correct 3",
			args: args{
				column:  testCol,
				queries: []SearchQuery{&textQuery{testCol, "horst1", TextEquals}, &textQuery{testCol, "horst2", TextEquals}, &textQuery{testCol, "horst3", TextEquals}},
			},
			want: &SubSelect{
				Column:  testCol,
				Queries: []SearchQuery{&textQuery{testCol, "horst1", TextEquals}, &textQuery{testCol, "horst2", TextEquals}, &textQuery{testCol, "horst3", TextEquals}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSubSelect(tt.args.column, tt.args.queries)
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

func TestSubSelect_comp(t *testing.T) {
	type fields struct {
		Column  Column
		Queries []SearchQuery
	}
	type want struct {
		query interface{}
		isNil bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "no queries",
			fields: fields{
				Column:  testCol,
				Queries: []SearchQuery{},
			},
			want: want{
				query: sq.Select("test_table.test_col").From("test_table"),
			},
		},
		{
			name: "queries 1",
			fields: fields{
				Column:  testCol,
				Queries: []SearchQuery{&textQuery{testCol, "horst", TextEquals}},
			},
			want: want{
				query: sq.Select("test_table.test_col").From("test_table").Where(sq.Eq{"test_table.test_col": interface{}("horst")}),
			},
		},
		{
			name: "queries 1 with alias",
			fields: fields{
				Column:  testColAlias,
				Queries: []SearchQuery{&textQuery{testColAlias, "horst", TextEquals}},
			},
			want: want{
				query: sq.Select("test_alias.test_col").From("test_table AS test_alias").Where(sq.Eq{"test_alias.test_col": interface{}("horst")}),
			},
		},
		{
			name: "queries 3",
			fields: fields{
				Column:  testCol,
				Queries: []SearchQuery{&textQuery{testCol, "horst1", TextEquals}, &textQuery{testCol, "horst2", TextEquals}, &textQuery{testCol, "horst3", TextEquals}},
			},
			want: want{
				query: sq.Select("test_table.test_col").From("test_table").From("test_table").Where(sq.Eq{"test_table.test_col": "horst1"}).From("test_table").Where(sq.Eq{"test_table.test_col": "horst2"}).From("test_table").Where(sq.Eq{"test_table.test_col": "horst3"}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SubSelect{
				Column:  tt.fields.Column,
				Queries: tt.fields.Queries,
			}
			query := s.comp()
			if query == nil && tt.want.isNil {
				return
			} else if tt.want.isNil && query != nil {
				t.Error("query should not be nil")
			}

			if !reflect.DeepEqual(query, tt.want.query) {
				t.Errorf("wrong query: want: %v, (%T), got: %v, (%T)", tt.want.query, tt.want.query, query, query)
			}
		})
	}
}

func TestNewColumnComparisonQuery(t *testing.T) {
	type args struct {
		column        Column
		columnCompare Column
		compare       ColumnComparison
	}
	tests := []struct {
		name    string
		args    args
		want    *ColumnComparisonQuery
		wantErr func(error) bool
	}{
		{
			name: "too low compare",
			args: args{
				column:        testCol,
				columnCompare: testCol2,
				compare:       -1,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrInvalidCompare)
			},
		},
		{
			name: "too high compare",
			args: args{
				column:        testCol,
				columnCompare: testCol2,
				compare:       columnCompareMax,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrInvalidCompare)
			},
		},
		{
			name: "no column 1",
			args: args{
				column:        Column{},
				columnCompare: testCol2,
				compare:       ColumnEquals,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "no column 2",
			args: args{
				column:        testCol,
				columnCompare: Column{},
				compare:       ColumnEquals,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "no column name 1",
			args: args{
				column:        testNoCol,
				columnCompare: testCol2,
				compare:       ColumnEquals,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "no column name 2",
			args: args{
				column:        testCol,
				columnCompare: testNoCol,
				compare:       ColumnEquals,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "correct",
			args: args{
				column:        testCol,
				columnCompare: testCol2,
				compare:       ColumnEquals,
			},
			want: &ColumnComparisonQuery{
				Column1: testCol,
				Column2: testCol2,
				Compare: ColumnEquals,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewColumnComparisonQuery(tt.args.column, tt.args.columnCompare, tt.args.compare)
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

func TestColumnComparisonQuery_comp(t *testing.T) {
	type fields struct {
		Column        Column
		ColumnCompare Column
		Compare       ColumnComparison
	}
	type want struct {
		query interface{}
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
				Column:        testCol,
				ColumnCompare: testCol2,
				Compare:       ColumnEquals,
			},
			want: want{
				query: sq.Expr("test_table.test_col = test_table2.test_col2"),
			},
		},
		{
			name: "not equals",
			fields: fields{
				Column:        testCol,
				ColumnCompare: testCol2,
				Compare:       ColumnNotEquals,
			},
			want: want{
				query: sq.Expr("test_table.test_col != test_table2.test_col2"),
			},
		},
		{
			name: "too high comparison",
			fields: fields{
				Column:        testCol,
				ColumnCompare: testCol2,
				Compare:       columnCompareMax,
			},
			want: want{
				isNil: true,
			},
		},
		{
			name: "too low comparison",
			fields: fields{
				Column:        testCol,
				ColumnCompare: testCol2,
				Compare:       -1,
			},
			want: want{
				isNil: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ColumnComparisonQuery{
				Column1: tt.fields.Column,
				Column2: tt.fields.ColumnCompare,
				Compare: tt.fields.Compare,
			}
			query := s.comp()
			if query == nil && tt.want.isNil {
				return
			} else if tt.want.isNil && query != nil {
				t.Error("query should not be nil")
			}

			if !reflect.DeepEqual(query, tt.want.query) {
				t.Errorf("wrong query: want: %v, (%T), got: %v, (%T)", tt.want.query, tt.want.query, query, query)
			}
		})
	}
}

func TestNewListQuery(t *testing.T) {
	type args struct {
		column  Column
		data    interface{}
		compare ListComparison
	}
	tests := []struct {
		name    string
		args    args
		want    *ListQuery
		wantErr func(error) bool
	}{
		{
			name: "too low compare",
			args: args{
				column:  testCol,
				data:    []interface{}{"hurst"},
				compare: -1,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrInvalidCompare)
			},
		},
		{
			name: "too high compare",
			args: args{
				column:  testCol,
				data:    []interface{}{"hurst"},
				compare: listCompareMax,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrInvalidCompare)
			},
		},
		{
			name: "no column",
			args: args{
				column:  Column{},
				data:    []interface{}{"hurst"},
				compare: ListIn,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "no column name",
			args: args{
				column:  testNoCol,
				data:    []interface{}{"hurst"},
				compare: ListIn,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "correct slice",
			args: args{
				column:  testCol,
				data:    []interface{}{"hurst"},
				compare: ListIn,
			},
			want: &ListQuery{
				Column:  testCol,
				Data:    []interface{}{"hurst"},
				Compare: ListIn,
			},
		},
		{
			name: "correct",
			args: args{
				column:  testCol,
				data:    &SubSelect{Column: testCol, Queries: []SearchQuery{&textQuery{testCol, "horst1", TextEquals}}},
				compare: ListIn,
			},
			want: &ListQuery{
				Column:  testCol,
				Data:    &SubSelect{Column: testCol, Queries: []SearchQuery{&textQuery{testCol, "horst1", TextEquals}}},
				Compare: ListIn,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewListQuery(tt.args.column, tt.args.data, tt.args.compare)
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

func TestListQuery_comp(t *testing.T) {
	type fields struct {
		Column  Column
		Data    interface{}
		Compare ListComparison
	}
	type want struct {
		query interface{}
		isNil bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "in list one element",
			fields: fields{
				Column:  testCol,
				Data:    []interface{}{"hurst"},
				Compare: ListIn,
			},
			want: want{
				query: sq.Eq{"test_table.test_col": []interface{}{"hurst"}},
			},
		},
		{
			name: "in list three elements",
			fields: fields{
				Column:  testCol,
				Data:    []interface{}{"hurst1", "hurst2", "hurst3"},
				Compare: ListIn,
			},
			want: want{
				query: sq.Eq{"test_table.test_col": []interface{}{"hurst1", "hurst2", "hurst3"}},
			},
		},
		{
			name: "in string list one element",
			fields: fields{
				Column:  testCol,
				Data:    []string{"hurst"},
				Compare: ListIn,
			},
			want: want{
				query: sq.Eq{"test_table.test_col": []string{"hurst"}},
			},
		},
		{
			name: "in string list three elements",
			fields: fields{
				Column:  testCol,
				Data:    []string{"hurst1", "hurst2", "hurst3"},
				Compare: ListIn,
			},
			want: want{
				query: sq.Eq{"test_table.test_col": []string{"hurst1", "hurst2", "hurst3"}},
			},
		},
		{
			name: "in int list one element",
			fields: fields{
				Column:  testCol,
				Data:    []int{1},
				Compare: ListIn,
			},
			want: want{
				query: sq.Eq{"test_table.test_col": []int{1}},
			},
		},
		{
			name: "in int list three elements",
			fields: fields{
				Column:  testCol,
				Data:    []int{1, 2, 3},
				Compare: ListIn,
			},
			want: want{
				query: sq.Eq{"test_table.test_col": []int{1, 2, 3}},
			},
		},
		{
			name: "in subquery text",
			fields: fields{
				Column:  testCol,
				Data:    &SubSelect{Column: testCol, Queries: []SearchQuery{&textQuery{testCol, "horst", TextEquals}}},
				Compare: ListIn,
			},
			want: want{
				query: sq.Expr("test_table.test_col IN ( SELECT test_table.test_col FROM test_table WHERE test_table.test_col = ? )", "horst"),
			},
		},
		{
			name: "in subquery number",
			fields: fields{
				Column:  testCol,
				Data:    &SubSelect{Column: testCol, Queries: []SearchQuery{&NumberQuery{testCol, 1, NumberEquals}}},
				Compare: ListIn,
			},
			want: want{
				query: sq.Expr("test_table.test_col IN ( SELECT test_table.test_col FROM test_table WHERE test_table.test_col = ? )", 1),
			},
		},
		{
			name: "in subquery column",
			fields: fields{
				Column:  testCol,
				Data:    &SubSelect{Column: testCol, Queries: []SearchQuery{&ColumnComparisonQuery{testCol, ColumnEquals, testCol2}}},
				Compare: ListIn,
			},
			want: want{
				query: sq.Expr("test_table.test_col IN ( SELECT test_table.test_col FROM test_table WHERE test_table.test_col = test_table2.test_col2 )"),
			},
		},
		{
			name: "too high comparison",
			fields: fields{
				Column:  testCol,
				Data:    []interface{}{"hurst"},
				Compare: listCompareMax,
			},
			want: want{
				isNil: true,
			},
		},
		{
			name: "too low comparison",
			fields: fields{
				Column:  testCol,
				Data:    []interface{}{"hurst"},
				Compare: -1,
			},
			want: want{
				isNil: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ListQuery{
				Column:  tt.fields.Column,
				Data:    tt.fields.Data,
				Compare: tt.fields.Compare,
			}
			query := s.comp()
			if query == nil && tt.want.isNil {
				return
			} else if tt.want.isNil && query != nil {
				t.Error("query should not be nil")
			}

			if !reflect.DeepEqual(query, tt.want.query) {
				t.Errorf("wrong query: want: %v, (%T), got: %v, (%T)", tt.want.query, tt.want.query, query, query)
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
		want    *textQuery
		wantErr func(error) bool
	}{
		{
			name: "too low compare",
			args: args{
				column:  testCol,
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
				column:  testCol,
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
				column:  Column{},
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
				column:  testNoCol,
				value:   "hurst",
				compare: TextEquals,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "equals",
			args: args{
				column:  testCol,
				value:   "hurst",
				compare: TextEquals,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hurst",
				Compare: TextEquals,
			},
		},
		{
			name: "equals ignore case",
			args: args{
				column:  testCol,
				value:   "hurst",
				compare: TextEqualsIgnoreCase,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hurst",
				Compare: TextEqualsIgnoreCase,
			},
		},
		{
			name: "equals ignore case % wildcard",
			args: args{
				column:  testCol,
				value:   "hu%rst",
				compare: TextEqualsIgnoreCase,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hu\\%rst",
				Compare: TextEqualsIgnoreCase,
			},
		},
		{
			name: "equals ignore case _ wildcard",
			args: args{
				column:  testCol,
				value:   "hu_rst",
				compare: TextEqualsIgnoreCase,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hu\\_rst",
				Compare: TextEqualsIgnoreCase,
			},
		},
		{
			name: "equals ignore case _, % wildcards",
			args: args{
				column:  testCol,
				value:   "h_urst%",
				compare: TextEqualsIgnoreCase,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "h\\_urst\\%",
				Compare: TextEqualsIgnoreCase,
			},
		},
		{
			name: "not equal",
			args: args{
				column:  testCol,
				value:   "hurst",
				compare: TextNotEquals,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hurst",
				Compare: TextNotEquals,
			},
		},
		{
			name: "starts with",
			args: args{
				column:  testCol,
				value:   "hurst",
				compare: TextStartsWith,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hurst",
				Compare: TextStartsWith,
			},
		},
		{
			name: "starts with _ wildcard",
			args: args{
				column:  testCol,
				value:   "_hurst",
				compare: TextStartsWith,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "\\_hurst",
				Compare: TextStartsWith,
			},
		},
		{
			name: "starts with % wildcard",
			args: args{
				column:  testCol,
				value:   "hurst%",
				compare: TextStartsWith,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hurst\\%",
				Compare: TextStartsWith,
			},
		},
		{
			name: "starts with %, % wildcard",
			args: args{
				column:  testCol,
				value:   "hu%%rst",
				compare: TextStartsWith,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hu\\%\\%rst",
				Compare: TextStartsWith,
			},
		},
		{
			name: "starts with ignore case",
			args: args{
				column:  testCol,
				value:   "hurst",
				compare: TextStartsWithIgnoreCase,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hurst",
				Compare: TextStartsWithIgnoreCase,
			},
		},
		{
			name: "starts with ignore case _ wildcard",
			args: args{
				column:  testCol,
				value:   "hur_st",
				compare: TextStartsWithIgnoreCase,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hur\\_st",
				Compare: TextStartsWithIgnoreCase,
			},
		},
		{
			name: "starts with ignore case % wildcard",
			args: args{
				column:  testCol,
				value:   "hurst%",
				compare: TextStartsWithIgnoreCase,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hurst\\%",
				Compare: TextStartsWithIgnoreCase,
			},
		},
		{
			name: "starts with ignore case _, _ wildcard",
			args: args{
				column:  testCol,
				value:   "h_r_t",
				compare: TextStartsWithIgnoreCase,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "h\\_r\\_t",
				Compare: TextStartsWithIgnoreCase,
			},
		},
		{
			name: "ends with",
			args: args{
				column:  testCol,
				value:   "hurst",
				compare: TextEndsWith,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hurst",
				Compare: TextEndsWith,
			},
		},
		{
			name: "ends with % wildcard",
			args: args{
				column:  testCol,
				value:   "%hurst",
				compare: TextEndsWith,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "\\%hurst",
				Compare: TextEndsWith,
			},
		},
		{
			name: "ends with _ wildcard",
			args: args{
				column:  testCol,
				value:   "hurst_",
				compare: TextEndsWith,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hurst\\_",
				Compare: TextEndsWith,
			},
		},
		{
			name: "ends with _, % wildcard",
			args: args{
				column:  testCol,
				value:   "hurst_%",
				compare: TextEndsWith,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hurst\\_\\%",
				Compare: TextEndsWith,
			},
		},
		{
			name: "ends with ignore case",
			args: args{
				column:  testCol,
				value:   "hurst",
				compare: TextEndsWithIgnoreCase,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hurst",
				Compare: TextEndsWithIgnoreCase,
			},
		},
		{
			name: "ends with ignore case _, %, _ wildcards",
			args: args{
				column:  testCol,
				value:   "h_r_t%",
				compare: TextEndsWithIgnoreCase,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "h\\_r\\_t\\%",
				Compare: TextEndsWithIgnoreCase,
			},
		},
		{
			name: "contains",
			args: args{
				column:  testCol,
				value:   "hurst",
				compare: TextContains,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hurst",
				Compare: TextContains,
			},
		},
		{
			name: "contains % wildcard",
			args: args{
				column:  testCol,
				value:   "%",
				compare: TextContains,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "\\%",
				Compare: TextContains,
			},
		},
		{
			name: "contains ignore csae",
			args: args{
				column:  testCol,
				value:   "hurst",
				compare: TextContainsIgnoreCase,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hurst",
				Compare: TextContainsIgnoreCase,
			},
		},
		{
			name: "contains ignore csae _ wildcard",
			args: args{
				column:  testCol,
				value:   "hurs_",
				compare: TextContainsIgnoreCase,
			},
			want: &textQuery{
				Column:  testCol,
				Text:    "hurs\\_",
				Compare: TextContainsIgnoreCase,
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
				Column:  testCol,
				Text:    "Hurst",
				Compare: TextEquals,
			},
			want: want{
				query: sq.Eq{"test_table.test_col": "Hurst"},
			},
		},
		{
			name: "equals ignore case",
			fields: fields{
				Column:  testCol,
				Text:    "Hurst",
				Compare: TextEqualsIgnoreCase,
			},
			want: want{
				query: sq.ILike{"test_table.test_col": "Hurst"},
			},
		},
		{
			name: "equals ignore case wildcard",
			fields: fields{
				Column:  testCol,
				Text:    "Hu%%rst",
				Compare: TextEqualsIgnoreCase,
			},
			want: want{
				query: sq.ILike{"test_table.test_col": "Hu\\%\\%rst"},
			},
		},
		{
			name: "starts with",
			fields: fields{
				Column:  testCol,
				Text:    "Hurst",
				Compare: TextStartsWith,
			},
			want: want{
				query: sq.Like{"test_table.test_col": "Hurst%"},
			},
		},
		{
			name: "starts with wildcards",
			fields: fields{
				Column:  testCol,
				Text:    "_Hurst%",
				Compare: TextStartsWith,
			},
			want: want{
				query: sq.Like{"test_table.test_col": "\\_Hurst\\%%"},
			},
		},
		{
			name: "starts with ignore case",
			fields: fields{
				Column:  testCol,
				Text:    "Hurst",
				Compare: TextStartsWithIgnoreCase,
			},
			want: want{
				query: sq.ILike{"test_table.test_col": "Hurst%"},
			},
		},
		{
			name: "starts with ignore case wildcards",
			fields: fields{
				Column:  testCol,
				Text:    "Hurst%",
				Compare: TextStartsWithIgnoreCase,
			},
			want: want{
				query: sq.ILike{"test_table.test_col": "Hurst\\%%"},
			},
		},
		{
			name: "ends with",
			fields: fields{
				Column:  testCol,
				Text:    "Hurst",
				Compare: TextEndsWith,
			},
			want: want{
				query: sq.Like{"test_table.test_col": "%Hurst"},
			},
		},
		{
			name: "ends with wildcards",
			fields: fields{
				Column:  testCol,
				Text:    "Hurst%",
				Compare: TextEndsWith,
			},
			want: want{
				query: sq.Like{"test_table.test_col": "%Hurst\\%"},
			},
		},
		{
			name: "ends with ignore case",
			fields: fields{
				Column:  testCol,
				Text:    "Hurst",
				Compare: TextEndsWithIgnoreCase,
			},
			want: want{
				query: sq.ILike{"test_table.test_col": "%Hurst"},
			},
		},
		{
			name: "ends with ignore case wildcards",
			fields: fields{
				Column:  testCol,
				Text:    "%Hurst",
				Compare: TextEndsWithIgnoreCase,
			},
			want: want{
				query: sq.ILike{"test_table.test_col": "%\\%Hurst"},
			},
		},
		{
			name: "contains",
			fields: fields{
				Column:  testCol,
				Text:    "Hurst",
				Compare: TextContains,
			},
			want: want{
				query: sq.Like{"test_table.test_col": "%Hurst%"},
			},
		},
		{
			name: "contains wildcards",
			fields: fields{
				Column:  testCol,
				Text:    "Hu%rst%",
				Compare: TextContains,
			},
			want: want{
				query: sq.Like{"test_table.test_col": "%Hu\\%rst\\%%"},
			},
		},
		{
			name: "containts ignore case",
			fields: fields{
				Column:  testCol,
				Text:    "Hurst",
				Compare: TextContainsIgnoreCase,
			},
			want: want{
				query: sq.ILike{"test_table.test_col": "%Hurst%"},
			},
		},
		{
			name: "contains ignore case wildcards",
			fields: fields{
				Column:  testCol,
				Text:    "%Hurst%",
				Compare: TextContainsIgnoreCase,
			},
			want: want{
				query: sq.ILike{"test_table.test_col": "%\\%Hurst\\%%"},
			},
		},
		{
			name: "list containts",
			fields: fields{
				Column:  testCol,
				Text:    "Hurst",
				Compare: TextListContains,
			},
			want: want{
				query: &listContains{
					col:  testCol,
					args: []interface{}{"Hurst"},
				},
			},
		},
		{
			name: "too high comparison",
			fields: fields{
				Column:  testCol,
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
				Column:  testCol,
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
			s, _ := NewTextQuery(tt.fields.Column, tt.fields.Text, tt.fields.Compare)
			if s == nil {
				// used to check correct behavior of comp
				s = &textQuery{Column: tt.fields.Column, Text: tt.fields.Text, Compare: tt.fields.Compare}
			}
			query := s.comp()
			if query == nil && tt.want.isNil {
				return
			} else if tt.want.isNil && query != nil {
				t.Error("query should not be nil")
			}

			if !reflect.DeepEqual(query, tt.want.query) {
				t.Errorf("wrong query: want: %v, (%T), got: %v, (%T)", tt.want.query, tt.want.query, query, query)
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
				column:  testCol,
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
				column:  testCol,
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
				column:  Column{},
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
				column:  testNoCol,
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
				column:  testCol,
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
				column:  testCol,
				value:   5,
				compare: NumberEquals,
			},
			want: &NumberQuery{
				Column:  testCol,
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
				Column:  testCol,
				Number:  42,
				Compare: NumberEquals,
			},
			want: want{
				query: sq.Eq{"test_table.test_col": 42},
			},
		},
		{
			name: "not equals",
			fields: fields{
				Column:  testCol,
				Number:  42,
				Compare: NumberNotEquals,
			},
			want: want{
				query: sq.NotEq{"test_table.test_col": 42},
			},
		},
		{
			name: "less",
			fields: fields{
				Column:  testCol,
				Number:  42,
				Compare: NumberLess,
			},
			want: want{
				query: sq.Lt{"test_table.test_col": 42},
			},
		},
		{
			name: "greater",
			fields: fields{
				Column:  testCol,
				Number:  42,
				Compare: NumberGreater,
			},
			want: want{
				query: sq.Gt{"test_table.test_col": 42},
			},
		},
		{
			name: "list containts",
			fields: fields{
				Column:  testCol,
				Number:  42,
				Compare: NumberListContains,
			},
			want: want{
				query: &listContains{
					col:  testCol,
					args: []interface{}{42},
				},
			},
		},
		{
			name: "too high comparison",
			fields: fields{
				Column:  testCol,
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
				Column:  testCol,
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
			query := s.comp()
			if query == nil && tt.want.isNil {
				return
			} else if tt.want.isNil && query != nil {
				t.Error("query should not be nil")
			}

			if !reflect.DeepEqual(query, tt.want.query) {
				t.Errorf("wrong query: want: %v, (%T), got: %v, (%T)", tt.want.query, tt.want.query, query, query)
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

func TestNewOrQuery(t *testing.T) {

	type args struct {
		queries []SearchQuery
	}

	singleCorrectQuery, _ := NewTextQuery(testCol, "hello", TextEquals)

	tests := []struct {
		name    string
		args    args
		want    *OrQuery
		wantErr func(error) bool
	}{
		{
			name: "empty values",
			args: args{
				queries: []SearchQuery{},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "correct",
			args: args{
				queries: []SearchQuery{singleCorrectQuery},
			},
			want: &OrQuery{
				queries: []SearchQuery{singleCorrectQuery},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOrQuery(tt.args.queries...)
			if err != nil && tt.wantErr == nil {
				t.Errorf("NewOrQuery() no error expected got %v", err)
				return
			} else if tt.wantErr != nil && !tt.wantErr(err) {
				t.Errorf("NewOrQuery() unexpeted error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestOrQuery_comp(t *testing.T) {
	q1, _ := NewTextQuery(testCol, "hello1", TextEquals)
	q2, _ := NewTextQuery(testCol, "hello2", TextEquals)
	q3, _ := NewTextQuery(testCol2, "world1", TextEquals)
	q4, _ := NewTextQuery(testCol2, "world2", TextEquals)
	orq, _ := NewOrQuery(q3, q4)

	type fields struct {
		queries []SearchQuery
	}
	type want struct {
		query interface{}
		isNil bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "single input",
			fields: fields{
				queries: []SearchQuery{q1},
			},
			want: want{
				query: sq.Or{sq.Eq{"test_table.test_col": "hello1"}},
			},
		},
		{
			name: "multi input",
			fields: fields{
				queries: []SearchQuery{q1, q2},
			},
			want: want{
				query: sq.Or{sq.Eq{"test_table.test_col": "hello1"}, sq.Eq{"test_table.test_col": "hello2"}},
			},
		},
		{
			name: "nested inputs",
			fields: fields{
				queries: []SearchQuery{q1, orq},
			},
			want: want{
				query: sq.Or{sq.Eq{"test_table.test_col": "hello1"}, sq.Or{sq.Eq{"test_table2.test_col2": "world1"}, sq.Eq{"test_table2.test_col2": "world2"}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &OrQuery{
				queries: tt.fields.queries,
			}
			query := s.comp()
			if query == nil && tt.want.isNil {
				return
			} else if tt.want.isNil && query != nil {
				t.Error("query should not be nil")
			}

			if !reflect.DeepEqual(query, tt.want.query) {
				t.Errorf("wrong query: want: %v, (%T), got: %v, (%T)", tt.want.query, tt.want.query, query, query)
			}
		})
	}
}
func TestNewAndQuery(t *testing.T) {

	type args struct {
		queries []SearchQuery
	}

	singleCorrectQuery, _ := NewTextQuery(testCol, "hello", TextEquals)

	tests := []struct {
		name    string
		args    args
		want    *AndQuery
		wantErr func(error) bool
	}{
		{
			name: "empty values",
			args: args{
				queries: []SearchQuery{},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "correct",
			args: args{
				queries: []SearchQuery{singleCorrectQuery},
			},
			want: &AndQuery{
				queries: []SearchQuery{singleCorrectQuery},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAndQuery(tt.args.queries...)
			if err != nil && tt.wantErr == nil {
				t.Errorf("NewAndQuery() no error expected got %v", err)
				return
			} else if tt.wantErr != nil && !tt.wantErr(err) {
				t.Errorf("NewAndQuery() unexpeted error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAndQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAndQuery_comp(t *testing.T) {
	q1, _ := NewTextQuery(testCol, "hello1", TextEquals)
	q2, _ := NewTextQuery(testCol, "hello2", TextEquals)
	q3, _ := NewTextQuery(testCol2, "world1", TextEquals)
	q4, _ := NewTextQuery(testCol2, "world2", TextEquals)
	andq, _ := NewAndQuery(q3, q4)

	type fields struct {
		queries []SearchQuery
	}
	type want struct {
		query interface{}
		isNil bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "single input",
			fields: fields{
				queries: []SearchQuery{q1},
			},
			want: want{
				query: sq.And{sq.Eq{"test_table.test_col": "hello1"}},
			},
		},
		{
			name: "multi input",
			fields: fields{
				queries: []SearchQuery{q1, q2},
			},
			want: want{
				query: sq.And{sq.Eq{"test_table.test_col": "hello1"}, sq.Eq{"test_table.test_col": "hello2"}},
			},
		},
		{
			name: "nested inputs",
			fields: fields{
				queries: []SearchQuery{q1, andq},
			},
			want: want{
				query: sq.And{sq.Eq{"test_table.test_col": "hello1"}, sq.And{sq.Eq{"test_table2.test_col2": "world1"}, sq.Eq{"test_table2.test_col2": "world2"}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AndQuery{
				queries: tt.fields.queries,
			}
			query := s.comp()
			if query == nil && tt.want.isNil {
				return
			} else if tt.want.isNil && query != nil {
				t.Error("query should not be nil")
			}

			if !reflect.DeepEqual(query, tt.want.query) {
				t.Errorf("wrong query: want: %v, (%T), got: %v, (%T)", tt.want.query, tt.want.query, query, query)
			}
		})
	}
}

func TestNewNotQuery(t *testing.T) {

	type args struct {
		query SearchQuery
	}

	singleCorrectQuery, _ := NewTextQuery(testCol, "hello", TextEquals)

	tests := []struct {
		name    string
		args    args
		want    *NotQuery
		wantErr func(error) bool
	}{
		{
			name: "empty query",
			args: args{
				query: nil,
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "correct",
			args: args{
				query: singleCorrectQuery,
			},
			want: &NotQuery{
				query: singleCorrectQuery,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNotQuery(tt.args.query)
			if err != nil && tt.wantErr == nil {
				t.Errorf("NewNotQuery() no error expected got %v", err)
				return
			} else if tt.wantErr != nil && !tt.wantErr(err) {
				t.Errorf("NewNotQuery() unexpeted error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotQuery_comp(t *testing.T) {
	q1, _ := NewTextQuery(testCol, "hello1", TextEquals)

	type fields struct {
		query SearchQuery
	}
	type want struct {
		query interface{}
		isNil bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "single input",
			fields: fields{
				query: q1,
			},
			want: want{
				query: &NotQuery{query: q1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &NotQuery{
				query: tt.fields.query,
			}
			query := s.comp()
			if query == nil && tt.want.isNil {
				return
			} else if tt.want.isNil && query != nil {
				t.Error("query should not be nil")
			}

			if !reflect.DeepEqual(query, tt.want.query) {
				t.Errorf("wrong query: want: %v, (%T), got: %v, (%T)", tt.want.query, tt.want.query, query, query)
			}
		})
	}
}
func TestNotQuery_ToSql(t *testing.T) {
	q, _ := NewTextQuery(testCol, "hello1", TextEquals)
	type fields struct {
		query SearchQuery
	}
	type want struct {
		query string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "single input",
			fields: fields{
				query: q,
			},
			want: want{
				query: "NOT (test_table.test_col = hello1)",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &NotQuery{
				query: tt.fields.query,
			}
			queryStr, _, err := s.ToSql()
			if err != nil {
				t.Error("an error raised ")
			}
			if queryStr == tt.want.query {
				t.Errorf("wrong query: want: %v, (%T), got: %v, (%T)", tt.want.query, tt.want.query, queryStr, queryStr)
			}
		})
	}
}

func TestAndOrQueryCombo(t *testing.T) {
	q1, _ := NewTextQuery(testCol, "hello1", TextEquals)
	q2, _ := NewTextQuery(testCol, "hello2", TextEquals)
	q3, _ := NewTextQuery(testCol2, "world1", TextEquals)
	q4, _ := NewTextQuery(testCol2, "world2", TextEquals)
	andq, _ := NewAndQuery(q3, q4)
	orq, _ := NewOrQuery(q1, q2, andq)

	type fields struct {
		query SearchQuery
	}
	type want struct {
		query interface{}
		isNil bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "OR containing AND query",
			fields: fields{
				query: orq,
			},
			want: want{
				query: sq.Or{sq.Eq{"test_table.test_col": "hello1"}, sq.Eq{"test_table.test_col": "hello2"}, sq.And{sq.Eq{"test_table2.test_col2": "world1"}, sq.Eq{"test_table2.test_col2": "world2"}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.query
			query := s.comp()
			if query == nil && tt.want.isNil {
				return
			} else if tt.want.isNil && query != nil {
				t.Error("query should not be nil")
			}

			if !reflect.DeepEqual(query, tt.want.query) {
				t.Errorf("wrong query: want: %v, (%T), got: %v, (%T)", tt.want.query, tt.want.query, query, query)
			}
		})
	}
}

func TestNewInTextQuery(t *testing.T) {
	type args struct {
		column Column
		value  []string
	}
	tests := []struct {
		name    string
		args    args
		want    *InTextQuery
		wantErr func(error) bool
	}{
		{
			name: "empty values",
			args: args{
				column: testCol,
				value:  []string{},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrEmptyValues)
			},
		},
		{
			name: "no column",
			args: args{
				column: Column{},
				value:  []string{"adler", "hurst"},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "no column name",
			args: args{
				column: testNoCol,
				value:  []string{"adler", "hurst"},
			},
			wantErr: func(err error) bool {
				return errors.Is(err, ErrMissingColumn)
			},
		},
		{
			name: "correct",
			args: args{
				column: testCol,
				value:  []string{"adler", "hurst"},
			},
			want: &InTextQuery{
				Column: testCol,
				Values: []string{"adler", "hurst"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewInTextQuery(tt.args.column, tt.args.value)
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

func TestInTextQuery_comp(t *testing.T) {
	type fields struct {
		Column Column
		Values []string
	}
	type want struct {
		query interface{}
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
				Column: testCol,
				Values: []string{"Adler", "Hurst"},
			},
			want: want{
				query: sq.Eq{"test_table.test_col": []string{"Adler", "Hurst"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InTextQuery{
				Column: tt.fields.Column,
				Values: tt.fields.Values,
			}
			query := s.comp()
			if query == nil && tt.want.isNil {
				return
			} else if tt.want.isNil && query != nil {
				t.Error("query should not be nil")
			}

			if !reflect.DeepEqual(query, tt.want.query) {
				t.Errorf("wrong query: want: %v, (%T), got: %v, (%T)", tt.want.query, tt.want.query, query, query)
			}
		})
	}
}
