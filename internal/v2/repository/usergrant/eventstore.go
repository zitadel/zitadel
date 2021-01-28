package usergrant

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(UserGrantAddedType, UserGrantAddedEventMapper).
		RegisterFilterEventMapper(UserGrantChangedType, UserGrantChangedEventMapper).
		RegisterFilterEventMapper(UserGrantCascadeChangedType, UserGrantCascadeChangedEventMapper).
		RegisterFilterEventMapper(UserGrantRemovedType, UserGrantRemovedEventMapper).
		RegisterFilterEventMapper(UserGrantCascadeRemovedType, UserGrantCascadeRemovedEventMapper).
		RegisterFilterEventMapper(UserGrantDeactivatedType, UserGrantDeactivatedEventMapper).
		RegisterFilterEventMapper(UserGrantReactivatedType, UserGrantReactivatedEventMapper)
}
