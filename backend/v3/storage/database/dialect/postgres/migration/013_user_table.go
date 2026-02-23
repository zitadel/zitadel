package migration

import (
	_ "embed"
)

var (
	//go:embed 013_user_table/up.sql
	up013UserTable string
	//go:embed 013_user_table/down.sql
	down013UserTable string
)

func init() {
	registerSQLMigration(13, up013UserTable, down013UserTable)
}
