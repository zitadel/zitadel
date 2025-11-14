//go:build integration

package setup_test

import (
	"encoding/json"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/permission"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func TestGetSystemPermissions(t *testing.T) {
	const query = "SELECT * FROM eventstore.get_system_permissions($1, $2);"
	t.Parallel()
	permissions := []authz.SystemUserPermissions{
		{
			MemberType:  authz.MemberTypeSystem,
			Permissions: []string{"iam.read", "iam.write", "iam.policy.read"},
		},
		{
			MemberType:  authz.MemberTypeIAM,
			AggregateID: "instanceID",
			Permissions: []string{"iam.read", "iam.write", "iam.policy.read", "org.read", "project.read", "project.write"},
		},
		{
			MemberType:  authz.MemberTypeOrganization,
			AggregateID: "orgID",
			Permissions: []string{"org.read", "org.write", "org.policy.read", "project.read", "project.write"},
		},
		{
			MemberType:  authz.MemberTypeProject,
			AggregateID: "projectID",
			Permissions: []string{"project.read", "project.write"},
		},
		{
			MemberType:  authz.MemberTypeProjectGrant,
			AggregateID: "projectID",
			ObjectID:    "grantID",
			Permissions: []string{"project.read", "project.write"},
		},
	}
	type result struct {
		MemberType  authz.MemberType
		AggregateID string
		ObjectID    string
	}
	tests := []struct {
		permm string
		want  []result
	}{
		{
			permm: "iam.read",
			want: []result{
				{
					MemberType: authz.MemberTypeSystem,
				},
				{
					MemberType:  authz.MemberTypeIAM,
					AggregateID: "instanceID",
				},
			},
		},
		{
			permm: "org.read",
			want: []result{
				{
					MemberType:  authz.MemberTypeIAM,
					AggregateID: "instanceID",
				},
				{
					MemberType:  authz.MemberTypeOrganization,
					AggregateID: "orgID",
				},
			},
		},
		{
			permm: "project.write",
			want: []result{
				{
					MemberType:  authz.MemberTypeIAM,
					AggregateID: "instanceID",
				},
				{
					MemberType:  authz.MemberTypeOrganization,
					AggregateID: "orgID",
				},
				{
					MemberType:  authz.MemberTypeProject,
					AggregateID: "projectID",
				},
				{
					MemberType:  authz.MemberTypeProjectGrant,
					AggregateID: "projectID",
					ObjectID:    "grantID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.permm, func(t *testing.T) {
			t.Parallel()
			rows, err := dbPool.Query(CTX, query, database.NewJSONArray(permissions), tt.permm)
			require.NoError(t, err)
			got, err := pgx.CollectRows(rows, pgx.RowToStructByPos[result])
			require.NoError(t, err)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestCheckSystemUserPerms(t *testing.T) {
	// Use JSON because of the composite project_grants SQL type
	const query = "SELECT row_to_json(eventstore.check_system_user_perms($1, $2, $3));"
	t.Parallel()
	type args struct {
		reqInstanceID string
		permissions   []authz.SystemUserPermissions
		permm         string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "iam.read, instance permitted from system",
			args: args{
				reqInstanceID: "instanceID",
				permissions: []authz.SystemUserPermissions{
					{
						MemberType:  authz.MemberTypeSystem,
						Permissions: []string{"iam.read", "iam.write", "iam.policy.read"},
					},
					{
						MemberType:  authz.MemberTypeIAM,
						AggregateID: "instanceID",
						Permissions: []string{"iam.read", "iam.write", "iam.policy.read", "org.read"},
					},
					{
						MemberType:  authz.MemberTypeOrganization,
						AggregateID: "orgID",
						Permissions: []string{"org.read", "org.write", "org.policy.read", "project.read", "project.write"},
					},
					{
						MemberType:  authz.MemberTypeProject,
						AggregateID: "projectID",
						Permissions: []string{"project.read", "project.write"},
					},
					{
						MemberType:  authz.MemberTypeProjectGrant,
						AggregateID: "projectID",
						ObjectID:    "grantID",
						Permissions: []string{"project.read", "project.write"},
					},
				},
				permm: "iam.read",
			},
			want: `{
				"instance_permitted": true,
				"org_ids": [],
				"project_grants": [],
				"project_ids": []
			}`,
		},
		{
			name: "org.read, instance permitted",
			args: args{
				reqInstanceID: "instanceID",
				permissions: []authz.SystemUserPermissions{
					{
						MemberType:  authz.MemberTypeSystem,
						Permissions: []string{"iam.read", "iam.write", "iam.policy.read"},
					},
					{
						MemberType:  authz.MemberTypeIAM,
						AggregateID: "instanceID",
						Permissions: []string{"iam.read", "iam.write", "iam.policy.read", "org.read"},
					},
					{
						MemberType:  authz.MemberTypeOrganization,
						AggregateID: "orgID",
						Permissions: []string{"org.read", "org.write", "org.policy.read", "project.read", "project.write"},
					},
					{
						MemberType:  authz.MemberTypeProject,
						AggregateID: "projectID",
						Permissions: []string{"project.read", "project.write"},
					},
					{
						MemberType:  authz.MemberTypeProjectGrant,
						AggregateID: "projectID",
						ObjectID:    "grantID",
						Permissions: []string{"project.read", "project.write"},
					},
				},
				permm: "org.read",
			},
			want: `{
				"instance_permitted": true,
				"org_ids": [],
				"project_grants": [],
				"project_ids": []
			}`,
		},
		{
			name: "project.read, org ID and project ID permitted",
			args: args{
				reqInstanceID: "instanceID",
				permissions: []authz.SystemUserPermissions{
					{
						MemberType:  authz.MemberTypeSystem,
						Permissions: []string{"iam.read", "iam.write", "iam.policy.read"},
					},
					{
						MemberType:  authz.MemberTypeIAM,
						AggregateID: "instanceID",
						Permissions: []string{"iam.read", "iam.write", "iam.policy.read", "org.read"},
					},
					{
						MemberType:  authz.MemberTypeOrganization,
						AggregateID: "orgID",
						Permissions: []string{"org.read", "org.write", "org.policy.read", "project.read", "project.write"},
					},
					{
						MemberType:  authz.MemberTypeProject,
						AggregateID: "projectID",
						Permissions: []string{"project.read", "project.write"},
					},
					{
						MemberType:  authz.MemberTypeProjectGrant,
						AggregateID: "projectID",
						ObjectID:    "grantID",
						Permissions: []string{"project_grant.read", "project_grant.write"},
					},
				},
				permm: "project.read",
			},
			want: `{
				"instance_permitted": false,
				"org_ids": ["orgID"],
				"project_ids": ["projectID"],
				"project_grants": []
			}`,
		},
		{
			name: "project_grant.read, project grant ID permitted",
			args: args{
				reqInstanceID: "instanceID",
				permissions: []authz.SystemUserPermissions{
					{
						MemberType:  authz.MemberTypeSystem,
						Permissions: []string{"iam.read", "iam.write", "iam.policy.read"},
					},
					{
						MemberType:  authz.MemberTypeIAM,
						AggregateID: "instanceID",
						Permissions: []string{"iam.read", "iam.write", "iam.policy.read", "org.read"},
					},
					{
						MemberType:  authz.MemberTypeOrganization,
						AggregateID: "orgID",
						Permissions: []string{"org.read", "org.write", "org.policy.read", "project.read", "project.write"},
					},
					{
						MemberType:  authz.MemberTypeProject,
						AggregateID: "projectID",
						Permissions: []string{"project.read", "project.write"},
					},
					{
						MemberType:  authz.MemberTypeProjectGrant,
						AggregateID: "projectID",
						ObjectID:    "grantID",
						Permissions: []string{"project_grant.read", "project_grant.write"},
					},
				},
				permm: "project_grant.read",
			},
			want: `{
				"instance_permitted": false,
				"org_ids": [],
				"project_ids": [],
				"project_grants": [
					{
						"project_id": "projectID",
						"grant_id": "grantID"
					}
				]
			}`,
		},
		{
			name: "instance without aggregate ID",
			args: args{
				reqInstanceID: "instanceID",
				permissions: []authz.SystemUserPermissions{
					{
						MemberType:  authz.MemberTypeIAM,
						AggregateID: "",
						Permissions: []string{"foo.bar", "bar.foo"},
					},
				},
				permm: "foo.bar",
			},
			want: `{
				"instance_permitted": false,
				"org_ids": [],
				"project_ids": [],
				"project_grants": []
			}`,
		},
		{
			name: "wrong instance ID",
			args: args{
				reqInstanceID: "instanceID",
				permissions: []authz.SystemUserPermissions{
					{
						MemberType:  authz.MemberTypeIAM,
						AggregateID: "wrong",
						Permissions: []string{"foo.bar", "bar.foo"},
					},
				},
				permm: "foo.bar",
			},
			want: `{
				"instance_permitted": false,
				"org_ids": [],
				"project_ids": [],
				"project_grants": []
			}`,
		},
		{
			name: "permission on other instance",
			args: args{
				reqInstanceID: "instanceID",
				permissions: []authz.SystemUserPermissions{
					{
						MemberType:  authz.MemberTypeIAM,
						AggregateID: "instanceID",
						Permissions: []string{"bar.foo"},
					},
					{
						MemberType:  authz.MemberTypeIAM,
						AggregateID: "wrong",
						Permissions: []string{"foo.bar"},
					},
				},
				permm: "foo.bar",
			},
			want: `{
				"instance_permitted": false,
				"org_ids": [],
				"project_ids": [],
				"project_grants": []
			}`,
		},
		{
			name: "org ID missing",
			args: args{
				reqInstanceID: "instanceID",
				permissions: []authz.SystemUserPermissions{
					{
						MemberType:  authz.MemberTypeOrganization,
						AggregateID: "",
						Permissions: []string{"foo.bar"},
					},
				},
				permm: "foo.bar",
			},
			want: `{
				"instance_permitted": false,
				"org_ids": [],
				"project_ids": [],
				"project_grants": []
			}`,
		},
		{
			name: "multiple org IDs",
			args: args{
				reqInstanceID: "instanceID",
				permissions: []authz.SystemUserPermissions{
					{
						MemberType:  authz.MemberTypeOrganization,
						AggregateID: "Org1",
						Permissions: []string{"foo.bar"},
					},
					{
						MemberType:  authz.MemberTypeOrganization,
						AggregateID: "Org2",
						Permissions: []string{"foo.bar"},
					},
				},
				permm: "foo.bar",
			},
			want: `{
				"instance_permitted": false,
				"org_ids": ["Org1", "Org2"],
				"project_ids": [],
				"project_grants": []
			}`,
		},
		{
			name: "project ID missing",
			args: args{
				reqInstanceID: "instanceID",
				permissions: []authz.SystemUserPermissions{
					{
						MemberType:  authz.MemberTypeProject,
						AggregateID: "",
						Permissions: []string{"foo.bar"},
					},
				},
				permm: "foo.bar",
			},
			want: `{
				"instance_permitted": false,
				"org_ids": [],
				"project_ids": [],
				"project_grants": []
			}`,
		},
		{
			name: "multiple project IDs",
			args: args{
				reqInstanceID: "instanceID",
				permissions: []authz.SystemUserPermissions{
					{
						MemberType:  authz.MemberTypeProject,
						AggregateID: "P1",
						Permissions: []string{"foo.bar"},
					},
					{
						MemberType:  authz.MemberTypeProject,
						AggregateID: "P2",
						Permissions: []string{"foo.bar"},
					},
				},
				permm: "foo.bar",
			},
			want: `{
				"instance_permitted": false,
				"org_ids": [],
				"project_ids": ["P1", "P2"],
				"project_grants": []
			}`,
		},
		{
			name: "project grant ID missing",
			args: args{
				reqInstanceID: "instanceID",
				permissions: []authz.SystemUserPermissions{
					{
						MemberType:  authz.MemberTypeProjectGrant,
						AggregateID: "",
						ObjectID:    "",
						Permissions: []string{"foo.bar"},
					},
				},
				permm: "foo.bar",
			},
			want: `{
				"instance_permitted": false,
				"org_ids": [],
				"project_ids": [],
				"project_grants": []
			}`,
		},
		{
			name: "multiple project IDs",
			args: args{
				reqInstanceID: "instanceID",
				permissions: []authz.SystemUserPermissions{
					{
						MemberType:  authz.MemberTypeProjectGrant,
						AggregateID: "P1",
						ObjectID:    "O1",
						Permissions: []string{"foo.bar"},
					},
					{
						MemberType:  authz.MemberTypeProjectGrant,
						AggregateID: "P2",
						ObjectID:    "O2",
						Permissions: []string{"foo.bar"},
					},
				},
				permm: "foo.bar",
			},
			want: `{
				"instance_permitted": false,
				"org_ids": [],
				"project_ids": [],
				"project_grants": [
					{
						"project_id": "P1",
						"grant_id": "O1"
					},
					{
						"project_id": "P2",
						"grant_id": "O2"
					}
				]
			}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rows, err := dbPool.Query(CTX, query, database.NewJSONArray(tt.args.permissions), tt.args.reqInstanceID, tt.args.permm)
			require.NoError(t, err)
			got, err := pgx.CollectOneRow(rows, pgx.RowTo[string])
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, got)
		})
	}
}

const (
	instanceID = "instanceID"
	orgID      = "orgID"
	projectID  = "projectID"
)

func TestPermittedOrgs(t *testing.T) {
	t.Parallel()

	tx, err := dbPool.Begin(CTX)
	require.NoError(t, err)
	defer tx.Rollback(CTX)

	// Insert a couple of deterministic field rows to test the function.
	// Data will not persist, because the transaction is rolled back.
	createRolePermission(t, tx, "IAM_OWNER", []string{"org.write", "org.read"})
	createRolePermission(t, tx, "ORG_OWNER", []string{"org.write", "org.read"})
	createMember(t, tx, instance.AggregateType, "instance_user")
	createMember(t, tx, org.AggregateType, "org_user")

	const query = "SELECT instance_permitted, org_ids FROM eventstore.permitted_orgs($1,$2,$3,$4,$5);"
	type args struct {
		reqInstanceID   string
		authUserID      string
		systemUserPerms []authz.SystemUserPermissions
		perm            string
		filterOrg       *string
	}
	type result struct {
		InstancePermitted bool
		OrgIDs            pgtype.FlatArray[string]
	}
	tests := []struct {
		name string
		args args
		want result
	}{
		{
			name: "system user, instance",
			args: args{
				reqInstanceID: instanceID,
				systemUserPerms: []authz.SystemUserPermissions{{
					MemberType:  authz.MemberTypeSystem,
					Permissions: []string{"org.write", "org.read"},
				}},
				perm: "org.read",
			},
			want: result{
				InstancePermitted: true,
			},
		},
		{
			name: "system user, orgs",
			args: args{
				reqInstanceID: instanceID,
				systemUserPerms: []authz.SystemUserPermissions{{
					MemberType:  authz.MemberTypeOrganization,
					AggregateID: orgID,
					Permissions: []string{"org.read", "org.write", "org.policy.read", "project.read", "project.write"},
				}},
				perm: "org.read",
			},
			want: result{
				OrgIDs: pgtype.FlatArray[string]{orgID},
			},
		},
		{
			name: "instance member",
			args: args{
				reqInstanceID: instanceID,
				authUserID:    "instance_user",
				perm:          "org.read",
			},
			want: result{
				InstancePermitted: true,
			},
		},
		{
			name: "org member",
			args: args{
				reqInstanceID: instanceID,
				authUserID:    "org_user",
				perm:          "org.read",
			},
			want: result{
				OrgIDs: pgtype.FlatArray[string]{orgID},
			},
		},
		{
			name: "org member, filter",
			args: args{
				reqInstanceID: instanceID,
				authUserID:    "org_user",
				perm:          "org.read",
				filterOrg:     gu.Ptr(orgID),
			},
			want: result{
				OrgIDs: pgtype.FlatArray[string]{orgID},
			},
		},
		{
			name: "org member, filter wrong org",
			args: args{
				reqInstanceID: instanceID,
				authUserID:    "org_user",
				perm:          "org.read",
				filterOrg:     gu.Ptr("foobar"),
			},
			want: result{},
		},
		{
			name: "no permission",
			args: args{
				reqInstanceID: instanceID,
				authUserID:    "foobar",
				perm:          "org.read",
				filterOrg:     gu.Ptr(orgID),
			},
			want: result{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := tx.Query(CTX, query, tt.args.reqInstanceID, tt.args.authUserID, database.NewJSONArray(tt.args.systemUserPerms), tt.args.perm, tt.args.filterOrg)
			require.NoError(t, err)
			got, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[result])
			require.NoError(t, err)
			assert.Equal(t, tt.want.InstancePermitted, got.InstancePermitted)
			assert.ElementsMatch(t, tt.want.OrgIDs, got.OrgIDs)
		})
	}
}

func TestPermittedProjects(t *testing.T) {
	t.Parallel()

	tx, err := dbPool.Begin(CTX)
	require.NoError(t, err)
	defer tx.Rollback(CTX)

	// Insert a couple of deterministic field rows to test the function.
	// Data will not persist, because the transaction is rolled back.
	createRolePermission(t, tx, "IAM_OWNER", []string{"project.write", "project.read"})
	createRolePermission(t, tx, "ORG_OWNER", []string{"project.write", "project.read"})
	createRolePermission(t, tx, "PROJECT_OWNER", []string{"project.write", "project.read"})
	createMember(t, tx, instance.AggregateType, "instance_user")
	createMember(t, tx, org.AggregateType, "org_user")
	createMember(t, tx, project.AggregateType, "project_user")

	const query = "SELECT instance_permitted, org_ids, project_ids FROM eventstore.permitted_projects($1,$2,$3,$4,$5);"
	type args struct {
		reqInstanceID   string
		authUserID      string
		systemUserPerms []authz.SystemUserPermissions
		perm            string
		filterOrg       *string
	}
	type result struct {
		InstancePermitted bool
		OrgIDs            pgtype.FlatArray[string]
		ProjectIDs        pgtype.FlatArray[string]
	}
	tests := []struct {
		name string
		args args
		want result
	}{
		{
			name: "system user, instance",
			args: args{
				reqInstanceID: instanceID,
				systemUserPerms: []authz.SystemUserPermissions{{
					MemberType:  authz.MemberTypeSystem,
					Permissions: []string{"project.write", "project.read"},
				}},
				perm: "project.read",
			},
			want: result{
				InstancePermitted: true,
			},
		},
		{
			name: "system user, orgs",
			args: args{
				reqInstanceID: instanceID,
				systemUserPerms: []authz.SystemUserPermissions{{
					MemberType:  authz.MemberTypeOrganization,
					AggregateID: orgID,
					Permissions: []string{"project.read", "project.write"},
				}},
				perm: "project.read",
			},
			want: result{
				OrgIDs: pgtype.FlatArray[string]{orgID},
			},
		},
		{
			name: "system user, projects",
			args: args{
				reqInstanceID: instanceID,
				systemUserPerms: []authz.SystemUserPermissions{{
					MemberType:  authz.MemberTypeProject,
					AggregateID: projectID,
					Permissions: []string{"project.read", "project.write"},
				}},
				perm: "project.read",
			},
			want: result{
				ProjectIDs: pgtype.FlatArray[string]{projectID},
			},
		},
		{
			name: "system user, org and project",
			args: args{
				reqInstanceID: instanceID,
				systemUserPerms: []authz.SystemUserPermissions{
					{
						MemberType:  authz.MemberTypeOrganization,
						AggregateID: orgID,
						Permissions: []string{"project.read", "project.write"},
					},
					{
						MemberType:  authz.MemberTypeProject,
						AggregateID: projectID,
						Permissions: []string{"project.read", "project.write"},
					},
				},
				perm: "project.read",
			},
			want: result{
				OrgIDs:     pgtype.FlatArray[string]{orgID},
				ProjectIDs: pgtype.FlatArray[string]{projectID},
			},
		},
		{
			name: "instance member",
			args: args{
				reqInstanceID: instanceID,
				authUserID:    "instance_user",
				perm:          "project.read",
			},
			want: result{
				InstancePermitted: true,
			},
		},
		{
			name: "org member",
			args: args{
				reqInstanceID: instanceID,
				authUserID:    "org_user",
				perm:          "project.read",
			},
			want: result{
				InstancePermitted: false,
				OrgIDs:            pgtype.FlatArray[string]{orgID},
			},
		},
		{
			name: "org member, filter",
			args: args{
				reqInstanceID: instanceID,
				authUserID:    "org_user",
				perm:          "project.read",
				filterOrg:     gu.Ptr(orgID),
			},
			want: result{
				InstancePermitted: false,
				OrgIDs:            pgtype.FlatArray[string]{orgID},
			},
		},
		{
			name: "org member, filter wrong org",
			args: args{
				reqInstanceID: instanceID,
				authUserID:    "org_user",
				perm:          "project.read",
				filterOrg:     gu.Ptr("foobar"),
			},
			want: result{},
		},
		{
			name: "project member",
			args: args{
				reqInstanceID: instanceID,
				authUserID:    "project_user",
				perm:          "project.read",
			},
			want: result{
				ProjectIDs: pgtype.FlatArray[string]{projectID},
			},
		},
		{
			name: "no permission",
			args: args{
				reqInstanceID: instanceID,
				authUserID:    "foobar",
				perm:          "project.read",
				filterOrg:     gu.Ptr(orgID),
			},
			want: result{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := tx.Query(CTX, query, tt.args.reqInstanceID, tt.args.authUserID, database.NewJSONArray(tt.args.systemUserPerms), tt.args.perm, tt.args.filterOrg)
			require.NoError(t, err)
			got, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[result])
			require.NoError(t, err)
			assert.Equal(t, tt.want.InstancePermitted, got.InstancePermitted)
			assert.ElementsMatch(t, tt.want.OrgIDs, got.OrgIDs)
		})
	}
}

func createRolePermission(t *testing.T, tx pgx.Tx, role string, permissions []string) {
	for _, perm := range permissions {
		createTestField(t, tx, instanceID, permission.AggregateType, instanceID, "role_permission", role, "permission", perm)
	}
}

func createMember(t *testing.T, tx pgx.Tx, aggregateType eventstore.AggregateType, userID string) {
	var err error
	switch aggregateType {
	case instance.AggregateType:
		createTestField(t, tx, instanceID, aggregateType, instanceID, "instance_member_role", userID, "instance_role", "IAM_OWNER")
	case org.AggregateType:
		createTestField(t, tx, orgID, aggregateType, orgID, "org_member_role", userID, "org_role", "ORG_OWNER")
	case project.AggregateType:
		createTestField(t, tx, orgID, aggregateType, orgID, "project_member_role", userID, "project_role", "PROJECT_OWNER")
	default:
		panic("unknown aggregate type " + aggregateType)
	}
	require.NoError(t, err)
}

func createTestField(t *testing.T, tx pgx.Tx, resourceOwner string, aggregateType eventstore.AggregateType, aggregateID, objectType, objectID, fieldName string, value any) {
	const query = `INSERT INTO eventstore.fields(
		instance_id, resource_owner, aggregate_type, aggregate_id, object_type, object_id, field_name, value, value_must_be_unique, should_index, object_revision)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, false, true, 1);`
	encValue, err := json.Marshal(value)
	require.NoError(t, err)
	_, err = tx.Exec(CTX, query, instanceID, resourceOwner, aggregateType, aggregateID, objectType, objectID, fieldName, encValue)
	require.NoError(t, err)

}
