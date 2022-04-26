package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func AddPasswordAgePolicy(
	a *instance.Aggregate,
	expireWarnDays,
	maxAgeDays uint64,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			//TODO: check if already exists
			return []eventstore.Command{
				instance.NewPasswordAgePolicyAddedEvent(ctx, &a.Aggregate,
					expireWarnDays,
					maxAgeDays,
				),
			}, nil
		}, nil
	}
}
