package migration

import (
	_ "embed"
)

var (
	//go:embed 007_user_table/up.sql
	up004UserTable string
	//go:embed 007_user_table/down.sql
	down007UserTable string
)

func init() {
	registerSQLMigration(5, up004UserTable, down007UserTable)
}
