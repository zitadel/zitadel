package session

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	AggregateType    = "session"
	AggregateVersion = "v2" //TODO: ?
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
