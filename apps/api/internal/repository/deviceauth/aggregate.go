package deviceauth

import "github.com/zitadel/zitadel/internal/eventstore"

const (
	AggregateType    = "device_auth"
	AggregateVersion = "v1"
)

func NewAggregate(aggrID, instanceID string) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:   aggrID,
		Type: AggregateType,
		// we use the id because we don't know the resource owner yet
		ResourceOwner: instanceID,
		InstanceID:    instanceID,
		Version:       AggregateVersion,
	}
}
