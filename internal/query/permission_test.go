package query

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

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
		orgIDCol   Column
		permission string
		options    []PermissionOption
	}
	tests := []struct {
		name     string
		args     args
		wantSql  string
		wantArgs []any
	}{
		{
			name: "org, no options",
			args: args{
				ctx:        ctx,
				orgIDCol:   UserResourceOwnerCol,
				permission: "permission1",
			},
			wantSql: "INNER JOIN eventstore.permitted_orgs(?, ?, ?, ?, ?) permissions ON (permissions.instance_permitted OR projections.users14.resource_owner = ANY(permissions.org_ids))",
			wantArgs: []any{
				"instanceID",
				"userID",
				database.NewJSONArray(permissions),
				"permission1",
				(*string)(nil),
			},
		},
		{
			name: "org, owned rows option",
			args: args{
				ctx:        ctx,
				orgIDCol:   UserResourceOwnerCol,
				permission: "permission1",
				options: []PermissionOption{
					OwnedRowsPermissionOption(UserIDCol),
				},
			},
			wantSql: "INNER JOIN eventstore.permitted_orgs(?, ?, ?, ?, ?) permissions ON (permissions.instance_permitted OR projections.users14.resource_owner = ANY(permissions.org_ids) OR projections.users14.id = ?)",
			wantArgs: []any{
				"instanceID",
				"userID",
				database.NewJSONArray(permissions),
				"permission1",
				(*string)(nil),
				"userID",
			},
		},
		{
			name: "org, connection rows option",
			args: args{
				ctx:        ctx,
				orgIDCol:   UserResourceOwnerCol,
				permission: "permission1",
				options: []PermissionOption{
					OwnedRowsPermissionOption(UserIDCol),
					ConnectionPermissionOption(UserStateCol, "bar"),
				},
			},
			wantSql: "INNER JOIN eventstore.permitted_orgs(?, ?, ?, ?, ?) permissions ON (permissions.instance_permitted OR projections.users14.resource_owner = ANY(permissions.org_ids) OR projections.users14.id = ? OR projections.users14.state = ?)",
			wantArgs: []any{
				"instanceID",
				"userID",
				database.NewJSONArray(permissions),
				"permission1",
				(*string)(nil),
				"userID",
				"bar",
			},
		},
		{
			name: "org, with ID",
			args: args{
				ctx:        ctx,
				orgIDCol:   UserResourceOwnerCol,
				permission: "permission1",
				options: []PermissionOption{
					SingleOrgPermissionOption([]SearchQuery{
						mustSearchQuery(NewUserDisplayNameSearchQuery("zitadel", TextContains)),
						mustSearchQuery(NewUserResourceOwnerSearchQuery("orgID", TextEquals)),
					}),
				},
			},
			wantSql: "INNER JOIN eventstore.permitted_orgs(?, ?, ?, ?, ?) permissions ON (permissions.instance_permitted OR projections.users14.resource_owner = ANY(permissions.org_ids))",
			wantArgs: []any{
				"instanceID",
				"userID",
				database.NewJSONArray(permissions),
				"permission1",
				gu.Ptr("orgID"),
			},
		},
		{
			name: "project",
			args: args{
				ctx:        ctx,
				orgIDCol:   ProjectColumnResourceOwner,
				permission: "permission1",
				options: []PermissionOption{
					WithProjectsPermissionOption(ProjectColumnID),
				},
			},
			wantSql: "INNER JOIN eventstore.permitted_projects(?, ?, ?, ?, ?) permissions ON (permissions.instance_permitted OR projections.projects4.resource_owner = ANY(permissions.org_ids) OR projections.projects4.id = ANY(permissions.project_ids))",
			wantArgs: []any{
				"instanceID",
				"userID",
				database.NewJSONArray(permissions),
				"permission1",
				(*string)(nil),
			},
		},
		{
			name: "project, single org",
			args: args{
				ctx:        ctx,
				orgIDCol:   ProjectColumnResourceOwner,
				permission: "permission1",
				options: []PermissionOption{
					WithProjectsPermissionOption(ProjectColumnID),
					SingleOrgPermissionOption([]SearchQuery{
						mustSearchQuery(NewProjectResourceOwnerSearchQuery("orgID")),
					}),
				},
			},
			wantSql: "INNER JOIN eventstore.permitted_projects(?, ?, ?, ?, ?) permissions ON (permissions.instance_permitted OR projections.projects4.resource_owner = ANY(permissions.org_ids) OR projections.projects4.id = ANY(permissions.project_ids))",
			wantArgs: []any{
				"instanceID",
				"userID",
				database.NewJSONArray(permissions),
				"permission1",
				gu.Ptr("orgID"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSql, gotArgs := PermissionClause(tt.args.ctx, tt.args.orgIDCol, tt.args.permission, tt.args.options...)
			assert.Equal(t, tt.wantSql, gotSql)
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
