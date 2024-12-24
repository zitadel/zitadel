package group

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	AggregateType    = "group"
	AggregateVersion = "v1" // v1 --> or v2 ?
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(id, resourceOwner string) *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:          AggregateType,
			Version:       AggregateVersion,
			ID:            id,
			ResourceOwner: resourceOwner,
		},
	}
}
