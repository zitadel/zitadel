package execution

import "github.com/zitadel/zitadel/internal/eventstore"

const (
	AggregateType    = "execution"
	AggregateVersion = "v1"
)

func NewAggregate(aggrID, resourceOwner, instanceID string) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            aggrID,
		Type:          AggregateType,
		ResourceOwner: resourceOwner,
		InstanceID:    instanceID,
		Version:       AggregateVersion,
	}
}
