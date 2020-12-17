package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func AggregateFromWriteModel(wm *eventstore.WriteModel) *user.Aggregate {
	return &user.Aggregate{
		Aggregate: *eventstore.AggregateFromWriteModel(wm, user.AggregateType, user.AggregateVersion),
	}
}
