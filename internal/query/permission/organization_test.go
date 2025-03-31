package permission

import (
	"context"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
)

const (
	instanceID = "instanceID"
	orgID      = "orgID"
	userID     = "userID"
	orgIDCol   = "resource_owner"
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

func prepareCtx() context.Context {
	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	return authz.SetCtxData(ctx, authz.CtxData{
		UserID:                "userID",
		SystemUserPermissions: permissions,
	})
}

func TestOrgsFilter(t *testing.T) {
	baseQuery := sq.Select("foo", "bar").From("users").Where(sq.Eq{"instance_id": instanceID})
	tests := []struct {
		name      string
		options   []OrgsOption
		wantQuery string
		wantArgs  []any
	}{
		{
			name:      "no options",
			wantQuery: "SELECT foo, bar FROM users WHERE instance_id = ? AND (resource_owner = ANY(eventstore.permitted_orgs(?, ?, ?, ?, ?)))",
			wantArgs:  []any{instanceID, instanceID, userID, database.NewJSONArray(permissions), permission, orgID},
		},
		{
			name: "owned rows option",
			options: []OrgsOption{
				OwnedRowsOption("user_id"),
			},
			wantQuery: "SELECT foo, bar FROM users WHERE instance_id = ? AND (resource_owner = ANY(eventstore.permitted_orgs(?, ?, ?, ?, ?)) OR user_id = ?)",
			wantArgs:  []any{instanceID, instanceID, userID, database.NewJSONArray(permissions), permission, orgID, userID},
		},
		{
			name: "override rows option",
			options: []OrgsOption{
				OwnedRowsOption("user_id"),
				OverrideOption("foo", "bar"),
			},
			wantQuery: "SELECT foo, bar FROM users WHERE instance_id = ? AND (resource_owner = ANY(eventstore.permitted_orgs(?, ?, ?, ?, ?)) OR user_id = ? OR foo = ?)",
			wantArgs:  []any{instanceID, instanceID, userID, database.NewJSONArray(permissions), permission, orgID, userID, "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := prepareCtx()
			query := OrgsFilter(ctx, baseQuery, orgIDCol, orgID, permission, tt.options...)
			gotQuery, gotArgs, err := query.ToSql()
			require.NoError(t, err)
			assert.Equal(t, tt.wantQuery, gotQuery)
			assert.Equal(t, tt.wantArgs, gotArgs)
		})
	}
}
