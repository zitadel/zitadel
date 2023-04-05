package deviceauth

import "github.com/zitadel/zitadel/internal/eventstore"

/*
Note: I've included this as per other "repository" packages,
but somehow it isn't used or required anywhere?
*/

const (
	AggregateType    = "device_auth"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(id string) *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:    AggregateType,
			Version: AggregateVersion,
			ID:      id,
		},
	}
}
