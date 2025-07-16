package organization_settings

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, OrganizationSettingsSetEventType, eventstore.GenericEventMapper[OrganizationSettingsSetEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, OrganizationSettingsRemovedEventType, eventstore.GenericEventMapper[OrganizationSettingsRemovedEvent])
}
