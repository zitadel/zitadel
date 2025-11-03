package migration

import (
	_ "embed"
)

var (
	//go:embed 006_projects_table/up.sql
	up006ProjectsTable string
	//go:embed 006_projects_table/down.sql
	down006ProjectsTable string
)

func init() {
	registerSQLMigration(6, up006ProjectsTable, down006ProjectsTable)
}
