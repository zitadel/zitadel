package limits

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, SetEventType, SetEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ResetEventType, ResetEventMapper)
}
