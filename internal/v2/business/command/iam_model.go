package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

func IAMAggregateFromWriteModel(wm *eventstore.WriteModel) *iam.Aggregate {
	return &iam.Aggregate{
		Aggregate: *eventstore.AggregateFromWriteModel(wm, iam.AggregateType, iam.AggregateVersion),
	}
}
