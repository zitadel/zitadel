package migrations

import (
	_ "embed"

	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/migration"
)

// Projections Schema
var (
	//go:embed projections_schema/up.sql
	projectionsSchemaUp string
	//go:embed projections_schema/down.sql
	projectionsSchemaDown string

	//go:embed current_states/up.sql
	currentStatesUp string
	//go:embed current_states/down.sql
	currentStatesDown string

	//go:embed failed_events/up.sql
	failedEventsUp string
	//go:embed failed_events/down.sql
	failedEventsDown string
)

// Eventstore Schema
var (
	//go:embed eventstore_schema/up.sql
	eventstoreUp string
	//go:embed eventstore_schema/down.sql
	eventstoreDown string

	//go:embed events/up.sql
	eventsUp string
	//go:embed events/down.sql
	eventsDown string

	//go:embed fields/up.sql
	fieldsUp string
	//go:embed fields/down.sql
	fieldsDown string

	//go:embed unique_constraints/up.sql
	uniqueConstraintsUp string
	//go:embed unique_constraints/down.sql
	uniqueConstraintsDown string
)

func init() {
	projectionsMigUp := projectionsSchemaUp + currentStatesUp + failedEventsUp
	projectionsMigDown := failedEventsDown + currentStatesDown + projectionsSchemaDown
	migration.RegisterSQLMigrationNoSequence(projectionsMigUp, projectionsMigDown)

	evenstoreMigUp := eventstoreUp + eventsUp + fieldsUp + uniqueConstraintsUp
	evenstoreMigDown := uniqueConstraintsDown + fieldsDown + eventsDown + eventstoreDown
	migration.RegisterSQLMigrationNoSequence(evenstoreMigUp, evenstoreMigDown)
}
