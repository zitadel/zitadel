package action

import "github.com/zitadel/zitadel/internal/eventstore"

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AddedEventType, AddedEventMapper).
		RegisterFilterEventMapper(ChangedEventType, ChangedEventMapper).
		RegisterFilterEventMapper(DeactivatedEventType, DeactivatedEventMapper).
		RegisterFilterEventMapper(ReactivatedEventType, ReactivatedEventMapper).
		RegisterFilterEventMapper(RemovedEventType, RemovedEventMapper)
}
