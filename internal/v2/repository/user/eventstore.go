package user

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(HumanAddedEventType, HumanAddedMapper).
		RegisterFilterEventMapper(HumanRegisteredEventType, HumanRegisteredMapper).
		RegisterFilterEventMapper(HumanInitialCodeAddedType, HumanInitialCodeAddedMapper).
		RegisterFilterEventMapper(HumanInitialCodeSentType, HumanInitialCodeSentMapper).
		RegisterFilterEventMapper(HumanInitializedCheckSucceededType, HumanInitializedCheckSucceededMapper).
		RegisterFilterEventMapper(HumanInitializedCheckFailedType, HumanInitializedCheckFailedEventMapper)
}
