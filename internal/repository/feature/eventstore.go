package feature

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, DefaultLoginInstanceEventType, eventstore.GenericEventMapper[SetEvent[Boolean]])
}
