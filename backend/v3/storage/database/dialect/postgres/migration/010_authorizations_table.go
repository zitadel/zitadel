package migration

import (
	_ "embed"
)

var (
	//go:embed 010_authorizations_table/up.sql
	up010AuthorizationsTable string
	//go:embed 010_authorizations_table/down.sql
	down010AuthorizationsTable string
)

func init() {
	RegisterSQLMigration(10, up010AuthorizationsTable, down010AuthorizationsTable)
}
