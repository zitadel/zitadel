package database

import (
	"reflect"
	"testing"
)

func TestNewNumberConstructors(t *testing.T) {
	type args struct {
		constructor func(t int8) *NumberFilter[int8]
		t           int8
	}
	tests := []struct {
		name string
		args args
		want *NumberFilter[int8]
	}{
		{
			name: "NewNumberEqual",
			args: args{
				constructor: NewNumberEquals[int8],
				t:           10,
			},
			want: &NumberFilter[int8]{
				Filter: Filter[numberCompare, int8]{
					comp:  numberEqual,
					value: 10,
				},
			},
		},
		{
			name: "NewNumberAtLeast",
			args: args{
				constructor: NewNumberAtLeast[int8],
				t:           10,
			},
			want: &NumberFilter[int8]{
				Filter: Filter[numberCompare, int8]{
					comp:  numberAtLeast,
					value: 10,
				},
			},
		},
		{
			name: "NewNumberAtMost",
			args: args{
				constructor: NewNumberAtMost[int8],
				t:           10,
			},
			want: &NumberFilter[int8]{
				Filter: Filter[numberCompare, int8]{
					comp:  numberAtMost,
					value: 10,
				},
			},
		},
		{
			name: "NewNumberGreater",
			args: args{
				constructor: NewNumberGreater[int8],
				t:           10,
			},
			want: &NumberFilter[int8]{
				Filter: Filter[numberCompare, int8]{
					comp:  numberGreater,
					value: 10,
				},
			},
		},
		{
			name: "NewNumberLess",
			args: args{
				constructor: NewNumberLess[int8],
				t:           10,
			},
			want: &NumberFilter[int8]{
				Filter: Filter[numberCompare, int8]{
					comp:  numberLess,
					value: 10,
				},
			},
		},
		{
			name: "NewNumberUnequal",
			args: args{
				constructor: NewNumberUnequal[int8],
				t:           10,
			},
			want: &NumberFilter[int8]{
				Filter: Filter[numberCompare, int8]{
					comp:  numberUnequal,
					value: 10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.constructor(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("number constructor = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNumberConditionWrite(t *testing.T) {
	type args struct {
		constructor func(t int8) *NumberFilter[int8]
		t           int8
	}
	tests := []struct {
		name string
		args args
		want wantQuery
	}{
		{
			name: "NewNumberEqual",
			args: args{
				constructor: NewNumberEquals[int8],
				t:           10,
			},
			want: wantQuery{
				query: "test = $1",
				args:  []any{int8(10)},
			},
		},
		{
			name: "NewNumberAtLeast",
			args: args{
				constructor: NewNumberAtLeast[int8],
				t:           10,
			},
			want: wantQuery{
				query: "test >= $1",
				args:  []any{int8(10)},
			},
		},
		{
			name: "NewNumberAtMost",
			args: args{
				constructor: NewNumberAtMost[int8],
				t:           10,
			},
			want: wantQuery{
				query: "test <= $1",
				args:  []any{int8(10)},
			},
		},
		{
			name: "NewNumberGreater",
			args: args{
				constructor: NewNumberGreater[int8],
				t:           10,
			},
			want: wantQuery{
				query: "test > $1",
				args:  []any{int8(10)},
			},
		},
		{
			name: "NewNumberLess",
			args: args{
				constructor: NewNumberLess[int8],
				t:           10,
			},
			want: wantQuery{
				query: "test < $1",
				args:  []any{int8(10)},
			},
		},
		{
			name: "NewNumberUnequal",
			args: args{
				constructor: NewNumberUnequal[int8],
				t:           10,
			},
			want: wantQuery{
				query: "test <> $1",
				args:  []any{int8(10)},
			},
		},
	}
	for _, tt := range tests {
		var stmt Statement
		t.Run(tt.name, func(t *testing.T) {
			tt.args.constructor(tt.args.t).Write(&stmt, "test")
			assertQuery(t, &stmt, tt.want)
		})
	}
}

func TestNumberBetween(t *testing.T) {
	filter := NewNumberBetween[int8](10, 20)

	if !reflect.DeepEqual(filter, &NumberBetweenFilter[int8]{min: 10, max: 20}) {
		t.Errorf("unexpected filter: %v", filter)
	}

	var stmt Statement
	filter.Write(&stmt, "test")
	if stmt.String() != "test >= $1 AND test <= $2" {
		t.Errorf("unexpected query: got: %q", stmt.String())
	}

	if len(stmt.Args()) != 2 {
		t.Errorf("unexpected length of args: got %d", len(stmt.Args()))
		return
	}

	if !reflect.DeepEqual(int8(10), stmt.Args()[0]) {
		t.Errorf("unexpected arg at position 0: want: 10, got: %v", stmt.Args()[0])
	}
	if !reflect.DeepEqual(int8(20), stmt.Args()[1]) {
		t.Errorf("unexpected arg at position 1: want: 20, got: %v", stmt.Args()[1])
	}
}
