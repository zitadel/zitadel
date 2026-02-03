package setup

import (
	_ "embed"

	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/migration"
)

const (
	downFieldsTable        = `DROP TABLE IF EXISTS eventstore.fields;`
	downFailedEventsTable  = `DROP TABLE IF EXISTS projections.failed_events2;`
	downCurrentStatesTable = `DROP TABLE IF EXISTS projections.current_states;`
)

var (
	//go:embed 15/01_new_failed_events.sql
	upFailedEventsTable string
	//go:embed 15/05_current_states.sql
	upCurrentStatesTable string
)

func init() {
	migration.RegisterSQLMigrationNoSequence(addFieldTable, downFieldsTable)
	migration.RegisterSQLMigrationNoSequence(upFailedEventsTable, downFailedEventsTable)
	migration.RegisterSQLMigrationNoSequence(upCurrentStatesTable+addOffsetField, downCurrentStatesTable)
}
