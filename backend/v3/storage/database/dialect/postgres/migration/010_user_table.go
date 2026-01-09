package migration

import (
	_ "embed"
)

var (
	//go:embed 010_user_table/up.sql
	up010UserTable string
	//go:embed 010_user_table/down.sql
	down010UserTable string
)

func init() {
	registerSQLMigration(10, up010UserTable, down010UserTable)
}
