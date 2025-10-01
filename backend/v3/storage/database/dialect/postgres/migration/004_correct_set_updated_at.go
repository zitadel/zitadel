package migration

import (
	_ "embed"
)

var (
	//go:embed 004_correct_set_updated_at/up.sql
	up004CorrectSetUpdatedAt string
	//go:embed 004_correct_set_updated_at/down.sql
	down004CorrectSetUpdatedAt string
)

func init() {
	registerSQLMigration(4, up004CorrectSetUpdatedAt, down004CorrectSetUpdatedAt)
}
