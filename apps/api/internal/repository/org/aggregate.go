package org

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	orgEventTypePrefix = eventstore.EventType("org.")
)

const (
	AggregateType    = "org"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}

func NewAggregate(id string) *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Type:          AggregateType,
			Version:       AggregateVersion,
			ID:            id,
			ResourceOwner: id,
		},
	}
}

func AggregateFromWriteModel(ctx context.Context, wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModelCtx(ctx, wm, AggregateType, AggregateVersion)
}
