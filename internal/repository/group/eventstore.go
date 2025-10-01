package group

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, GroupAddedEventType, eventstore.GenericEventMapper[GroupAddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, GroupChangedEventType, eventstore.GenericEventMapper[GroupChangedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, GroupRemovedEventType, eventstore.GenericEventMapper[GroupRemovedEvent])
}
