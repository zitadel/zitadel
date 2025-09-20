package migration

import (
	_ "embed"
)

var (
	//go:embed 004_org_metadata_table/up.sql
	up004OrgMetadataTable string
	//go:embed 004_org_metadata_table/down.sql
	down004OrgMetadataTable string
)

func init() {
	registerSQLMigration(4, up004OrgMetadataTable, down004OrgMetadataTable)
}
