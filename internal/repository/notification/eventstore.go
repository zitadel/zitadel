package notification

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, RequestedType, eventstore.GenericEventMapper[RequestedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, SentType, eventstore.GenericEventMapper[SentEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, RetryRequestedType, eventstore.GenericEventMapper[RetryRequestedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, CanceledType, eventstore.GenericEventMapper[CanceledEvent])
}
