package deviceauth

import "github.com/zitadel/zitadel/internal/eventstore"

const (
	AggregateType    = "device_auth"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

/*
Note: I've included this as per other "repository" packages,
but somehow it isn't used or required anywhere?
*/
func NewAggregate(id string) *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:    AggregateType,
			Version: AggregateVersion,
			ID:      id,
		},
	}
}
