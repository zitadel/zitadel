package migration

import (
	_ "embed"
)

var (
	//go:embed 016_archived_sessions_table/up.sql
	up016ArchivedSessionsTable string
	//go:embed 016_archived_sessions_table/down.sql
	down016ArchivedSessionsTable string
)

func init() {
	registerSQLMigration(16, up016ArchivedSessionsTable, down016ArchivedSessionsTable)
}
