package instance

import (
	"context"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

func AddLockoutPolicy(
	a *instance.Aggregate,
	maxAttempts uint64,
	showLockoutFailure bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			//TODO: check if already exists
			return []eventstore.Command{
				instance.NewLockoutPolicyAddedEvent(ctx, &a.Aggregate, maxAttempts, showLockoutFailure),
			}, nil
		}, nil
	}
}
