package projection

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	repoDomain "github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/repository/permission"
)

type administratorRolePermission struct {
	RoleName   string `db:"role_name"`
	Permission string `db:"permission"`
}

func TestAdministratorRolePermissionReducers(t *testing.T) {
	handler := &relationalTablesProjection{}
	rawTx, tx := getTransactions(t)
	t.Cleanup(func() {
		require.NoError(t, rawTx.Rollback())
	})

	instanceID := seedAdministratorRolePermissions(t, tx)

	repo := repository.AdministratorRoleRepository()

	t.Run("added event inserts a permission row", func(t *testing.T) {
		event := permission.NewAddedEvent(t.Context(), permission.NewAggregate(instanceID), "ORG_OWNER", "org.read")
		require.True(t, callReduce(t, rawTx, handler, event))
		assert.Equal(t,
			[]administratorRolePermission{{RoleName: "ORG_OWNER", Permission: "org.read"}},
			listReducedAdministratorRolePermissions(t, tx),
		)
	})

	t.Run("removed event deletes only the matching permission row", func(t *testing.T) {
		_, err := repo.AddPermissions(t.Context(), tx, instanceID, "INSTANCE_OWNER", "instance.read", "instance.write")
		require.NoError(t, err)

		event := permission.NewRemovedEvent(t.Context(), permission.NewAggregate(instanceID), "INSTANCE_OWNER", "instance.read")
		require.True(t, callReduce(t, rawTx, handler, event))
		assert.Equal(t,
			[]administratorRolePermission{{RoleName: "INSTANCE_OWNER", Permission: "instance.write"}},
			listReducedAdministratorRolePermissions(t, tx, repo.RoleNameCondition(database.TextOperationEqual, "INSTANCE_OWNER")),
		)
	})
}

func listReducedAdministratorRolePermissions(t *testing.T, tx database.QueryExecutor, conditions ...database.Condition) []administratorRolePermission {
	t.Helper()

	builder := database.NewStatementBuilder(`SELECT role_name, permission FROM zitadel.administrator_role_permissions`)
	if len(conditions) > 0 {
		builder.WriteString(" WHERE ")
		database.And(conditions...).Write(builder)
	}
	builder.WriteString(" ORDER BY role_name, permission")

	rows, err := tx.Query(t.Context(), builder.String(), builder.Args()...)
	require.NoError(t, err)

	var result []*administratorRolePermission
	require.NoError(t, rows.(database.CollectableRows).Collect(&result))

	out := make([]administratorRolePermission, len(result))
	for i, row := range result {
		out[i] = *row
	}
	return out
}

func seedAdministratorRolePermissions(t *testing.T, tx database.QueryExecutor) string {
	t.Helper()

	instanceID := fmt.Sprintf("instance-%d", time.Now().UnixNano())
	err := repository.InstanceRepository().Create(t.Context(), tx, &repoDomain.Instance{
		ID:              instanceID,
		Name:            "instance",
		DefaultOrgID:    "default-org",
		IAMProjectID:    "iam-project",
		ConsoleClientID: "console-client",
		ConsoleAppID:    "console-app",
		DefaultLanguage: "en",
	})
	require.NoError(t, err)

	return instanceID
}
