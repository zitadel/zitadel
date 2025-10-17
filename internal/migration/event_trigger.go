package migration

import (
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
)

const (
	eventTriggerTmpl = "event_trigger"
)

// EventTrigger registers triggers that record INSERT/UPDATE/DELETE
// occurrences into projections.service_ping_resource_events using
// projections.record_service_ping_resource_event.
//
// This migration repeats when any of the arguments are changed.
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

// // String implements [Migration] and [fmt.Stringer].
// func (m *triggerMigration) String() string {
// 	// keep existing format but distinguish by template name
// 	return fmt.Sprintf("repeatable_%s_%s", m.Resource, m.templateName)
// }

// Execute implements [Migration] (delegates to shared templates executor)
// func (m *triggerMigration) Execute(ctx context.Context, _ eventstore.Event) error {
// 	var query strings.Builder
// 	err := templates.ExecuteTemplate(&query, m.templateName, m.triggerConfig)
// 	if err != nil {
// 		return fmt.Errorf("%s: execute trigger template: %w", m, err)
// 	}
// 	_, err = m.db.ExecContext(ctx, query.String())
// 	if err != nil {
// 		return fmt.Errorf("%s: exec trigger query: %w", m, err)
// 	}
// 	return nil
// }
