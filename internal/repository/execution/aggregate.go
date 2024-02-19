package target

import "github.com/zitadel/zitadel/internal/eventstore"

const (
	AggregateType    = "execution"
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
