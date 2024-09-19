package milestone

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, PushedEventType, PushedEventMapper)
}
