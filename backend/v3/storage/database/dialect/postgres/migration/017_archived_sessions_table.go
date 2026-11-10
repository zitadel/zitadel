package migration

import (
	_ "embed"
)

var (
	//go:embed 017_archived_sessions_table/up.sql
	up017ArchivedSessionsTable string
	//go:embed 017_archived_sessions_table/down.sql
	down017ArchivedSessionsTable string
)

func init() {
	registerSQLMigration(17, up017ArchivedSessionsTable, down017ArchivedSessionsTable)
}
