package migration

import (
	_ "embed"
)

var (
	//go:embed 010_session_table/up.sql
	up010SessionTable string
	//go:embed 010_session_table/down.sql
	down010SessionTable string
)

func init() {
	registerSQLMigration(9, up010SessionTable, down010SessionTable) // TODO: needs to be set to 10 after user (9) is merged
}
