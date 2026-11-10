package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

type administratorRolePermissionRow struct {
	InstanceID string `db:"instance_id"`
	RoleName   string `db:"role_name"`
	Permission string `db:"permission"`
}

func TestAdministratorRoleRepository_AddPermissions(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	repo := repository.AdministratorRoleRepository()

	tests := []struct {
		name        string
		role        string
		permissions []string
		wantRows    []administratorRolePermissionRow
		wantErr     error
		wantCount   int64
	}{
		{
			name:        "add deduplicated permissions",
			role:        "ORG_OWNER",
			permissions: []string{"org.read", "org.write", "org.read"},
			wantCount:   2,
			wantRows: []administratorRolePermissionRow{
				{InstanceID: instanceID, RoleName: "ORG_OWNER", Permission: "org.read"},
				{InstanceID: instanceID, RoleName: "ORG_OWNER", Permission: "org.write"},
			},
		},
		{
			name:      "no permissions",
			role:      "ORG_OWNER",
			wantErr:   database.ErrNoChanges,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollbackSavepoint := savepointForRollback(t, tx)
			defer rollbackSavepoint()

			count, err := repo.AddPermissions(t.Context(), savepoint, instanceID, tt.role, tt.permissions...)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantCount, count)

			if tt.wantErr == nil {
				assert.Equal(t,
					tt.wantRows,
					listAdministratorRolePermissions(t, savepoint, repo.RoleNameCondition(database.TextOperationEqual, tt.role)),
				)
			}
		})
	}
}

func TestAdministratorRoleRepository_RemovePermissions(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)

	repo := repository.AdministratorRoleRepository()
	count, err := repo.AddPermissions(t.Context(), tx, instanceID, "INSTANCE_OWNER", "instance.read", "instance.write")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	tests := []struct {
		name        string
		role        string
		permissions []string
		wantRows    []administratorRolePermissionRow
		wantErr     error
		wantCount   int64
	}{
		{
			name:    "no permissions",
			role:    "INSTANCE_OWNER",
			wantErr: database.ErrNoChanges,
		},
		{
			name:        "remove permissions",
			role:        "INSTANCE_OWNER",
			permissions: []string{"instance.read"},
			wantCount:   1,
			wantRows: []administratorRolePermissionRow{
				{InstanceID: instanceID, RoleName: "INSTANCE_OWNER", Permission: "instance.write"},
			},
		},
		{
			name:        "remove exact pair and missing values",
			role:        "INSTANCE_OWNER",
			permissions: []string{"instance.write", "missing"},
			wantCount:   1,
			wantRows:    []administratorRolePermissionRow{{InstanceID: instanceID, RoleName: "INSTANCE_OWNER", Permission: "instance.read"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollbackSavepoint := savepointForRollback(t, tx)
			defer rollbackSavepoint()

			count, err = repo.RemovePermissions(t.Context(), savepoint, instanceID, tt.role, tt.permissions...)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantCount, count)
			if tt.wantErr == nil {
				assert.Equal(t,
					tt.wantRows,
					listAdministratorRolePermissions(t, savepoint, repo.RoleNameCondition(database.TextOperationEqual, "INSTANCE_OWNER")),
				)
			}
		})
	}
}

func TestAdministratorRoleRepository_AddAndRemoveAcrossCalls(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)

	repo := repository.AdministratorRoleRepository()

	count, err := repo.AddPermissions(t.Context(), tx, instanceID, "INSTANCE_OWNER", "instance.read", "instance.write")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	count, err = repo.AddPermissions(t.Context(), tx, instanceID, "INSTANCE_OWNER", "instance.manage", "instance.audit", "instance.write")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	count, err = repo.RemovePermissions(t.Context(), tx, instanceID, "INSTANCE_OWNER", "instance.read")
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	assert.Equal(t,
		[]administratorRolePermissionRow{
			{InstanceID: instanceID, RoleName: "INSTANCE_OWNER", Permission: "instance.audit"},
			{InstanceID: instanceID, RoleName: "INSTANCE_OWNER", Permission: "instance.manage"},
			{InstanceID: instanceID, RoleName: "INSTANCE_OWNER", Permission: "instance.write"},
		},
		listAdministratorRolePermissions(t, tx, repo.RoleNameCondition(database.TextOperationEqual, "INSTANCE_OWNER")),
	)
}

func TestAdministratorRoleRepository_PrimaryKeyCondition(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)

	repo := repository.AdministratorRoleRepository()
	count, err := repo.AddPermissions(t.Context(), tx, instanceID, "PROJECT_OWNER", "project.read", "project.write")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	rows := listAdministratorRolePermissions(t, tx, repo.PrimaryKeyCondition(instanceID, "PROJECT_OWNER", "project.read"))
	assert.Equal(t,
		[]administratorRolePermissionRow{{InstanceID: instanceID, RoleName: "PROJECT_OWNER", Permission: "project.read"}},
		rows,
	)
}

func listAdministratorRolePermissions(t *testing.T, tx database.QueryExecutor, condition database.Condition) []administratorRolePermissionRow {
	t.Helper()

	builder := database.NewStatementBuilder(`SELECT instance_id, role_name, permission FROM zitadel.administrator_role_permissions`)
	if condition != nil {
		builder.WriteString(" WHERE ")
		condition.Write(builder)
	}
	builder.WriteString(" ORDER BY instance_id, role_name, permission")

	rows, err := tx.Query(t.Context(), builder.String(), builder.Args()...)
	require.NoError(t, err)

	var result []*administratorRolePermissionRow
	require.NoError(t, rows.(database.CollectableRows).Collect(&result))

	out := make([]administratorRolePermissionRow, len(result))
	for i, row := range result {
		out[i] = *row
	}
	return out
}
