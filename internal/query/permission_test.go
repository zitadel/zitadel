package query

import (
	"context"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	domain_pkg "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/feature"
)

func TestPermissionClause(t *testing.T) {
	var permissions = []authz.SystemUserPermissions{
		{
			MemberType:  authz.MemberTypeOrganization,
			AggregateID: "orgID",
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

	type args struct {
		ctx        context.Context
		query      sq.SelectBuilder
		enabled    bool
		orgIDCol   Column
		permission string
		options    []PermissionOption
	}
	tests := []struct {
		name      string
		args      args
		wantQuery string
		wantArgs  []any
	}{
		{
			name: "disabled",
			args: args{
				ctx:        ctx,
				query:      sq.Select("foo", "bar").From("users").Where(sq.Eq{"instance_id": "instanceID"}),
				enabled:    false,
				orgIDCol:   UserResourceOwnerCol,
				permission: "permission1",
			},
			wantQuery: "SELECT foo, bar FROM users WHERE instance_id = ?",
			wantArgs:  []any{"instanceID"},
		},
		{
			name: "no options",
			args: args{
				ctx:        ctx,
				query:      sq.Select("foo", "bar").From("users").Where(sq.Eq{"instance_id": "instanceID"}),
				enabled:    true,
				orgIDCol:   UserResourceOwnerCol,
				permission: "permission1",
			},
			wantQuery: "SELECT foo, bar FROM users WHERE instance_id = ? AND (projections.users14.resource_owner = ANY(eventstore.permitted_orgs(?, ?, ?, ?, ?)))",
			wantArgs:  []any{"instanceID", "instanceID", "userID", database.NewJSONArray(permissions), "permission1", ""},
		},
		{
			name: "owned rows option",
			args: args{
				ctx:        ctx,
				query:      sq.Select("foo", "bar").From("users").Where(sq.Eq{"instance_id": "instanceID"}),
				enabled:    true,
				orgIDCol:   UserResourceOwnerCol,
				permission: "permission1",
				options: []PermissionOption{
					OwnedRowsPermissionOption(UserIDCol),
				},
			},
			wantQuery: "SELECT foo, bar FROM users WHERE instance_id = ? AND (projections.users14.resource_owner = ANY(eventstore.permitted_orgs(?, ?, ?, ?, ?)) OR projections.users14.id = ?)",
			wantArgs:  []any{"instanceID", "instanceID", "userID", database.NewJSONArray(permissions), "permission1", "", "userID"},
		},
		{
			name: "override rows option",
			args: args{
				ctx:        ctx,
				query:      sq.Select("foo", "bar").From("users").Where(sq.Eq{"instance_id": "instanceID"}),
				enabled:    true,
				orgIDCol:   UserResourceOwnerCol,
				permission: "permission1",
				options: []PermissionOption{
					OwnedRowsPermissionOption(UserIDCol),
					OverridePermissionOption(UserStateCol, "bar"),
				},
			},
			wantQuery: "SELECT foo, bar FROM users WHERE instance_id = ? AND (projections.users14.resource_owner = ANY(eventstore.permitted_orgs(?, ?, ?, ?, ?)) OR projections.users14.id = ? OR projections.users14.state = ?)",
			wantArgs:  []any{"instanceID", "instanceID", "userID", database.NewJSONArray(permissions), "permission1", "", "userID", "bar"},
		},
		{
			name: "single org option",
			args: args{
				ctx:        ctx,
				query:      sq.Select("foo", "bar").From("users").Where(sq.Eq{"instance_id": "instanceID"}),
				enabled:    true,
				orgIDCol:   UserResourceOwnerCol,
				permission: "permission1",
				options: []PermissionOption{
					SingleOrgPermissionOption([]SearchQuery{
						mustSearchQuery(NewUserDisplayNameSearchQuery("zitadel", TextContains)),
						mustSearchQuery(NewUserResourceOwnerSearchQuery("orgID", TextEquals)),
					}),
				},
			},
			wantQuery: "SELECT foo, bar FROM users WHERE instance_id = ? AND (projections.users14.resource_owner = ANY(eventstore.permitted_orgs(?, ?, ?, ?, ?)))",
			wantArgs:  []any{"instanceID", "instanceID", "userID", database.NewJSONArray(permissions), "permission1", "orgID"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := PermissionClause(tt.args.ctx, tt.args.query, tt.args.enabled, tt.args.orgIDCol, tt.args.permission, tt.args.options...)
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

func TestPermissionV2(t *testing.T) {
	type args struct {
		ctx context.Context
		cf  domain_pkg.PermissionCheck
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "feature disabled, no permission check",
			args: args{
				ctx: context.Background(),
				cf:  nil,
			},
			want: false,
		},
		{
			name: "feature enabled, no permission check",
			args: args{
				ctx: authz.WithFeatures(context.Background(), feature.Features{
					PermissionCheckV2: true,
				}),
				cf: nil,
			},
			want: false,
		},
		{
			name: "feature enabled, with permission check",
			args: args{
				ctx: authz.WithFeatures(context.Background(), feature.Features{
					PermissionCheckV2: true,
				}),
				cf: func(context.Context, string, string, string) error {
					return nil
				},
			},
			want: true,
		},
		{
			name: "feature disabled, with permission check",
			args: args{
				ctx: authz.WithFeatures(context.Background(), feature.Features{
					PermissionCheckV2: false,
				}),
				cf: func(context.Context, string, string, string) error {
					return nil
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PermissionV2(tt.args.ctx, tt.args.cf)
			assert.Equal(t, tt.want, got)
		})
	}
}
