package command

import (
	"context"

	"github.com/caos/zitadel/internal/command/preparation"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

func AddDefaultLockoutPolicy(
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
