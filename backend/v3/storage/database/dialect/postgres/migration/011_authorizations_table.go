package migration

import (
	_ "embed"
)

var (
	//go:embed 011_authorizations_table/up.sql
	up009AuthorizationsTable string
	//go:embed 011_authorizations_table/down.sql
	down009AuthorizationsTable string
)

func init() {
	registerSQLMigration(10, up009AuthorizationsTable, down009AuthorizationsTable)
}
