package org

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

var (
	AggregateType = "org"
)

func NewAggregate(ctx context.Context, id string) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:       id,
		Type:     "org",
		Instance: authz.GetInstance(ctx).InstanceID(),
		Owner:    authz.GetCtxData(ctx).OrgID,
	}
}
