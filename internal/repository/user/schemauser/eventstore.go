package schemauser

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, CreatedType, eventstore.GenericEventMapper[CreatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, UpdatedType, eventstore.GenericEventMapper[UpdatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, DeletedType, eventstore.GenericEventMapper[DeletedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, LockedType, eventstore.GenericEventMapper[LockedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, UnlockedType, eventstore.GenericEventMapper[UnlockedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, ActivatedType, eventstore.GenericEventMapper[ActivatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, DeactivatedType, eventstore.GenericEventMapper[DeactivatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, EmailUpdatedType, eventstore.GenericEventMapper[EmailUpdatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, EmailCodeAddedType, eventstore.GenericEventMapper[EmailCodeAddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, EmailCodeSentType, eventstore.GenericEventMapper[EmailCodeSentEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, EmailVerifiedType, eventstore.GenericEventMapper[EmailVerifiedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, EmailVerificationFailedType, eventstore.GenericEventMapper[EmailVerificationFailedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PhoneUpdatedType, eventstore.GenericEventMapper[PhoneUpdatedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PhoneCodeAddedType, eventstore.GenericEventMapper[PhoneCodeAddedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PhoneCodeSentType, eventstore.GenericEventMapper[PhoneCodeSentEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PhoneVerifiedType, eventstore.GenericEventMapper[PhoneVerifiedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, PhoneVerificationFailedType, eventstore.GenericEventMapper[PhoneVerificationFailedEvent])
}
