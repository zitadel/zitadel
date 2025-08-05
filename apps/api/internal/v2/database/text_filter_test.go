package database

import (
	"reflect"
	"testing"
)

func TestNewTextEqual(t *testing.T) {
	type args struct {
		constructor func(t string) *TextFilter[string]
		t           string
	}
	tests := []struct {
		name string
		args args
		want *TextFilter[string]
	}{
		{
			name: "NewTextEqual",
			args: args{
				constructor: NewTextEqual[string],
				t:           "text",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textEqual,
					value: "text",
				},
			},
		},
		{
			name: "NewTextUnequal",
			args: args{
				constructor: NewTextUnequal[string],
				t:           "text",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textUnequal,
					value: "text",
				},
			},
		},
		{
			name: "NewTextEqualInsensitive",
			args: args{
				constructor: NewTextEqualInsensitive[string],
				t:           "text",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textEqualInsensitive,
					value: "text",
				},
			},
		},
		{
			name: "NewTextEqualInsensitive check lower",
			args: args{
				constructor: NewTextEqualInsensitive[string],
				t:           "tEXt",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textEqualInsensitive,
					value: "text",
				},
			},
		},
		{
			name: "NewTextUnequalInsensitive",
			args: args{
				constructor: NewTextUnequalInsensitive[string],
				t:           "text",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textUnequalInsensitive,
					value: "text",
				},
			},
		},
		{
			name: "NewTextUnequalInsensitive check lower",
			args: args{
				constructor: NewTextUnequalInsensitive[string],
				t:           "tEXt",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textUnequalInsensitive,
					value: "text",
				},
			},
		},
		{
			name: "NewTextStartsWith",
			args: args{
				constructor: NewTextStartsWith[string],
				t:           "text",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textStartsWith,
					value: "text",
				},
			},
		},
		{
			name: "NewTextStartsWithInsensitive",
			args: args{
				constructor: NewTextStartsWithInsensitive[string],
				t:           "text",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textStartsWithInsensitive,
					value: "text",
				},
			},
		},
		{
			name: "NewTextStartsWithInsensitive check lower",
			args: args{
				constructor: NewTextStartsWithInsensitive[string],
				t:           "tEXt",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textStartsWithInsensitive,
					value: "text",
				},
			},
		},
		{
			name: "NewTextEndsWith",
			args: args{
				constructor: NewTextEndsWith[string],
				t:           "text",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textEndsWith,
					value: "text",
				},
			},
		},
		{
			name: "NewTextEndsWithInsensitive",
			args: args{
				constructor: NewTextEndsWithInsensitive[string],
				t:           "text",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textEndsWithInsensitive,
					value: "text",
				},
			},
		},
		{
			name: "NewTextEndsWithInsensitive check lower",
			args: args{
				constructor: NewTextEndsWithInsensitive[string],
				t:           "tEXt",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textEndsWithInsensitive,
					value: "text",
				},
			},
		},
		{
			name: "NewTextContains",
			args: args{
				constructor: NewTextContains[string],
				t:           "text",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textContains,
					value: "text",
				},
			},
		},
		{
			name: "NewTextContainsInsensitive",
			args: args{
				constructor: NewTextContainsInsensitive[string],
				t:           "text",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textContainsInsensitive,
					value: "text",
				},
			},
		},
		{
			name: "NewTextContainsInsensitive to lower",
			args: args{
				constructor: NewTextContainsInsensitive[string],
				t:           "tEXt",
			},
			want: &TextFilter[string]{
				Filter: Filter[textCompare, string]{
					comp:  textContainsInsensitive,
					value: "text",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.constructor(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTextEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextConditionWrite(t *testing.T) {
	type args struct {
		constructor func(t string) *TextFilter[string]
		t           string
	}
	tests := []struct {
		name string
		args args
		want wantQuery
	}{
		{
			name: "NewTextEqual",
			args: args{
				constructor: NewTextEqual[string],
				t:           "text",
			},
			want: wantQuery{
				query: "test = $1",
				args:  []any{"text"},
			},
		},
		{
			name: "NewTextUnequal",
			args: args{
				constructor: NewTextUnequal[string],
				t:           "text",
			},
			want: wantQuery{
				query: "test <> $1",
				args:  []any{"text"},
			},
		},
		{
			name: "NewTextEqualInsensitive",
			args: args{
				constructor: NewTextEqualInsensitive[string],
				t:           "text",
			},
			want: wantQuery{
				query: "LOWER(test) = $1",
				args:  []any{"text"},
			},
		},
		{
			name: "NewTextUnequalInsensitive",
			args: args{
				constructor: NewTextUnequalInsensitive[string],
				t:           "text",
			},
			want: wantQuery{
				query: "test <> $1",
				args:  []any{"text"},
			},
		},
		{
			name: "NewTextStartsWith",
			args: args{
				constructor: NewTextStartsWith[string],
				t:           "text",
			},
			want: wantQuery{
				query: "test LIKE $1",
				args:  []any{"text"},
			},
		},
		{
			name: "NewTextStartsWithInsensitive",
			args: args{
				constructor: NewTextStartsWithInsensitive[string],
				t:           "text",
			},
			want: wantQuery{
				query: "LOWER(test) LIKE $1",
				args:  []any{"text"},
			},
		},
		{
			name: "NewTextEndsWith",
			args: args{
				constructor: NewTextEndsWith[string],
				t:           "text",
			},
			want: wantQuery{
				query: "test LIKE $1",
				args:  []any{"text"},
			},
		},
		{
			name: "NewTextEndsWithInsensitive",
			args: args{
				constructor: NewTextEndsWithInsensitive[string],
				t:           "text",
			},
			want: wantQuery{
				query: "LOWER(test) LIKE $1",
				args:  []any{"text"},
			},
		},
		{
			name: "NewTextContains",
			args: args{
				constructor: NewTextContains[string],
				t:           "text",
			},
			want: wantQuery{
				query: "test LIKE $1",
				args:  []any{"text"},
			},
		},
		{
			name: "NewTextContainsInsensitive",
			args: args{
				constructor: NewTextContainsInsensitive[string],
				t:           "text",
			},
			want: wantQuery{
				query: "LOWER(test) LIKE $1",
				args:  []any{"text"},
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
