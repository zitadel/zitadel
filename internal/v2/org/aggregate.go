package org

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/v2/eventstore"
)

const (
	AggregateType   = "org"
	eventTypePrefix = AggregateType + "."
)

func NewAggregate(ctx context.Context, id string) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:       id,
		Type:     AggregateType,
		Instance: authz.GetInstance(ctx).InstanceID(),
		Owner:    authz.GetCtxData(ctx).OrgID,
	}
}
