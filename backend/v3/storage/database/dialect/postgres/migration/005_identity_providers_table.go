package migration

import (
	_ "embed"
)

var (
	//go:embed 005_identity_providers_table/up.sql
	up005IdentityProvidersTable string
	//go:embed 005_identity_providers_table/down.sql
	down005IdentityProvidersTable string
)

func init() {
	RegisterSQLMigration(5, up005IdentityProvidersTable, down005IdentityProvidersTable)
}
