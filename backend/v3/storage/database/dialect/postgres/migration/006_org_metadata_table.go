package migration

import (
	_ "embed"
)

var (
	//go:embed 006_org_metadata_table/up.sql
	up006OrgMetadataTable string
	//go:embed 006_org_metadata_table/down.sql
	down006OrgMetadataTable string
)

func init() {
	registerSQLMigration(6, up006OrgMetadataTable, down006OrgMetadataTable)
}
