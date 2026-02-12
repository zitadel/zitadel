package migration

import (
	_ "embed"
)

var (
	//go:embed 012_identity_provider_intents_table/up.sql
	up012IDPIntentsTable string
	//go:embed 012_identity_provider_intents_table/down.sql
	down012IDPIntentsTable string
)

func init() {
	RegisterSQLMigration(12, up012IDPIntentsTable, down012IDPIntentsTable)
}
