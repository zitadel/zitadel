package migration

import (
	_ "embed"
)

var (
	//go:embed 003_domains_table/up.sql
	up003DomainsTable string
	//go:embed 003_domains_table/down.sql
	down003DomainsTable string
)

func init() {
	registerSQLMigration(3, up003DomainsTable, down003DomainsTable)
}
