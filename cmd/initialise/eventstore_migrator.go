package initialise

import (
	_ "embed"

	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/migration"
)

const (
	downEventsTable = `
DROP TABLE IF EXISTS eventstore.events2;
DROP FUNCTION IF EXISTS eventstore.commands_to_events(commands eventstore.command[]), eventstore.push(commands eventstore.command[]);
DROP TYPE IF EXISTS eventstore.command;
`
	downUniqueConstraintsTable = "DROP TABLE IF EXISTS eventstore.unique_constraints;"
	upSchemaCreation           = `CREATE SCHEMA IF NOT EXISTS eventstore;`
	downSchemaCreation         = `DROP SCHEMA IF EXISTS eventstore;`
)

var (
	//go:embed sql/08_events_table.sql
	upEventsTable string
	//go:embed sql/10_unique_constraints_table.sql
	upUniqueConstraintsTable string
)

func init() {
	migration.RegisterSQLMigrationNoSequence(upSchemaCreation, downSchemaCreation)
	migration.RegisterSQLMigrationNoSequence(upEventsTable, downEventsTable)
	migration.RegisterSQLMigrationNoSequence(upUniqueConstraintsTable, downUniqueConstraintsTable)
}
