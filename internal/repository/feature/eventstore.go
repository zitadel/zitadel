package feature

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, DefaultLoginInstanceEventType, eventstore.GenericEventMapper[SetEvent[Boolean]])
}
