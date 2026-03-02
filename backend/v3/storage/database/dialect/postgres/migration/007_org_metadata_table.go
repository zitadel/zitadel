package migration

import (
	_ "embed"
)

var (
	//go:embed 007_organization_metadata_table/up.sql
	up007OrganizationMetadataTable string
	//go:embed 007_organization_metadata_table/down.sql
	down007OrganizationMetadataTable string
)

func init() {
	registerSQLMigration(7, up007OrganizationMetadataTable, down007OrganizationMetadataTable)
}
