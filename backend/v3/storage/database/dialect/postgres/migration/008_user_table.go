package migration

import (
	_ "embed"
)

var (
	//go:embed 008_user_table/up.sql
	up008UserTable string
	//go:embed 008_user_table/down.sql
	down008UserTable string
)

func init() {
	registerSQLMigration(8, up008UserTable, down008UserTable)
}
