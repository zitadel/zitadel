package webkey

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	AggregateType    = "web_key"
	AggregateVersion = "v1"
)

func NewAggregate(id, resourceOwner string) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		Type:          AggregateType,
		Version:       AggregateVersion,
		ID:            id,
		ResourceOwner: resourceOwner,
	}
}

func AggregateFromWriteModel(ctx context.Context, wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModelCtx(ctx, wm, AggregateType, AggregateVersion)
}
