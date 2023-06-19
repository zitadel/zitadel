package milestone

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	AggregateType    = "milestone"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

// Each data point receives its own aggregate
func newAggregate(id, instanceId, resourceOwner string) *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:          AggregateType,
			Version:       AggregateVersion,
			ID:            id,
			InstanceID:    instanceId,
			ResourceOwner: resourceOwner,
		},
	}
}
