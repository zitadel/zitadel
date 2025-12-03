package migration

import (
	_ "embed"
)

var (
	//go:embed 011_session_table/up.sql
	up010SessionTable string
	//go:embed 011_session_table/down.sql
	down010SessionTable string
)

func init() {
	registerSQLMigration(10, up010SessionTable, down010SessionTable) // TODO: needs to be set to 11 after user (10) is merged
}
