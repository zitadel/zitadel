package migration

import (
	_ "embed"
)

var (
	//go:embed 019_password_locked_at/up.sql
	up019PasswordLockedAt string
	//go:embed 019_password_locked_at/down.sql
	down019PasswordLockedAt string
)

func init() {
	registerSQLMigration(19, up019PasswordLockedAt, down019PasswordLockedAt)
}
