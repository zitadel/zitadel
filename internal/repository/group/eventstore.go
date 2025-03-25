package group

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, GroupAddedType, GroupAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GroupChangedType, GroupChangeEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GroupDeactivatedType, GroupDeactivatedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GroupReactivatedType, GroupReactivatedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GroupRemovedType, GroupRemovedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, MemberAddedType, MemberAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, MemberChangedType, MemberChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, MemberRemovedType, MemberRemovedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, MemberCascadeRemovedType, MemberCascadeRemovedEventMapper)
}
