package migration

import (
	_ "embed"
)

var (
	//go:embed 006_settings_table/up.sql
	up004SettingsTable string
	//go:embed 006_settings_table/down.sql
	down004SettingsTable string
)

func init() {
	registerSQLMigration(6, up004SettingsTable, down004SettingsTable)
}
