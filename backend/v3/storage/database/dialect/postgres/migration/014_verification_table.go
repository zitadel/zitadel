package migration

import (
	_ "embed"
)

var (
	//go:embed 014_verification_table/up.sql
	up014VerificationTable string
	//go:embed 014_verification_table/down.sql
	down014VerificationTable string
)

func init() {
	registerSQLMigration(14, up014VerificationTable, down014VerificationTable)
}
