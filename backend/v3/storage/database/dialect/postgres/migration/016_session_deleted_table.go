package migration

import (
	_ "embed"
)

var (
	//go:embed 016_session_deleted_table/up.sql
	up016SessionDeletedTable string
	//go:embed 016_session_deleted_table/down.sql
	down016SessionDeletedTable string
)

func init() {
	registerSQLMigration(16, up016SessionDeletedTable, down016SessionDeletedTable)
}
