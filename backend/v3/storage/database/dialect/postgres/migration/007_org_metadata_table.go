package migration

import (
	_ "embed"
)

var (
	//go:embed 007_org_metadata_table/up.sql
	up007OrgMetadataTable string
	//go:embed 007_org_metadata_table/down.sql
	down007OrgMetadataTable string
)

func init() {
	RegisterSQLMigration(7, up007OrgMetadataTable, down007OrgMetadataTable)
}
