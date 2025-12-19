package migration

import (
	_ "embed"
)

var (
	//go:embed 011_session_table/up.sql
	up011SessionTable string
	//go:embed 011_session_table/down.sql
	down011SessionTable string
)

func init() {
	registerSQLMigration(11, up011SessionTable, down011SessionTable)
}
