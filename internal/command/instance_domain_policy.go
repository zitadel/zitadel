package command

import (
	"context"

	"github.com/caos/zitadel/internal/command/preparation"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

func AddDefaultDomainPolicy(
	a *instance.Aggregate,
	userLoginMustBeDomain,
	validateOrgDomains bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			//TODO: check if already exists
			return []eventstore.Command{
				instance.NewDomainPolicyAddedEvent(ctx, &a.Aggregate,
					userLoginMustBeDomain,
					validateOrgDomains,
				),
			}, nil
		}, nil
	}
}
