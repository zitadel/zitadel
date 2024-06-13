package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleOrgIDsFromScope(t *testing.T) {
	type args struct {
		scopes []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "nil",
			args: args{nil},
			want: nil,
		},
		{
			name: "unrelated scope",
			args: args{[]string{"foo", "bar"}},
			want: nil,
		},
		{
			name: "orgID role scope",
			args: args{[]string{OrgRoleIDScope + "123"}},
			want: []string{"123"},
		},
		{
			name: "mixed scope",
			args: args{[]string{"foo", OrgRoleIDScope + "123"}},
			want: []string{"123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RoleOrgIDsFromScope(tt.args.scopes)
			assert.Equal(t, tt.want, got)
		})
	}
}
