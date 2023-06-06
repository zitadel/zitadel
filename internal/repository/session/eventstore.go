package session

import "github.com/zitadel/zitadel/internal/eventstore"

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, AddedType, AddedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserCheckedType, UserCheckedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasswordCheckedType, PasswordCheckedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasskeyChallengedType, eventstore.GenericEventMapper[PasskeyChallengedEvent]).
		RegisterFilterEventMapper(AggregateType, PasskeyCheckedType, eventstore.GenericEventMapper[PasskeyCheckedEvent]).
		RegisterFilterEventMapper(AggregateType, TokenSetType, TokenSetEventMapper).
		RegisterFilterEventMapper(AggregateType, MetadataSetType, MetadataSetEventMapper).
		RegisterFilterEventMapper(AggregateType, TerminateType, TerminateEventMapper)
}
