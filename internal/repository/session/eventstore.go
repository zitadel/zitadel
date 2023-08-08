package session

import "github.com/zitadel/zitadel/internal/eventstore"

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, AddedType, AddedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserCheckedType, UserCheckedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasswordCheckedType, PasswordCheckedEventMapper).
		RegisterFilterEventMapper(AggregateType, IntentCheckedType, IntentCheckedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasskeyChallengedType, eventstore.GenericEventMapper[PasskeyChallengedEvent]).
		RegisterFilterEventMapper(AggregateType, PasskeyCheckedType, eventstore.GenericEventMapper[PasskeyCheckedEvent]).
		RegisterFilterEventMapper(AggregateType, U2FChallengedType, eventstore.GenericEventMapper[U2FChallengedEvent]).
		RegisterFilterEventMapper(AggregateType, U2FCheckedType, eventstore.GenericEventMapper[U2FCheckedEvent]).
		RegisterFilterEventMapper(AggregateType, TokenSetType, TokenSetEventMapper).
		RegisterFilterEventMapper(AggregateType, MetadataSetType, MetadataSetEventMapper).
		RegisterFilterEventMapper(AggregateType, TerminateType, TerminateEventMapper)
}
