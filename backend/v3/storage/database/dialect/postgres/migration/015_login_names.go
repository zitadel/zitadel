package migration

import (
	_ "embed"
)

var (
	//go:embed 015_login_names/up.sql
	up015LoginNames string
	//go:embed 015_login_names/down.sql
	down015LoginNames string
)

func init() {
	registerSQLMigration(15, up015LoginNames, down015LoginNames)
}
