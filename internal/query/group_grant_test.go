package query

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
)

// Test_groupGrantPermissionCheckV2 locks down the v2 SQL shape: enabling the
// permission check must add a well-formed INNER JOIN against the permitted_orgs
// function so a caller with zero accessible orgs returns an empty result via
// the join — not via a degenerate WHERE clause that would be a SQL syntax error.
func Test_groupGrantPermissionCheckV2(t *testing.T) {
	t.Parallel()

	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	ctx = authz.SetCtxData(ctx, authz.CtxData{UserID: "userID"})

	t.Run("disabled returns the base query untouched", func(t *testing.T) {
		t.Parallel()
		base, _ := prepareGroupGrantsQuery()
		got := groupGrantPermissionCheckV2(ctx, base, false)
		gotSQL, _, err := got.ToSql()
		require.NoError(t, err)
		baseSQL, _, err := base.ToSql()
		require.NoError(t, err)
		assert.Equal(t, baseSQL, gotSQL)
	})

	t.Run("enabled appends a well-formed permitted_orgs join", func(t *testing.T) {
		t.Parallel()
		base, _ := prepareGroupGrantsQuery()
		got := groupGrantPermissionCheckV2(ctx, base, true)
		gotSQL, args, err := got.ToSql()
		require.NoError(t, err)
		assert.Contains(t, gotSQL, "INNER JOIN eventstore.permitted_orgs")
		assert.Contains(t, gotSQL, "permissions.instance_permitted")
		assert.Contains(t, gotSQL, "ANY(permissions.org_ids)")
		// the dangerous regression: an empty IN list would render as IN ()
		// and postgres rejects it as a syntax error
		assert.False(t, strings.Contains(gotSQL, "IN ()"), "no empty IN clause: %q", gotSQL)
		// PermissionClause appends 5 join args: instance, user, system perms, permission, org filter
		require.GreaterOrEqual(t, len(args), 5)
		assert.Equal(t, "instanceID", args[0])
		assert.Equal(t, "userID", args[1])
	})
}

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
