package migration

import (
	_ "embed"
)

var (
	//go:embed 009_verification_table/up.sql
	up009VerificationTable string
	//go:embed 009_verification_table/down.sql
	down009VerificationTable string
)

func init() {
	registerSQLMigration(9, up009VerificationTable, down009VerificationTable)
}
