package migration

import (
	_ "embed"
)

var (
	//go:embed 007_project_grants_table/up.sql
	up007ProjectGrantsTable string
	//go:embed 007_project_grants_table/down.sql
	down007ProjectGrantsTable string
)

func init() {
	registerSQLMigration(7, up007ProjectGrantsTable, down007ProjectGrantsTable)
}
