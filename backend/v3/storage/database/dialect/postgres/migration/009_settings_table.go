package migration

import (
	_ "embed"
)

var (
	//go:embed 009_settings_table/up.sql
	up009SettingsTable string
	//go:embed 009_settings_table/down.sql
	down009SettingsTable string
)

func init() {
	registerSQLMigration(9, up009SettingsTable, down009SettingsTable)
}
