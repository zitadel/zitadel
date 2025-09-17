package group

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, GroupAddedEventType, GroupAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GroupChangedEventType, GroupChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GroupRemovedEventType, GroupRemovedEventMapper)
}
