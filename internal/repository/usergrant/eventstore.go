package usergrant

import (
	"github.com/zitadel/zitadel/internal/eventstore"
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
