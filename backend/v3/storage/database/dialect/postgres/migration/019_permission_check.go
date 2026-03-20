package migration

import (
	_ "embed"
)

var (
	//go:embed 019_permission_check/up.sql
	up019PermissionCheck string
	//go:embed 019_permission_check/down.sql
	down019PermissionCheck string
)

func init() {
	registerSQLMigration(19, up019PermissionCheck, down019PermissionCheck)
}
