package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type administratorRoleColumns interface {
	PrimaryKeyColumns() []database.Column
	InstanceIDColumn() database.Column
	RoleNameColumn() database.Column
	PermissionColumn() database.Column
}

type administratorRoleConditions interface {
	PrimaryKeyCondition(instanceID, roleName, permission string) database.Condition
	InstanceIDCondition(instanceID string) database.Condition
	RoleNameCondition(op database.TextOperation, name string) database.Condition
	PermissionCondition(op database.TextOperation, permission string) database.Condition
}

//go:generate mockgen -typed -package domainmock -destination ./mock/administrator_role.mock.go . AdministratorRoleRepository
type AdministratorRoleRepository interface {
	Repository

	administratorRoleColumns
	administratorRoleConditions

	AddPermissions(ctx context.Context, client database.QueryExecutor, instanceID, role string, permissions ...string) (int64, error)
	RemovePermissions(ctx context.Context, client database.QueryExecutor, instanceID, role string, permissionsToRemove ...string) (int64, error)
}
