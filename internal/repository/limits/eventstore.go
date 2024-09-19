package limits

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, SetEventType, SetEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ResetEventType, ResetEventMapper)
}
