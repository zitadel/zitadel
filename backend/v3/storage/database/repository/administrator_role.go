package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.AdministratorRoleRepository = (*administratorRole)(nil)

type administratorRole struct{}

func AdministratorRoleRepository() domain.AdministratorRoleRepository {
	return new(administratorRole)
}

func (administratorRole) unqualifiedTableName() string {
	return "administrator_role_permissions"
}

func (a administratorRole) qualifiedTableName() string {
	return "zitadel." + a.unqualifiedTableName()
}

// AddPermission implements [domain.AdministratorRoleRepository].
func (a administratorRole) AddPermissions(ctx context.Context, client database.QueryExecutor, instanceID, role string, permissions ...string) (int64, error) {
	if len(permissions) == 0 {
		return 0, database.ErrNoChanges
	}
	builder := database.NewStatementBuilder(
		"INSERT INTO zitadel.administrator_role_permissions (instance_id, role_name, permission)"+
			" SELECT $1::text, $2::text, unnest($3::text[])"+
			" ON CONFLICT (instance_id, permission, role_name) DO NOTHING", instanceID, role, permissions,
	)
	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// RemovePermissions implements [domain.AdministratorRoleRepository].
func (a administratorRole) RemovePermissions(ctx context.Context, client database.QueryExecutor, instanceID, role string, permissionsToRemove ...string) (int64, error) {
	if len(permissionsToRemove) == 0 {
		return 0, database.ErrNoChanges
	}
	builder := database.NewStatementBuilder(
		"DELETE FROM zitadel.administrator_role_permissions"+
			" WHERE instance_id = $1 AND role_name = $2 AND permission = ANY($3::text[])", instanceID, role, permissionsToRemove,
	)
	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// PrimaryKeyCondition implements [domain.AdministratorRoleRepository].
func (a administratorRole) PrimaryKeyCondition(instanceID, roleName, permission string) database.Condition {
	return database.And(
		a.InstanceIDCondition(instanceID),
		a.RoleNameCondition(database.TextOperationEqual, roleName),
		a.PermissionCondition(database.TextOperationEqual, permission),
	)
}

// InstanceIDCondition implements [domain.AdministratorRoleRepository].
func (a administratorRole) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(a.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// RoleNameCondition implements [domain.AdministratorRoleRepository].
func (a administratorRole) RoleNameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(a.RoleNameColumn(), op, name)
}

// PermissionCondition implements [domain.AdministratorRoleRepository].
func (a administratorRole) PermissionCondition(op database.TextOperation, permission string) database.Condition {
	return database.NewTextCondition(a.PermissionColumn(), op, permission)
}

// PrimaryKeyColumns implements [domain.AdministratorRoleRepository].
func (a administratorRole) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		a.InstanceIDColumn(),
		a.PermissionColumn(),
		a.RoleNameColumn(),
	}
}

// RoleNameColumn implements [domain.AdministratorRoleRepository].
func (a administratorRole) RoleNameColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "role_name")
}

// PermissionColumn implements [domain.AdministratorRoleRepository].
func (a administratorRole) PermissionColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "permission")
}

// InstanceIDColumn implements [domain.AdministratorRoleRepository].
func (a administratorRole) InstanceIDColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "instance_id")
}
