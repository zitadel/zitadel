package query

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GroupGrantsCheckPermission(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		want        []*GroupGrant
		grants      *GroupGrants
		permissions []string
	}{
		{
			name: "no permissions",
			want: []*GroupGrant{},
			grants: &GroupGrants{
				GroupGrants: []*GroupGrant{
					{ID: "grant1", GroupID: "group1"},
					{ID: "grant2", GroupID: "group2"},
					{ID: "grant3", GroupID: "group3"},
				},
			},
			permissions: []string{},
		},
		{
			name: "permissions for group1",
			want: []*GroupGrant{
				{ID: "grant1", GroupID: "group1"},
			},
			grants: &GroupGrants{
				GroupGrants: []*GroupGrant{
					{ID: "grant1", GroupID: "group1"},
					{ID: "grant2", GroupID: "group2"},
					{ID: "grant3", GroupID: "group3"},
				},
			},
			permissions: []string{"group1"},
		},
		{
			name: "permissions for multiple groups keeps order",
			want: []*GroupGrant{
				{ID: "grant1", GroupID: "group1"},
				{ID: "grant3", GroupID: "group3"},
			},
			grants: &GroupGrants{
				GroupGrants: []*GroupGrant{
					{ID: "grant1", GroupID: "group1"},
					{ID: "grant2", GroupID: "group2"},
					{ID: "grant3", GroupID: "group3"},
				},
			},
			permissions: []string{"group1", "group3"},
		},
		{
			name: "permissions for all groups keeps every grant",
			want: []*GroupGrant{
				{ID: "grant1", GroupID: "group1"},
				{ID: "grant2", GroupID: "group2"},
			},
			grants: &GroupGrants{
				GroupGrants: []*GroupGrant{
					{ID: "grant1", GroupID: "group1"},
					{ID: "grant2", GroupID: "group2"},
				},
			},
			permissions: []string{"group1", "group2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			checkPermission := func(ctx context.Context, permission, orgID, resourceID string) error {
				for _, perm := range tt.permissions {
					if resourceID == perm {
						return nil
					}
				}
				return errors.New("not found")
			}
			groupGrantsCheckPermission(context.Background(), tt.grants, checkPermission)
			require.Equal(t, tt.want, tt.grants.GroupGrants)
		})
	}
}
