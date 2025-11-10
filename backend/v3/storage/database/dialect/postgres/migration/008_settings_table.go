package migration

import (
	_ "embed"
)

var (
	//go:embed 008_settings_table/up.sql
	up008SettingsTable string
	//go:embed 008_settings_table/down.sql
	down008SettingsTable string
)

func init() {
	registerSQLMigration(8, up008SettingsTable, down008SettingsTable)
}
