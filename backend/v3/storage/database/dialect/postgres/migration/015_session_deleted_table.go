package migration

import (
	_ "embed"
)

var (
	//go:embed 015_session_deleted_table/up.sql
	up015SessionDeletedTable string
	//go:embed 015_session_deleted_table/down.sql
	down015SessionDeletedTable string
)

func init() {
	registerSQLMigration(15, up015SessionDeletedTable, down015SessionDeletedTable)
}
