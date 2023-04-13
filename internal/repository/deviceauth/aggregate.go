package deviceauth

import "github.com/zitadel/zitadel/internal/eventstore"

const (
	AggregateType    = "device_auth"
	AggregateVersion = "v1"
)

func NewAggregate(aggrID, instanceID string) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            aggrID,
		Type:          AggregateType,
		ResourceOwner: instanceID,
		InstanceID:    instanceID,
		Version:       AggregateVersion,
	}
}
