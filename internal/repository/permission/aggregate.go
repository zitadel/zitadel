package permission

import "github.com/zitadel/zitadel/internal/eventstore"

const (
	AggregateType    eventstore.AggregateType = "permission"
	AggregateVersion eventstore.Version       = "v1"
)

func NewAggregate(aggregateID string) *eventstore.Aggregate {
	var instanceID string
	if aggregateID != "SYSTEM" {
		instanceID = aggregateID
	}
	return &eventstore.Aggregate{
		ID:            aggregateID,
		Type:          AggregateType,
		ResourceOwner: aggregateID,
		InstanceID:    instanceID,
		Version:       AggregateVersion,
	}
}
