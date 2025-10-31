package migration

import (
	_ "embed"
)

var (
	//go:embed 008_project_grants_table/up.sql
	up008ProjectGrantsTable string
	//go:embed 008_project_grants_table/down.sql
	down008ProjectGrantsTable string
)

func init() {
	registerSQLMigration(8, up008ProjectGrantsTable, down008ProjectGrantsTable)
}
