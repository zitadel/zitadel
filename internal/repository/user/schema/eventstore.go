package schema

import "github.com/zitadel/zitadel/internal/eventstore"

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, CreatedType, eventstore.GenericEventMapper[CreatedEvent])
}
