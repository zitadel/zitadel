package instance

import (
	"context"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
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
