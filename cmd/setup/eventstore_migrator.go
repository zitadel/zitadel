package setup

import (
	_ "embed"

	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/migration"
)

const (
	downFieldsTable = `DROP TABLE IF EXISTS eventstore.fields;`
)

func init() {
	migration.RegisterSQLMigrationNoSequence(addFieldTable, downFieldsTable)
}
