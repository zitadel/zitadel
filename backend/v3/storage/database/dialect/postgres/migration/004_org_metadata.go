package migration

import (
	_ "embed"
)

var (
	//go:embed 004_org_metadata/up.sql
	up004OrgMetadata string
	//go:embed 004_org_metadata/down.sql
	down004OrgMetadata string
)

func init() {
	registerSQLMigration(4, up004OrgMetadata, down004OrgMetadata)
}
