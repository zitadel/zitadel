package migration

import (
	_ "embed"
)

var (
	//go:embed 012_verification_table/up.sql
	up012VerificationTable string
	//go:embed 012_verification_table/down.sql
	down012VerificationTable string
)

func init() {
	registerSQLMigration(12, up012VerificationTable, down012VerificationTable)
}
