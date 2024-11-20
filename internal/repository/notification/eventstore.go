package notification

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, RequestedType, eventstore.GenericEventMapper[RequestedEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, SentType, eventstore.GenericEventMapper[SentEvent])
	eventstore.RegisterFilterEventMapper(AggregateType, FailedType, eventstore.GenericEventMapper[FailedEvent])

}
