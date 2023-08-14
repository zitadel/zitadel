package milestone

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	AggregateType    = "milestone"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(ctx context.Context, id string) *Aggregate {
	instanceID := authz.GetInstance(ctx).InstanceID()
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:          AggregateType,
			Version:       AggregateVersion,
			ID:            id,
			ResourceOwner: instanceID,
			InstanceID:    instanceID,
		},
	}
}
