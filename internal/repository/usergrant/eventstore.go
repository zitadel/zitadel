package usergrant

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, UserGrantAddedType, UserGrantAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserGrantChangedType, UserGrantChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserGrantCascadeChangedType, UserGrantCascadeChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserGrantRemovedType, UserGrantRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserGrantCascadeRemovedType, UserGrantCascadeRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserGrantDeactivatedType, UserGrantDeactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserGrantReactivatedType, UserGrantReactivatedEventMapper)
}
