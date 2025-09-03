package migration

import (
	_ "embed"
)

var (
	//go:embed 004_identity_providers_table/up.sql
	up004IdentityProvidersTable string
	//go:embed 004_identity_providers_table/down.sql
	down004IdentityProvidersTable string
)

func init() {
	registerSQLMigration(4, up004IdentityProvidersTable, down004IdentityProvidersTable)
}
