package database

import (
	"reflect"
	"testing"
)

func TestStatement_WriteArgs(t *testing.T) {
	type args struct {
		args []any
	}
	tests := []struct {
		name string
		args args
		want wantQuery
	}{
		{
			name: "no args",
			args: args{
				args: nil,
			},
		},
		{
			name: "1 arg",
			args: args{
				args: []any{"asdf"},
			},
			want: wantQuery{
				query: "$1",
				args:  []any{"asdf"},
			},
		},
		{
			name: "n args",
			args: args{
				args: []any{"asdf", "jkl", 1},
			},
			want: wantQuery{
				query: "$1, $2, $3",
				args:  []any{"asdf", "jkl", 1},
			},
		},
	}
	for _, tt := range tests {
		var stmt Statement
		t.Run(tt.name, func(t *testing.T) {
			stmt.WriteArgs(tt.args.args...)
			assertQuery(t, &stmt, tt.want)
		})
	}
}

type wantQuery struct {
	query string
	args  []any
}

func assertQuery(t *testing.T, stmt *Statement, want wantQuery) {
	if want.query != stmt.String() {
		t.Errorf("unexpected query: want: %q got: %q", want.query, stmt.String())
	}

	if len(want.args) != len(stmt.Args()) {
		t.Errorf("unexpected length of args: want %d, got %d", len(want.args), len(stmt.Args()))
		return
	}

	for i, wantArg := range want.args {
		if !reflect.DeepEqual(wantArg, stmt.Args()[i]) {
			t.Errorf("unexpected arg at position %d: want: %v, got: %v", i, wantArg, stmt.Args()[i])
		}
	}
}
