package debug_events

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix = eventstore.EventType("debug.")
)

const (
	AggregateType    = "debug"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(id, resourceOwner string) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		Type:          AggregateType,
		Version:       AggregateVersion,
		ID:            id,
		ResourceOwner: resourceOwner,
	}
}
