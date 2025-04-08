//go:build integration

package setup_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
)

const ConnString = "host=localhost port=5432 user=zitadel dbname=zitadel sslmode=disable"

var (
	CTX    context.Context
	dbPool *pgxpool.Pool
)

func TestMain(m *testing.M) {
	var cancel context.CancelFunc
	CTX, cancel = context.WithTimeout(context.Background(), time.Second*10)

	var err error
	dbPool, err = pgxpool.New(context.Background(), ConnString)
	if err != nil {
		panic(err)
	}
	exit := m.Run()
	cancel()
	dbPool.Close()
	os.Exit(exit)
}

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
