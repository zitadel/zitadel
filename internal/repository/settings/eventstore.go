package settings

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, SettingOrganizationSetEventType, eventstore.GenericEventMapper[SettingOrganizationSetEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, SettingOrganizationRemovedEventType, eventstore.GenericEventMapper[SettingOrganizationRemovedEvent])
}
