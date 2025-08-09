package migration

import (
	_ "embed"
)

var (
	//go:embed 002_organization_table/up.sql
	up002OrganizationTable string
	//go:embed 002_organization_table/down.sql
	down002OrganizationTable string
)

func init() {
	registerSQLMigration(2, up002OrganizationTable, down002OrganizationTable)
}
