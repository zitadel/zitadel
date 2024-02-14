package session

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, AddedType, AddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, UserCheckedType, UserCheckedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, PasswordCheckedType, PasswordCheckedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, IntentCheckedType, IntentCheckedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, WebAuthNChallengedType, eventstore.GenericEventMapper[WebAuthNChallengedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, WebAuthNCheckedType, eventstore.GenericEventMapper[WebAuthNCheckedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, TOTPCheckedType, eventstore.GenericEventMapper[TOTPCheckedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, OTPSMSChallengedType, eventstore.GenericEventMapper[OTPSMSChallengedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, OTPSMSSentType, eventstore.GenericEventMapper[OTPSMSSentEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, OTPSMSCheckedType, eventstore.GenericEventMapper[OTPSMSCheckedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, OTPEmailChallengedType, eventstore.GenericEventMapper[OTPEmailChallengedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, OTPEmailSentType, eventstore.GenericEventMapper[OTPEmailSentEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, OTPEmailCheckedType, eventstore.GenericEventMapper[OTPEmailCheckedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, TokenSetType, TokenSetEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, MetadataSetType, MetadataSetEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, LifetimeSetType, eventstore.GenericEventMapper[LifetimeSetEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, TerminateType, TerminateEventMapper)
}
