package migration_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/embedded"
)

func TestMigrate(t *testing.T) {
	tests := []struct {
		name string
		stmt string
		args []any
		res  []any
	}{
		{
			name: "schema",
			stmt: "SELECT EXISTS(SELECT 1 FROM information_schema.schemata where schema_name = 'zitadel') ;",
			res:  []any{true},
		},
		{
			name: "001",
			stmt: "SELECT EXISTS(SELECT 1 FROM pg_catalog.pg_tables WHERE schemaname = 'zitadel' and tablename=$1)",
			args: []any{"instances"},
			res:  []any{true},
		},
	}

	ctx := context.Background()

	connector, stop, err := embedded.StartEmbedded()
	require.NoError(t, err, "failed to start embedded postgres")
	defer stop()

	client, err := connector.Connect(ctx)
	require.NoError(t, err, "failed to connect to embedded postgres")

	err = client.(database.Migrator).Migrate(ctx)
	require.NoError(t, err, "failed to execute migration steps")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]any, len(tt.res))
			for i := range got {
				got[i] = new(any)
				tt.res[i] = gu.Ptr(tt.res[i])
			}

			require.NoError(t, client.QueryRow(ctx, tt.stmt, tt.args...).Scan(got...), "failed to execute check query")

			assert.Equal(t, tt.res, got, "query result does not match")
		})
	}

	_, err = client.Exec(ctx, "INSERT INTO zitadel.instances(id, name) VALUES ('i1', 'instance-1')")
	require.NoError(t, err)

	_, err = client.Exec(ctx, "INSERT INTO zitadel.organizations(instance_id, id, name, state) VALUES ('i1', 'o1', 'org-1', 'active'), ('i1', 'o2', 'org-2', 'active')")
	require.NoError(t, err)

	_, err = client.Exec(ctx, "INSERT INTO zitadel.projects(instance_id, organization_id, id, name, state) VALUES ('i1', 'o1', 'p1', 'project-1', 'active')")
	require.NoError(t, err)

	_, err = client.Exec(ctx, "INSERT INTO zitadel.project_grants(instance_id, id, granting_organization_id, project_id, granted_organization_id, state) VALUES ('i1', 'pg1', 'o1', 'p1', 'o2', 'active')")
	require.NoError(t, err)

	_, err = client.Exec(ctx, "INSERT INTO zitadel.users(instance_id, organization_id, id, username, type) VALUES ('i1', 'o1', 'u_instance', 'user-instance', 'machine'), ('i1', 'o1', 'u_org', 'user-org', 'machine'), ('i1', 'o1', 'u_project', 'user-project', 'machine'), ('i1', 'o1', 'u_project_grant', 'user-project-grant', 'machine'), ('i1', 'o1', 'u_none', 'user-none', 'machine')")
	require.NoError(t, err)

	_, err = client.Exec(ctx, "INSERT INTO zitadel.administrator_role_permissions(instance_id, role_name, permission) VALUES ('i1', 'instance_admin', 'perm.read'), ('i1', 'instance_admin', 'perm.write'), ('i1', 'org_admin', 'perm.read'), ('i1', 'project_admin', 'perm.read'), ('i1', 'project_grant_admin', 'perm.read')")
	require.NoError(t, err)

	_, err = client.Exec(ctx, "INSERT INTO zitadel.administrators(instance_id, user_id, scope) VALUES ('i1', 'u_instance', 'instance')")
	require.NoError(t, err)

	_, err = client.Exec(ctx, "INSERT INTO zitadel.administrators(instance_id, user_id, scope, organization_id) VALUES ('i1', 'u_org', 'organization', 'o1')")
	require.NoError(t, err)

	_, err = client.Exec(ctx, "INSERT INTO zitadel.administrators(instance_id, user_id, scope, project_id) VALUES ('i1', 'u_project', 'project', 'p1')")
	require.NoError(t, err)

	_, err = client.Exec(ctx, "INSERT INTO zitadel.administrators(instance_id, user_id, scope, project_grant_id) VALUES ('i1', 'u_project_grant', 'project_grant', 'pg1')")
	require.NoError(t, err)

	_, err = client.Exec(ctx, "INSERT INTO zitadel.administrator_roles(instance_id, administrator_id, role_name) VALUES ('i1', 'u_instance:instance:i1', 'instance_admin'), ('i1', 'u_org:organization:o1', 'org_admin'), ('i1', 'u_project:project:p1', 'project_admin'), ('i1', 'u_project_grant:project_grant:pg1', 'project_grant_admin')")
	require.NoError(t, err)

	// args order: instance_id, organization_id, project_id, project_grant_id, user_id, permission
	permissionTests := []struct {
		name string
		args []any
		res  bool
	}{
		{
			name: "instance scope: user has permission",
			args: []any{"i1", nil, nil, nil, "u_instance", "perm.read"},
			res:  true,
		},
		{
			name: "instance scope: user has no admin record",
			args: []any{"i1", nil, nil, nil, "u_none", "perm.read"},
			res:  false,
		},
		{
			name: "instance scope: org-scoped user without org context has no permission",
			args: []any{"i1", nil, nil, nil, "u_org", "perm.read"},
			res:  false,
		},
		{
			name: "organization scope: user has permission on matching org",
			args: []any{"i1", "o1", nil, nil, "u_org", "perm.read"},
			res:  true,
		},
		{
			name: "organization scope: user has no permission on different org",
			args: []any{"i1", "o2", nil, nil, "u_org", "perm.read"},
			res:  false,
		},
		{
			name: "organization scope: instance admin inherits permission",
			args: []any{"i1", "o1", nil, nil, "u_instance", "perm.read"},
			res:  true,
		},
		{
			name: "organization scope: role does not carry the requested permission",
			args: []any{"i1", "o1", nil, nil, "u_org", "perm.write"},
			res:  false,
		},
		{
			name: "project scope: user has permission on matching project",
			args: []any{"i1", nil, "p1", nil, "u_project", "perm.read"},
			res:  true,
		},
		{
			name: "project scope: user has no permission on different project",
			args: []any{"i1", nil, "p2", nil, "u_project", "perm.read"},
			res:  false,
		},
		{
			name: "project scope: instance admin inherits permission",
			args: []any{"i1", nil, "p1", nil, "u_instance", "perm.read"},
			res:  true,
		},
		{
			name: "project scope: org admin inherits permission when org context is provided",
			args: []any{"i1", "o1", "p1", nil, "u_org", "perm.read"},
			res:  true,
		},
		{
			name: "project grant scope: user has permission on matching project grant",
			args: []any{"i1", nil, nil, "pg1", "u_project_grant", "perm.read"},
			res:  true,
		},
		{
			name: "project grant scope: user has no permission on different project grant",
			args: []any{"i1", nil, nil, "pg2", "u_project_grant", "perm.read"},
			res:  false,
		},
		{
			name: "project grant scope: instance admin inherits permission",
			args: []any{"i1", nil, nil, "pg1", "u_instance", "perm.read"},
			res:  true,
		},
		{
			name: "project grant scope: org admin inherits permission when org context is provided",
			args: []any{"i1", "o1", nil, "pg1", "u_org", "perm.read"},
			res:  true,
		},
	}

	for _, tt := range permissionTests {
		t.Run(tt.name, func(t *testing.T) {
			builder := database.NewStatementBuilder("SELECT zitadel.check_permission(")
			builder.WriteArgs(tt.args...)
			builder.WriteString(")")

			var got bool
			require.NoError(t, client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&got), "failed to execute check query")
			assert.Equal(t, tt.res, got, "query result does not match")
		})
	}

	t.Run("permission check raises when requested", func(t *testing.T) {
		builder := database.NewStatementBuilder("SELECT zitadel.check_permission(")
		builder.WriteArgs("i1", nil, nil, nil, "u_none", "perm.read", true)
		builder.WriteString(")")

		var got bool
		err := client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&got)
		require.Error(t, err, "expected permission denial to raise an error")

		var pgErr *pgconn.PgError
		require.ErrorAs(t, err, &pgErr, "expected a pgconn.PgError")
		assert.Equal(t, "ZIT01", pgErr.Code, "unexpected Postgres error code")
		assert.Contains(t, pgErr.Message, "Permission denied", "unexpected error message")
	})
}
