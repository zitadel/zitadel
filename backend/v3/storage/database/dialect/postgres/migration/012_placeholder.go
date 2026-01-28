package migration

func init() {
	registerSQLMigration(12, "select now()", "select now()")
}
