package migration

import (
	_ "embed"
)

var (
	//go:embed 001_instance_table/up.sql
	up001InstanceTable string
	//go:embed 001_instance_table/down.sql
	down001InstanceTable string
)

func init() {
	registerSQLMigration(1, up001InstanceTable, down001InstanceTable)
}
