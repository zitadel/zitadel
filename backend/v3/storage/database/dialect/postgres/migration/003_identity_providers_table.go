package migration

import (
	_ "embed"
)

var (
	//go:embed 003_identity_providers_table/up.sql
	up003IdentityProvidersTable string
	//go:embed 003_identity_providers_table/down.sql
	down003IdentityProvidersTable string
)

func init() {
	registerSQLMigration(3, up003IdentityProvidersTable, down003IdentityProvidersTable)
}
