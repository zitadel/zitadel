package groupgrant

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, GroupGrantAddedType, GroupGrantAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GroupGrantChangedType, GroupGrantChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GroupGrantRemovedType, GroupGrantRemovedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GroupGrantCascadeRemovedType, GroupGrantCascadeRemovedEventMapper)
}
