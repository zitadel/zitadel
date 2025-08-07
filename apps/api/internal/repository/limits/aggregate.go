package limits

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	AggregateType    = "limits"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(id, instanceId string) *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:          AggregateType,
			Version:       AggregateVersion,
			ID:            id,
			InstanceID:    instanceId,
			ResourceOwner: instanceId,
		},
	}
}
