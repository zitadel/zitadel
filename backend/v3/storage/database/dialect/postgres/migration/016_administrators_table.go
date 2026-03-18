package migration

import (
	_ "embed"
)

var (
	//go:embed 016_administrators_table/up.sql
	up016AdministratorsTable string
	//go:embed 016_administrators_table/down.sql
	down016AdministratorsTable string
)

func init() {
	registerSQLMigration(16, up016AdministratorsTable, down016AdministratorsTable)
}
