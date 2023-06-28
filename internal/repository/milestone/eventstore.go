package milestone

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, eventstore.EventType(PushedProjectCreatedEventType), eventstore.GenericEventMapper[InstanceCreatedPushedEvent]).
		RegisterFilterEventMapper(AggregateType, eventstore.EventType(PushedAuthenticationSucceededOnInstanceEventType), eventstore.GenericEventMapper[AuthenticationSucceededOnInstancePushedEvent]).
		RegisterFilterEventMapper(AggregateType, eventstore.EventType(PushedProjectCreatedEventType), eventstore.GenericEventMapper[ProjectCreatedPushedEvent]).
		RegisterFilterEventMapper(AggregateType, eventstore.EventType(PushedApplicationCreatedEventType), eventstore.GenericEventMapper[ApplicationCreatedPushedEvent]).
		RegisterFilterEventMapper(AggregateType, eventstore.EventType(PushedAuthenticationSucceededOnApplicationEventType), eventstore.GenericEventMapper[AuthenticationSucceededOnApplicationPushedEvent]).
		RegisterFilterEventMapper(AggregateType, eventstore.EventType(PushedInstanceDeletedEventType), eventstore.GenericEventMapper[InstanceDeletedPushedEvent])
}
