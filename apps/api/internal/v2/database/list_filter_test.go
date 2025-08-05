package database

import (
	"reflect"
	"testing"
)

func TestNewListConstructors(t *testing.T) {
	type args struct {
		constructor func(t ...string) *ListFilter[string]
		t           []string
	}
	tests := []struct {
		name string
		args args
		want *ListFilter[string]
	}{
		{
			name: "NewListEquals",
			args: args{
				constructor: NewListEquals[string],
				t:           []string{"as", "df"},
			},
			want: &ListFilter[string]{
				comp: listEqual,
				list: []string{"as", "df"},
			},
		},
		{
			name: "NewListContains",
			args: args{
				constructor: NewListContains[string],
				t:           []string{"as", "df"},
			},
			want: &ListFilter[string]{
				comp: listContain,
				list: []string{"as", "df"},
			},
		},
		{
			name: "NewListNotContains",
			args: args{
				constructor: NewListNotContains[string],
				t:           []string{"as", "df"},
			},
			want: &ListFilter[string]{
				comp: listNotContain,
				list: []string{"as", "df"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.constructor(tt.args.t...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("number constructor = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewListConditionWrite(t *testing.T) {
	type args struct {
		constructor func(t ...string) *ListFilter[string]
		t           []string
	}
	tests := []struct {
		name string
		args args
		want wantQuery
	}{
		{
			name: "ListEquals",
			args: args{
				constructor: NewListEquals[string],
				t:           []string{"as", "df"},
			},
			want: wantQuery{
				query: "test = $1",
				args:  []any{[]string{"as", "df"}},
			},
		},
		{
			name: "ListContains",
			args: args{
				constructor: NewListContains[string],
				t:           []string{"as", "df"},
			},
			want: wantQuery{
				query: "test = ANY($1)",
				args:  []any{[]string{"as", "df"}},
			},
		},
		{
			name: "ListNotContains",
			args: args{
				constructor: NewListNotContains[string],
				t:           []string{"as", "df"},
			},
			want: wantQuery{
				query: "NOT(test = ANY($1))",
				args:  []any{[]string{"as", "df"}},
			},
		},
		{
			name: "empty list",
			args: args{
				constructor: NewListNotContains[string],
			},
			want: wantQuery{
				query: "",
				args:  nil,
			},
		},
	}
	for _, tt := range tests {
		var stmt Statement
		t.Run(tt.name, func(t *testing.T) {
			tt.args.constructor(tt.args.t...).Write(&stmt, "test")
			assertQuery(t, &stmt, tt.want)
		})
	}
}
