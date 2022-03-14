package iam

import (
	"context"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

func AddOrgIAMPolicy(
	a *iam.Aggregate,
	userLoginMustBeDomain bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				iam.NewOrgIAMPolicyAddedEvent(ctx, &a.Aggregate,
					userLoginMustBeDomain,
				),
			}, nil
		}, nil
	}
}
