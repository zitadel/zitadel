package sessionlogout

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	AggregateType    = "session_logout"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(id, instanceID string) *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:          AggregateType,
			Version:       AggregateVersion,
			ID:            id,
			ResourceOwner: instanceID,
			InstanceID:    instanceID,
		},
	}
}
