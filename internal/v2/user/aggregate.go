package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

const (
	AggregateType = "user"
	humanPrefix   = AggregateType + ".human"
	machinePrefix = AggregateType + ".machine"
)

func NewAggregate(ctx context.Context, id string) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:       id,
		Type:     AggregateType,
		Instance: authz.GetInstance(ctx).InstanceID(),
		Owner:    authz.GetCtxData(ctx).OrgID,
	}
}
