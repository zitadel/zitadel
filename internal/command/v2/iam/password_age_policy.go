package iam

import (
	"context"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

func AddPasswordAgePolicy(
	a *iam.Aggregate,
	expireWarnDays,
	maxAgeDays uint64,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				iam.NewPasswordAgePolicyAddedEvent(ctx, &a.Aggregate,
					expireWarnDays,
					maxAgeDays,
				),
			}, nil
		}, nil
	}
}
