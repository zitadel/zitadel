package iam

import (
	"context"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

func AddPrivacyPolicy(
	a *iam.Aggregate,
	tosLink,
	privacyLink string,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				iam.NewPrivacyPolicyAddedEvent(ctx, &a.Aggregate, tosLink, privacyLink),
			}, nil
		}, nil
	}
}
