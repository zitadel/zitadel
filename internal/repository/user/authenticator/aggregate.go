package authenticator

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	AggregateType    = "authenticator"
	AggregateVersion = "v1"

	eventPrefix = "authenticator."
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
