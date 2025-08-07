package authrequest

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, AddedType, AddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, SessionLinkedType, SessionLinkedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, CodeAddedType, CodeAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, CodeExchangedType, CodeExchangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, FailedType, FailedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, SucceededType, SucceededEventMapper)
}
