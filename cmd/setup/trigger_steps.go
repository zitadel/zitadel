package setup

import (
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/migration"
	"github.com/zitadel/zitadel/internal/query/projection"
)

// triggerSteps defines the repeatable migrations that set up triggers
// for counting resources in the database.
func triggerSteps(db *database.DB) []migration.RepeatableMigration {
	return []migration.RepeatableMigration{
		migration.DeleteParentCountsTrigger(db,
			projection.InstanceProjectionTable,
			domain.CountParentTypeInstance,
			projection.InstanceColumnID,
			projection.InstanceColumnID,
			"instance",
		),
		migration.DeleteParentCountsTrigger(db,
			projection.OrgProjectionTable,
			domain.CountParentTypeOrganization,
			projection.OrgColumnInstanceID,
			projection.OrgColumnID,
			"organization",
		),
		migration.CountTrigger(db,
			projection.UserTable,
			domain.CountParentTypeOrganization,
			projection.UserInstanceIDCol,
			projection.UserResourceOwnerCol,
			"user",
		),
	}
}
