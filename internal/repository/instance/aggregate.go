package instance

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
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

func NewAggregate() *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:          AggregateType,
			Version:       AggregateVersion,
			ID:            domain.IAMID,
			ResourceOwner: domain.IAMID,
		},
	}
}
