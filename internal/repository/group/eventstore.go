package group

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, GroupAddedEventType, eventstore.GenericEventMapper[GroupAddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, GroupChangedEventType, eventstore.GenericEventMapper[GroupChangedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, GroupRemovedEventType, eventstore.GenericEventMapper[GroupRemovedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, GroupUsersAddedEventType, eventstore.GenericEventMapper[GroupUsersAddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, GroupUsersChangedEventType, eventstore.GenericEventMapper[GroupUserChangedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, GroupUsersRemovedEventType, eventstore.GenericEventMapper[GroupUsersRemovedEvent])
	// TODO(legacy-cleanup): Remove these two registrations once all pre-migration group
	// memberships have been re-added via AddUsersToGroup and the old events are gone.
	// Legacy event types — registered for backward-compatible replay only; never written.
	eventstore.RegisterFilterEventMapper(AggregateType, LegacyGroupUserAddedEventType, eventstore.GenericEventMapper[LegacyGroupUserAddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, LegacyGroupUserRemovedEventType, eventstore.GenericEventMapper[LegacyGroupUserRemovedEvent])
	// TODO(legacy-cleanup): Remove these two registrations once all pre-migration group
	// memberships have been re-added via AddUsersToGroup and the old member events are gone.
	eventstore.RegisterFilterEventMapper(AggregateType, LegacyGroupMemberAddedEventType, eventstore.GenericEventMapper[LegacyGroupMemberAddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, LegacyGroupMemberRemovedEventType, eventstore.GenericEventMapper[LegacyGroupMemberRemovedEvent])
	// eventstore.RegisterFilterEventMapper(AggregateType, GroupUserCascadeRemovedEventType, eventstore.GenericEventMapper[GroupUserCascadeRemovedEvent]) // TODO: Check if need a cascade removed event.
	eventstore.RegisterFilterEventMapper(AggregateType, GroupDeactivatedEventType, eventstore.GenericEventMapper[GroupDeactivatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, GroupReactivatedEventType, eventstore.GenericEventMapper[GroupReactivatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, MetadataSetType, MetadataSetEventMapper)               // TODO: Replace with generic mapper after refactor
	eventstore.RegisterFilterEventMapper(AggregateType, MetadataRemovedType, MetadataRemovedEventMapper)       // TODO: Replace with generic mapper after refactor
	eventstore.RegisterFilterEventMapper(AggregateType, MetadataRemovedAllType, MetadataRemovedAllEventMapper) // TODO: Replace with generic mapper after refactor
}
