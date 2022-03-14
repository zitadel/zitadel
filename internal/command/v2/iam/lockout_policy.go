package iam

import (
	"context"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

func AddLockoutPolicy(
	a *iam.Aggregate,
	maxAttempts uint64,
	showLockoutFailure bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				iam.NewLockoutPolicyAddedEvent(ctx, &a.Aggregate, maxAttempts, showLockoutFailure),
			}, nil
		}, nil
	}
}
