package instance

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	instanceEventTypePrefix = eventstore.EventType("instance.")
)

const (
	AggregateType    = "instance"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(instanceID string) *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			InstanceID:    instanceID,
			Type:          AggregateType,
			Version:       AggregateVersion,
			ID:            instanceID,
			ResourceOwner: instanceID,
		},
	}
}
