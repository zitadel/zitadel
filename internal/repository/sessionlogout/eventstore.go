package sessionlogout

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	BackChannelLogoutRegisteredEventMapper = eventstore.GenericEventMapper[BackChannelLogoutRegisteredEvent]
	BackChannelLogoutSentEventMapper       = eventstore.GenericEventMapper[BackChannelLogoutSentEvent]

	// Federated logout event mappers
	FederatedLogoutStartedEventMapper              = eventstore.GenericEventMapper[StartedEvent]
	FederatedLogoutSAMLRequestCreatedEventMapper   = eventstore.GenericEventMapper[SAMLRequestCreatedEvent]
	FederatedLogoutSAMLResponseReceivedEventMapper = eventstore.GenericEventMapper[SAMLResponseReceivedEvent]
	FederatedLogoutCompletedEventMapper            = eventstore.GenericEventMapper[CompletedEvent]
	FederatedLogoutFailedEventMapper               = eventstore.GenericEventMapper[FailedEvent]
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, BackChannelLogoutRegisteredType, BackChannelLogoutRegisteredEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, BackChannelLogoutSentType, BackChannelLogoutSentEventMapper)

	// Register federated logout event mappers
	eventstore.RegisterFilterEventMapper(AggregateType, StartedEventType, FederatedLogoutStartedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, SAMLRequestCreatedEventType, FederatedLogoutSAMLRequestCreatedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, SAMLResponseReceivedEventType, FederatedLogoutSAMLResponseReceivedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, CompletedEventType, FederatedLogoutCompletedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, FailedEventType, FederatedLogoutFailedEventMapper)
}
