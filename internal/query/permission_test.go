package query

import (
	"context"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
)

func TestWherePermittedOrgs(t *testing.T) {
	const (
		instanceID = "instanceID"
		orgID      = "orgID"
		userID     = "userID"
		permission = "permission1"
	)
	var permissions = []authz.SystemUserPermissions{
		{
			MemberType:  authz.MemberTypeOrganization,
			AggregateID: orgID,
			Permissions: []string{"permission1", "permission2"},
		},
		{
			MemberType:  authz.MemberTypeIAM,
			Permissions: []string{"permission2", "permission3"},
		},
	}
	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	ctx = authz.SetCtxData(ctx, authz.CtxData{
		UserID:                "userID",
		SystemUserPermissions: permissions,
	})

	baseQuery := sq.Select("foo", "bar").From("users").Where(sq.Eq{"instance_id": instanceID})
	tests := []struct {
		name      string
		options   []PermittedOrgsOption
		wantQuery string
		wantArgs  []any
	}{
		{
			name:      "no options",
			wantQuery: "SELECT foo, bar FROM users WHERE instance_id = ? AND (projections.users14.resource_owner = ANY(eventstore.permitted_orgs(?, ?, ?, ?, ?)))",
			wantArgs:  []any{instanceID, instanceID, userID, database.NewJSONArray(permissions), permission, ""},
		},
		{
			name: "owned rows option",
			options: []PermittedOrgsOption{
				OwnedRowsOrgOption(UserIDCol),
			},
			wantQuery: "SELECT foo, bar FROM users WHERE instance_id = ? AND (projections.users14.resource_owner = ANY(eventstore.permitted_orgs(?, ?, ?, ?, ?)) OR projections.users14.id = ?)",
			wantArgs:  []any{instanceID, instanceID, userID, database.NewJSONArray(permissions), permission, "", userID},
		},
		{
			name: "override rows option",
			options: []PermittedOrgsOption{
				OwnedRowsOrgOption(UserIDCol),
				OverrideOrgOption(UserStateCol, "bar"),
			},
			wantQuery: "SELECT foo, bar FROM users WHERE instance_id = ? AND (projections.users14.resource_owner = ANY(eventstore.permitted_orgs(?, ?, ?, ?, ?)) OR projections.users14.id = ? OR projections.users14.state = ?)",
			wantArgs:  []any{instanceID, instanceID, userID, database.NewJSONArray(permissions), permission, "", userID, "bar"},
		},
		{
			name: "single org option",
			options: []PermittedOrgsOption{
				SingleOrgOption([]SearchQuery{
					mustSearchQuery(NewUserDisplayNameSearchQuery("zitadel", TextContains)),
					mustSearchQuery(NewUserResourceOwnerSearchQuery(orgID, TextEquals)),
				}),
			},
			wantQuery: "SELECT foo, bar FROM users WHERE instance_id = ? AND (projections.users14.resource_owner = ANY(eventstore.permitted_orgs(?, ?, ?, ?, ?)))",
			wantArgs:  []any{instanceID, instanceID, userID, database.NewJSONArray(permissions), permission, orgID},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := WherePermittedOrgs(ctx, baseQuery, UserResourceOwnerCol, permission, tt.options...)
			gotQuery, gotArgs, err := query.ToSql()
			require.NoError(t, err)
			assert.Equal(t, tt.wantQuery, gotQuery)
			assert.Equal(t, tt.wantArgs, gotArgs)
		})
	}
}

func mustSearchQuery(q SearchQuery, err error) SearchQuery {
	if err != nil {
		panic(err)
	}
	return q
}
