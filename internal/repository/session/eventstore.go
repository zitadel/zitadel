package session

import "github.com/zitadel/zitadel/internal/eventstore"

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, AddedType, AddedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserCheckedType, UserCheckedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasswordCheckedType, PasswordCheckedEventMapper).
		RegisterFilterEventMapper(AggregateType, IntentCheckedType, IntentCheckedEventMapper).
		RegisterFilterEventMapper(AggregateType, WebAuthNChallengedType, eventstore.GenericEventMapper[WebAuthNChallengedEvent]).
		RegisterFilterEventMapper(AggregateType, WebAuthNCheckedType, eventstore.GenericEventMapper[WebAuthNCheckedEvent]).
		RegisterFilterEventMapper(AggregateType, TOTPCheckedType, eventstore.GenericEventMapper[TOTPCheckedEvent]).
		RegisterFilterEventMapper(AggregateType, OTPSMSChallengedType, eventstore.GenericEventMapper[OTPSMSChallengedEvent]).
		RegisterFilterEventMapper(AggregateType, OTPSMSSentType, eventstore.GenericEventMapper[OTPSMSSentEvent]).
		RegisterFilterEventMapper(AggregateType, OTPSMSCheckedType, eventstore.GenericEventMapper[OTPSMSCheckedEvent]).
		RegisterFilterEventMapper(AggregateType, OTPEmailChallengedType, eventstore.GenericEventMapper[OTPEmailChallengedEvent]).
		RegisterFilterEventMapper(AggregateType, OTPEmailSentType, eventstore.GenericEventMapper[OTPEmailSentEvent]).
		RegisterFilterEventMapper(AggregateType, OTPEmailCheckedType, eventstore.GenericEventMapper[OTPEmailCheckedEvent]).
		RegisterFilterEventMapper(AggregateType, TokenSetType, TokenSetEventMapper).
		RegisterFilterEventMapper(AggregateType, MetadataSetType, MetadataSetEventMapper).
		RegisterFilterEventMapper(AggregateType, LifetimeSetType, eventstore.GenericEventMapper[LifetimeSetEvent]).
		RegisterFilterEventMapper(AggregateType, TerminateType, TerminateEventMapper)
}
