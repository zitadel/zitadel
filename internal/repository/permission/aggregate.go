package permission

import "github.com/zitadel/zitadel/internal/eventstore"

const (
	AggregateType    eventstore.AggregateType = "permission"
	AggregateVersion eventstore.Version       = "v1"
)

func NewAggregate(instanceID string) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            instanceID,
		Type:          AggregateType,
		ResourceOwner: instanceID,
		InstanceID:    instanceID,
		Version:       AggregateVersion,
	}
}
