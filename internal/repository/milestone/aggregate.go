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

func NewAggregate(id, resourceOwner, instanceID string) *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:          AggregateType,
			Version:       AggregateVersion,
			ID:            id,
			ResourceOwner: resourceOwner,
			InstanceID:    instanceID,
		},
	}
}
