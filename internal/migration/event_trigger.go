package migration

import (
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
)

const (
	eventTriggerTmpl = "event_trigger"
)

// EventTrigger registers triggers that record INSERT/UPDATE/DELETE
// occurrences into projections.service_ping_resource_events

func EventTrigger(
	db *database.DB,
	table string,
	parentType domain.CountParentType,
	instanceIDColumn string,
	parentIDColumn string,
	resource string,
) RepeatableMigration {
	return EventTriggerConditional(
		db,
		table,
		parentType,
		instanceIDColumn,
		parentIDColumn,
		resource,
		nil,
	)
}

// EventTriggerConditional registers the event triggers with optional conditions.
func EventTriggerConditional(
	db *database.DB,
	table string,
	parentType domain.CountParentType,
	instanceIDColumn string,
	parentIDColumn string,
	resource string,
	conditions TriggerConditions,
) RepeatableMigration {
	return &triggerMigration{
		triggerConfig: triggerConfig{
			Table:            table,
			ParentType:       parentType.String(),
			InstanceIDColumn: instanceIDColumn,
			ParentIDColumn:   parentIDColumn,
			Resource:         resource,
			Conditions:       conditions,
		},
		db:           db,
		templateName: eventTriggerTmpl,
	}
}
