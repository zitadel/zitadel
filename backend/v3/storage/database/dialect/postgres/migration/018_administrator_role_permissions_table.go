package migration

import (
	_ "embed"
)

var (
	//go:embed 018_administrator_role_permissions_table/up.sql
	up018AdministratorRolePermissionsTable string
	//go:embed 018_administrator_role_permissions_table/down.sql
	down018AdministratorRolePermissionsTable string
)

func init() {
	registerSQLMigration(18, up018AdministratorRolePermissionsTable, down018AdministratorRolePermissionsTable)
}
