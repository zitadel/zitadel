package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

func ORGAggregateFromWriteModel(wm *eventstore.WriteModel) *org.Aggregate {
	return &org.Aggregate{
		Aggregate: *eventstore.AggregateFromWriteModel(wm, org.AggregateType, org.AggregateVersion),
	}
}
